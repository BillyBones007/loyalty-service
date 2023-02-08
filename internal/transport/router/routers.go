package router

import (
	"github.com/BillyBones007/loyalty-service/internal/db"
	"github.com/BillyBones007/loyalty-service/internal/transport/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter(storage db.Store) *chi.Mux {
	h := handlers.Handler{Storage: storage}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", h.BadHandler)
	r.Delete("/", h.BadHandler)
	r.Put("/", h.BadHandler)
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.RegisterUserHandler)
		r.Post("/login", h.AuthUserHandler)
		r.Post("/orders", h.UploadNumberOrderHandler)
		r.Post("/balance/withdraw", h.WithdrawPointsHandler)
		r.Get("/orders", h.GetInfoHandler)
		r.Get("/balance", h.GetCurrentBalanceHandler)
		r.Get("/withdrawals", h.GetWithdrawalsHandler)
	})
	return r
}
