package router

import (
	"github.com/BillyBones007/loyalty-service/internal/transport/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// TODO: routers
	r.Get("/", handlers.BadHandler)
	r.Delete("/", handlers.BadHandler)
	r.Put("/", handlers.BadHandler)
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handlers.RegisterUserHandler)
		r.Post("/login", handlers.AuthUserHandler)
		r.Post("/orders", handlers.UploadNumberOrderHandler)
		r.Post("/balance/withdraw", handlers.WithdrawPointsHandler)
		r.Get("/orders", handlers.GetInfoHandler)
		r.Get("/balance", handlers.GetCurrentBalanceHandler)
		r.Get("/withdrawals", handlers.GetWithdrawalsHandler)
	})
	return r
}
