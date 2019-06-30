package router

import (
	"time"

	"SocialPaymentsFeed/src/controller"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

// Initialize chi mux router
func Initialize() *chi.Mux {

	muxRouter := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowCredentials: true,
	})

	muxRouter.Use(cors.Handler)
	muxRouter.Use(middleware.RequestID)
	muxRouter.Use(middleware.RealIP)
	muxRouter.Use(middleware.Logger)
	muxRouter.Use(middleware.Recoverer)
	muxRouter.Use(middleware.Timeout(200 * time.Second))

	muxRouter.Get("/v1/payments/", controller.GetPaymentsController)
	muxRouter.Post("/v1/transfers/", controller.CreateTransferController)
	return muxRouter
}
