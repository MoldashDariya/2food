package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"gitlab.com/2foood/backend/config"
	"gitlab.com/2foood/backend/internal/auth"
	authHTTP "gitlab.com/2foood/backend/internal/auth/transport/http"
	"gitlab.com/2foood/backend/internal/catalog"
	catalogHTTP "gitlab.com/2foood/backend/internal/catalog/transport/http"
	healthHTTP "gitlab.com/2foood/backend/internal/health/transport/http"
	"gitlab.com/2foood/backend/internal/infra/jwt"
	"gitlab.com/2foood/backend/internal/infra/oauth"
	oauthProvider "gitlab.com/2foood/backend/internal/infra/oauth/provider"
	"gitlab.com/2foood/backend/internal/infra/otp"
	otpProvider "gitlab.com/2foood/backend/internal/infra/otp/provider"
	"gitlab.com/2foood/backend/internal/infra/payment"
	paymentProvider "gitlab.com/2foood/backend/internal/infra/payment/provider"
	"gitlab.com/2foood/backend/internal/infra/push"
	pushProvider "gitlab.com/2foood/backend/internal/infra/push/provider"
	"gitlab.com/2foood/backend/internal/infra/queue"
	"gitlab.com/2foood/backend/internal/infra/repository/postgres"
	"gitlab.com/2foood/backend/internal/infra/storage"
	storageProvider "gitlab.com/2foood/backend/internal/infra/storage/provider"
	"gitlab.com/2foood/backend/internal/infra/transactor"
	"gitlab.com/2foood/backend/internal/order"
	orderHTTP "gitlab.com/2foood/backend/internal/order/transport/http"
	"gitlab.com/2foood/backend/internal/pkg/asynqlog"
	customMW "gitlab.com/2foood/backend/internal/pkg/middleware"
	"gitlab.com/2foood/backend/internal/pkg/translator"
	"gitlab.com/2foood/backend/internal/pkg/validation"
	"gitlab.com/2foood/backend/internal/product"
	productHTTP "gitlab.com/2foood/backend/internal/product/transport/http"
	"gitlab.com/2foood/backend/internal/restaurant"
	restaurantHTTP "gitlab.com/2foood/backend/internal/restaurant/transport/http"
	"gitlab.com/2foood/backend/internal/review"
	reviewHTTP "gitlab.com/2foood/backend/internal/review/transport/http"
	"gitlab.com/2foood/backend/internal/upload"
	uploadHTTP "gitlab.com/2foood/backend/internal/upload/transport/http"
	"gitlab.com/2foood/backend/internal/user"
	userHTTP "gitlab.com/2foood/backend/internal/user/transport/http"
	"gitlab.com/2foood/backend/internal/worker"
	"gitlab.com/2foood/backend/pkg/httpserver"
	pg "gitlab.com/2foood/backend/pkg/postgres"
	"gitlab.com/2foood/backend/pkg/redis"
)

// App wires everything together by hand (no DI framework):
// infra -> repositories -> services -> handlers. it's all in Run()
type App struct {
	ctx      context.Context
	log      *slog.Logger
	cfg      config.Config
	handlers *Handlers
	mw       *customMW.Manager
}

// one http handler per feature, built in Run() and used by newRouter()
type Handlers struct {
	Health     *healthHTTP.Handler
	Auth       *authHTTP.Handler
	User       *userHTTP.Handler
	Restaurant *restaurantHTTP.Handler
	Product    *productHTTP.Handler
	Order      *orderHTTP.Handler
	Upload     *uploadHTTP.Handler
	Catalog    *catalogHTTP.Handler
	Review     *reviewHTTP.Handler
}

// todo: looks unused, a.mw (customMW.Manager) is what the router actually uses
type MW struct {
	requireAuth             func(http.Handler) http.Handler
	requireRestaurantAccess func(http.Handler) http.Handler
}

func New(ctx context.Context, log *slog.Logger, cfg config.Config) *App {
	return &App{
		ctx: ctx,
		log: log,
		cfg: cfg,
	}
}

func (a *App) Run() error {
	const op = "app.Run"
	log := a.log.With(slog.String("op", op))

	// ctx gets cancelled on ctrl+c / SIGTERM (docker stop), that's what triggers the shutdown at the bottom
	// context with signal
	ctx, cancel := signal.NotifyContext(a.ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// postgres
	db, err := pg.New(a.cfg.Postgres)
	if err != nil {
		log.Error("failed to initialize database", slog.Any("error", err))
		return err
	}
	defer func() {
		log.Info("closing database connections...")
		db.Close()
	}()

	// redis
	rdb, err := redis.New(a.cfg.Redis)
	if err != nil {
		log.Error("failed to initialize redis", slog.Any("error", err))
		return err
	}
	defer func() {
		log.Info("closing redis connections...")
		rdb.Close()
	}()

	// asynq (background jobs) talks to the same redis as the cache, just through its own client
	// asynq queue setup
	redisAddr := net.JoinHostPort(a.cfg.Redis.Host, a.cfg.Redis.Port)
	asynqOpts := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: a.cfg.Redis.Password,
		DB:       a.cfg.Redis.DB,
	}

	queueClient := queue.NewClient(asynqOpts)
	defer queueClient.Close()

	// adapters so services can enqueue jobs without knowing about asynq
	asyncPushSender := queue.NewPushAdapter(queueClient, a.cfg.PushQueue.MaxRetries)
	asyncPaymentSender := queue.NewPaymentAdapter(queueClient, a.cfg.PaymentQueue.MaxRetries)

	// gives services a way to run several repo calls in one db transaction
	// transactor
	tm := transactor.New(db)

	// translator
	tr := translator.New()

	// phone otp manager
	smsProvider := otpProvider.NewTwilio(a.cfg.Twilio, tr)
	phoneOTPManager := otp.NewManager(a.cfg.PhoneOTP, smsProvider, rdb)

	// email otp manager
	mailProvider, err := otpProvider.NewSMTP(a.cfg.SMTP, tr, a.cfg.EmailOTP.TTL)
	if err != nil {
		log.Error("failed to initialize email otp provider", slog.Any("error", err))
		return err
	}
	emailOTPManager := otp.NewManager(a.cfg.EmailOTP, mailProvider, rdb)

	// token manager
	tokenManager := jwt.NewManager(a.cfg.JWT)

	// oauth manager
	googleOAuthProvider := oauthProvider.NewGoogle(a.cfg.GoogleOAuth)
	oauthManager := oauth.NewManager(a.cfg.OAuth, googleOAuthProvider, rdb)

	// payment manager
	stripeProvider := paymentProvider.NewStripe(a.cfg.Stripe)
	paymentManager := payment.NewManager(a.cfg.Payment, stripeProvider)

	// storage manager
	s3Provider, err := storageProvider.NewS3(a.cfg.S3)
	if err != nil {
		log.Error("failed to initialize storage provider", slog.Any("error", err))
		return err
	}
	storageManager := storage.NewManager(a.cfg.Storage, s3Provider)

	// repositories
	userRepo := postgres.NewUserRepo(db)
	restaurantRepo := postgres.NewRestaurantRepo(db)
	staffRepo := postgres.NewStaffRepo(db)
	productCategoryRepo := postgres.NewProductCategoryRepo(db)
	productRepo := postgres.NewProductRepo(db)
	orderRepo := postgres.NewOrderRepo(db)
	deviceTokenRepo := postgres.NewDeviceTokenRepo(db)
	reviewRepo := postgres.NewReviewRepo(db)

	// catalog repositories
	cityRepo := postgres.NewCityRepo(db)
	cuisineTypeRepo := postgres.NewCuisineTypeRepo(db)

	// push manager
	fcmProvider, err := pushProvider.NewFCM(a.cfg.FCM)
	if err != nil {
		log.Error("failed to initialize fcm provider", slog.Any("error", err))
		return err
	}
	pushManager := push.NewManager(fcmProvider, deviceTokenRepo)

	// business logic lives here, they only see interfaces (easy to mock in tests)
	// services
	authService := auth.NewService(userRepo, phoneOTPManager, emailOTPManager, oauthManager, tokenManager)
	userService := user.NewService(tm, userRepo, deviceTokenRepo, phoneOTPManager, emailOTPManager)
	restService := restaurant.NewService(tm, userRepo, restaurantRepo, staffRepo)
	productService := product.NewService(tm, productRepo)
	// orderService is the big one: creates orders, talks to stripe and pushes
	orderService := order.NewService(
		tm, a.log, orderRepo, productRepo,
		paymentManager, asyncPaymentSender, asyncPushSender,
		a.cfg.Order.DisplayNumberLength, a.cfg.Order.PickupCodeLength, a.cfg.Order.AuditDelay, a.cfg.Payment.Currency,
	)
	uploadService := upload.NewService(storageManager)
	catalogService := catalog.NewService(cityRepo, cuisineTypeRepo, productCategoryRepo)
	reviewService := review.NewService(tm, reviewRepo, orderRepo, restaurantRepo)

	// handlers for the background jobs (they run in this same process)
	// workers
	roleCleaner := worker.NewRoleCleanerHandler(a.log, userRepo)
	orderCleaner := worker.NewOrderCleanerHandler(a.log, orderService, a.cfg.Worker.PendingOrderTTL)
	pushHandler := worker.NewPushHandler(a.log, pushManager, tr, userRepo)
	paymentHandler := worker.NewPaymentHandler(a.log, paymentManager, orderService, orderService)

	// mux setup
	mux := asynq.NewServeMux()
	mux.Handle(queue.TypeRoleCleanup, roleCleaner)
	mux.Handle(queue.TypeOrderCancelExpired, orderCleaner)
	mux.Handle(queue.TypePushSend, pushHandler)

	mux.Handle(queue.TypePaymentConfirm, paymentHandler)
	mux.Handle(queue.TypePaymentCancel, paymentHandler)
	mux.Handle(queue.TypePaymentRefund, paymentHandler)
	mux.Handle(queue.TypePaymentAudit, paymentHandler)

	// queue server
	queueServer := asynq.NewServer(
		asynqOpts,
		asynq.Config{
			// how many jobs can run at the same time
			Concurrency: a.cfg.Worker.Concurrency,
			Logger:      asynqlog.New(a.log),
		},
	)

	// scheduler
	scheduler := asynq.NewScheduler(asynqOpts, &asynq.SchedulerOpts{
		Logger: asynqlog.New(a.log),
	})

	// cron jobs. if registering fails we only log it, the api should still start
	if _, err := scheduler.Register(a.cfg.Worker.OrderCleanerCron, asynq.NewTask(queue.TypeOrderCancelExpired, nil)); err != nil {
		log.Error("failed to register order cleaner cron", slog.Any("error", err))
	}

	if _, err := scheduler.Register(a.cfg.Worker.RoleCleanerCron, asynq.NewTask(queue.TypeRoleCleanup, nil)); err != nil {
		log.Error("failed to register role cleaner cron", slog.Any("error", err))
	}

	// run queue server & scheduler
	go func() {
		log.Info("starting asynq worker server")
		if err := queueServer.Run(mux); err != nil {
			log.Error("asynq worker server stopped", slog.Any("error", err))
		}
	}()

	go func() {
		log.Info("starting asynq scheduler")
		if err := scheduler.Run(); err != nil {
			log.Error("asynq scheduler stopped", slog.Any("error", err))
		}
	}()

	// validator
	validator := validation.New()

	// middleware manager
	a.mw = customMW.NewManager(a.log, tokenManager, staffRepo)

	// handlers
	a.handlers = &Handlers{
		Health:     healthHTTP.NewHandler(a.log),
		Auth:       authHTTP.NewHandler(a.log, authService, validator),
		User:       userHTTP.NewHandler(a.log, userService, validator),
		Restaurant: restaurantHTTP.NewHandler(a.log, restService, restService, validator),
		Product:    productHTTP.NewHandler(a.log, productService, validator),
		Order:      orderHTTP.NewHandler(a.log, orderService, validator),
		Upload:     uploadHTTP.NewHandler(a.log, uploadService, validator),
		Catalog:    catalogHTTP.NewHandler(a.log, catalogService, validator, translator.New()), // todo: makes a 2nd translator, could reuse `tr` from above
		Review:     reviewHTTP.NewHandler(a.log, reviewService, validator),
	}

	// add health checks
	a.handlers.Health.Add("postgres", db.PingContext)
	a.handlers.Health.Add("redis", rdb.Ping)

	// note: the http server, worker and scheduler all run in this one process
	// server setup
	router := a.newRouter()
	httpServer := httpserver.New(router, a.cfg.HTTP)

	log.Info("server starting", slog.String("port", a.cfg.HTTP.Port))
	httpServer.Start()

	// block here until the signal comes, then stop things in order:
	// http first (no new requests), then background jobs. db/redis close via defer
	// graceful shutdown
	<-ctx.Done()

	log.Info("shutting down application...")

	log.Info("closing http server...")
	if err := httpServer.Close(); err != nil {
		log.Error("failed to close http server", slog.Any("error", err))
	}

	log.Info("stopping background jobs...")
	queueServer.Shutdown()
	scheduler.Shutdown()

	log.Info("application stopped")
	return nil
}
