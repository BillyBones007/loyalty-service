package router

import (
	"github.com/BillyBones007/loyalty-service/internal/transport/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", h.BadHandler)
	r.Delete("/", h.BadHandler)
	r.Put("/", h.BadHandler)
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.RegisterUserHandler)
		r.Post("/login", h.AuthUserHandler)
		r.With(h.CheckCookies).Post("/orders", h.UploadNumberOrderHandler)
		r.With(h.CheckCookies).Post("/balance/withdraw", h.WithdrawPointsHandler)
		r.With(h.CheckCookies).Get("/orders", h.GetOrdersInfoHandler)
		r.With(h.CheckCookies).Get("/balance", h.GetCurrentBalanceHandler)
		r.With(h.CheckCookies).Get("/withdrawals", h.GetWithdrawalsHandler)
	})
	return r
}
