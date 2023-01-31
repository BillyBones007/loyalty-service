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
	r.Get("/", handlers.CoreHandler)
	return r
}
