package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/2foood/backend/docs"
	"gitlab.com/2foood/backend/internal/domain"
	"gitlab.com/2foood/backend/internal/pkg/params"
)

func (a *App) newRouter() http.Handler {
	r := chi.NewRouter()

	// order matters here: request id + real ip first so the logger can use them,
	// recoverer catches panics, timeout cancels slow requests
	// base middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(a.mw.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(a.cfg.Middleware.RequestTimeout))
	r.Use(a.mw.Language)

	// per-ip limit, numbers come from config (MIDDLEWARE_LIM_*)
	// rate limiting
	r.Use(httprate.Limit(
		a.cfg.Middleware.LimiterRequests,
		a.cfg.Middleware.LimiterPeriod,
		httprate.WithKeyFuncs(httprate.KeyByIP),
	))

	// cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   a.cfg.Middleware.CorsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// health checks
	r.Group(func(r chi.Router) {
		r.Get("/live", a.handlers.Health.Live)
		r.Get("/ready", a.handlers.Health.Readiness)
		// note: no auth on these, make sure /metrics isn't reachable from outside
		r.Handle("/metrics", promhttp.Handler())
	})

	// api v1 routes
	r.Route("/api/v1", func(r chi.Router) {

		// public routes = how you get a token (phone otp or google)
		// auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/phone/otp", a.handlers.Auth.PhoneSendOTP)
			r.Post("/phone/login", a.handlers.Auth.PhoneLogin)
			r.Post("/google/login", a.handlers.Auth.GoogleLogin)
			r.Post("/google/complete", a.handlers.Auth.CompleteGoogleLogin)
			r.Post("/refresh", a.handlers.Auth.Refresh)

			// these two need a token: attach google to an existing account / send otp to verify an email
			// protected routes
			r.Group(func(r chi.Router) {
				r.Use(a.mw.Auth)
				r.Post("/link/google", a.handlers.Auth.LinkGoogle)
				r.Post("/email/otp", a.handlers.Auth.EmailSendOTP)
			})
		})

		// user routes
		r.Route("/users", func(r chi.Router) {
			r.Use(a.mw.Auth)

			r.Route("/me", func(r chi.Router) {
				r.Get("/", a.handlers.User.Get)
				r.Patch("/", a.handlers.User.Update)
				r.Patch("/phone", a.handlers.User.UpdatePhone)
				r.Patch("/email", a.handlers.User.UpdateEmail)
				r.Delete("/", a.handlers.User.Delete)

				// device tokens
				r.Post("/devices", a.handlers.User.RegisterDeviceToken)
				r.Delete("/devices", a.handlers.User.DeleteDeviceToken)

				// reviews
				r.Get("/reviews", a.handlers.Review.GetByUser)

				// payment history
				r.Get("/payment-history", a.handlers.Order.GetPaymentHistoryByUser)
			})
		})

		// restaurant routes
		r.Route("/restaurants", func(r chi.Router) {
			r.Use(a.mw.Auth)
			r.Get("/", a.handlers.Restaurant.GetAll)

			r.Route("/{"+params.RestaurantID+"}", func(r chi.Router) {
				r.Get("/", a.handlers.Restaurant.GetByID)
				r.Get("/products", a.handlers.Product.GetActiveByRestaurant)
				r.Get("/categories", a.handlers.Product.GetCategoriesByRestaurant)

				r.Route("/reviews", func(r chi.Router) {
					r.Get("/", a.handlers.Review.GetByRestaurant)
					r.Get("/photos", a.handlers.Review.GetPhotosByRestaurant)
				})
			})
		})

		// partner = restaurant staff. RequireRole(partner) first (admins pass too),
		// then RequireRestaurantAccess checks the user actually belongs to that restaurant
		// partner routes
		r.Route("/partner", func(r chi.Router) {
			r.Use(a.mw.Auth)
			r.Use(a.mw.RequireRole(domain.UserRolePartner))
			r.Get("/restaurants", a.handlers.Restaurant.GetByPartner)

			r.Route("/restaurants/{"+params.RestaurantID+"}", func(r chi.Router) {
				r.Use(a.mw.RequireRestaurantAccess)
				r.Patch("/", a.handlers.Restaurant.UpdateInfo)

				r.Route("/staff", func(r chi.Router) {
					r.Get("/", a.handlers.Restaurant.GetStaff)
					r.Post("/", a.handlers.Restaurant.AddStaff)
					r.Patch("/{"+params.UserID+"}/role", a.handlers.Restaurant.UpdateStaffRole)
					r.Patch("/{"+params.UserID+"}/ownership", a.handlers.Restaurant.TransferOwnership)
					r.Delete("/{"+params.UserID+"}", a.handlers.Restaurant.RemoveStaff)
				})

				r.Route("/products", func(r chi.Router) {
					r.Post("/", a.handlers.Product.Create)
					r.Get("/", a.handlers.Product.GetByRestaurant)
					r.Get("/{"+params.ProductID+"}", a.handlers.Product.GetByID)
					r.Patch("/{"+params.ProductID+"}", a.handlers.Product.Update)
					r.Patch("/{"+params.ProductID+"}/status", a.handlers.Product.UpdateStatus)
					r.Patch("/{"+params.ProductID+"}/stock", a.handlers.Product.UpdateStock)
					r.Delete("/{"+params.ProductID+"}", a.handlers.Product.Delete)
				})

				r.Route("/orders", func(r chi.Router) {
					r.Get("/", a.handlers.Order.GetByRestaurant)
					r.Get("/active", a.handlers.Order.GetActiveByRestaurant)
					r.Get("/pickup", a.handlers.Order.GetByPickupCode)
					r.Get("/{"+params.OrderID+"}/history", a.handlers.Order.GetStatusHistoryByRestaurant)
					r.Patch("/{"+params.OrderID+"}/ready", a.handlers.Order.MarkAsReady)
					r.Patch("/{"+params.OrderID+"}/collect", a.handlers.Order.MarkAsCollected)
					r.Delete("/{"+params.OrderID+"}", a.handlers.Order.CancelByRestaurant)
				})

				r.Route("/reviews", func(r chi.Router) {
					r.Get("/", a.handlers.Review.GetByRestaurantAll)
					r.Post("/{"+params.ReviewID+"}/reply", a.handlers.Review.AddReply)
				})
			})
		})

		// admin only
		// admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(a.mw.Auth)
			r.Use(a.mw.RequireRole(domain.UserRoleAdmin))

			// restaurant routes
			r.Route("/restaurants", func(r chi.Router) {
				r.Post("/", a.handlers.Restaurant.Create)
				r.Delete("/{"+params.RestaurantID+"}", a.handlers.Restaurant.Delete)

				r.Route("/{"+params.RestaurantID+"}", func(r chi.Router) {
					r.Patch("/owner", a.handlers.Restaurant.AssignOwner)
				})
			})

			// upload routes
			r.Delete("/upload", a.handlers.Upload.Delete)

			// catalog routes
			r.Route("/catalog", func(r chi.Router) {
				// cities routes
				r.Post("/cities", a.handlers.Catalog.CreateCity)
				r.Delete("/cities/{"+params.CityID+"}", a.handlers.Catalog.DeleteCity)
				// cuisine types routes
				r.Post("/cuisine-types", a.handlers.Catalog.CreateCuisineType)
				r.Delete("/cuisine-types/{"+params.CuisineTypeID+"}", a.handlers.Catalog.DeleteCuisineType)
				// product categories routes
				r.Post("/product-categories", a.handlers.Catalog.CreateProductCategory)
				r.Delete("/product-categories/{"+params.ProductCategoryID+"}", a.handlers.Catalog.DeleteProductCategory)
			})

			// review routes
			r.Route("/reviews", func(r chi.Router) {
				r.Get("/", a.handlers.Review.GetAll)
				r.Patch("/{"+params.ReviewID+"}/hide", a.handlers.Review.Hide)
				r.Patch("/{"+params.ReviewID+"}/unhide", a.handlers.Review.Unhide)
			})
		})

		// customer side of orders. the service checks that the order belongs to the user
		// order routes
		r.Route("/orders", func(r chi.Router) {
			r.Use(a.mw.Auth)
			r.Post("/", a.handlers.Order.Create)
			r.Get("/", a.handlers.Order.GetByUser)

			r.Route("/{"+params.OrderID+"}", func(r chi.Router) {
				r.Get("/", a.handlers.Order.GetByID)
				r.Get("/history", a.handlers.Order.GetStatusHistoryByUser)
				r.Delete("/", a.handlers.Order.CancelByUser)
				r.Post("/review", a.handlers.Review.Create)
				r.Get("/review", a.handlers.Review.GetByOrder)
				r.Get("/payment-history", a.handlers.Order.GetPaymentHistoryByOrder)
			})
		})

		// no Auth middleware on purpose - stripe calls this, not a user.
		// the request is verified by the Stripe-Signature header inside the order service
		// webhooks routes
		r.Route("/webhooks", func(r chi.Router) {
			r.Route("/stripe", func(r chi.Router) {
				r.Post("/payment/confirm", a.handlers.Order.StripeConfirmPayment)
			})
		})

		// upload routes
		r.Route("/upload", func(r chi.Router) {
			r.Use(a.mw.Auth)
			r.Post("/presign", a.handlers.Upload.Presign)
		})

		// catalog routes
		r.Route("/catalog", func(r chi.Router) {
			r.Use(a.mw.Auth)
			// cities routes
			r.Get("/cities", a.handlers.Catalog.GetCities)
			r.Get("/cities/{"+params.CityID+"}", a.handlers.Catalog.GetCityByID)
			// cuisine types routes
			r.Get("/cuisine-types", a.handlers.Catalog.GetCuisineTypes)
			r.Get("/cuisine-types/{"+params.CuisineTypeID+"}", a.handlers.Catalog.GetCuisineTypeByID)
			// product categories routes
			r.Get("/product-categories", a.handlers.Catalog.GetProductCategories)
			r.Get("/product-categories/{"+params.ProductCategoryID+"}", a.handlers.Catalog.GetProductCategoryByID)
		})
	})

	// swagger json + scalar ui, generated from the annotations above the handlers
	// documentation routes
	r.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(docs.SwaggerJSON)
	})

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(docs.ScalarHTML)
	})

	return r
}
