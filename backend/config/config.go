package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/2foood/backend/internal/infra/jwt"
	"gitlab.com/2foood/backend/internal/infra/oauth"
	oauthProvider "gitlab.com/2foood/backend/internal/infra/oauth/provider"
	"gitlab.com/2foood/backend/internal/infra/otp"
	otpProvider "gitlab.com/2foood/backend/internal/infra/otp/provider"
	"gitlab.com/2foood/backend/internal/infra/payment"
	paymentProvider "gitlab.com/2foood/backend/internal/infra/payment/provider"
	pushProvider "gitlab.com/2foood/backend/internal/infra/push/provider"
	"gitlab.com/2foood/backend/internal/infra/storage"
	storageProvider "gitlab.com/2foood/backend/internal/infra/storage/provider"
	"gitlab.com/2foood/backend/pkg/httpserver"
	"gitlab.com/2foood/backend/pkg/logger"
	"gitlab.com/2foood/backend/pkg/postgres"
	"gitlab.com/2foood/backend/pkg/redis"
)

// everything comes from env vars (envconfig), defaults are in the struct tags.
// see .env.example for the full list
type Config struct {
	App          AppСonfig `envconfig:"APP"`
	HTTP         httpserver.Config
	Logger       logger.Config
	Postgres     postgres.Config
	JWT          jwt.Config
	Redis        redis.Config
	OAuth        oauth.Config
	PhoneOTP     otp.Config `envconfig:"PHONE"`
	EmailOTP     otp.Config `envconfig:"EMAIL"`
	Twilio       otpProvider.TwilioConfig
	SMTP         otpProvider.SMTPConfig
	GoogleOAuth  oauthProvider.GoogleConfig
	Middleware   MiddlewareConfig `envconfig:"MIDDLEWARE"`
	Worker       WorkerConfig     `envconfig:"WORKER"`
	Order        OrderConfig      `envconfig:"ORDER"`
	Payment      payment.Config
	Stripe       paymentProvider.StripeConfig
	Storage      storage.Config
	S3           storageProvider.S3Config
	FCM          pushProvider.FCMConfig
	PushQueue    QueueConfig `envconfig:"PUSH_QUEUE"`
	PaymentQueue QueueConfig `envconfig:"PAYMENT_QUEUE"`
}

// heads up: the first "С" in AppСonfig is a cyrillic letter, not latin. works fine,
// but it's easy to trip on if you ever search for it
type AppСonfig struct {
	Version string `envconfig:"VERSION" required:"true"`
}

type MiddlewareConfig struct {
	RequestTimeout     time.Duration `default:"30s" envconfig:"REQUEST_TIMEOUT"`
	LimiterRequests    int           `default:"100" envconfig:"LIM_REQUESTS"`
	LimiterPeriod      time.Duration `default:"1m" envconfig:"LIM_PERIOD"`
	CorsAllowedOrigins []string      `default:"*" envconfig:"CORS_ALLOWED_ORIGINS"` // * = allow everyone, set real origins in prod
}

type WorkerConfig struct {
	Concurrency int `default:"10" envconfig:"CONCURRENCY"`

	// cron jobs
	OrderCleanerCron string `default:"@every 15m" envconfig:"ORDER_CLEANER_CRON"`
	RoleCleanerCron  string `default:"@every 1h" envconfig:"ROLE_CLEANER_CRON"`

	// timeouts
	PendingOrderTTL time.Duration `default:"1h" envconfig:"PENDING_ORDER_TTL"` // how long an order can stay PENDING (unpaid) before the cleaner cancels it
}

type QueueConfig struct {
	MaxRetries int `default:"5" envconfig:"MAX_RETRIES"`
}

type OrderConfig struct {
	DisplayNumberLength int           `default:"4" envconfig:"DISPLAY_NUMBER_LENGTH"`
	PickupCodeLength    int           `default:"4" envconfig:"PICKUP_CODE_LENGTH"`
	AuditDelay          time.Duration `default:"15m" envconfig:"AUDIT_DELAY"` // how long after confirming a payment we re-check it with stripe
}

func Init() (Config, error) {
	var cfg Config

	// the .env file is optional, in docker/prod everything is passed as real env vars
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
