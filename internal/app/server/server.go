package server

import (
	"net/http"

	"github.com/BillyBones007/loyalty-service/internal/db"
	"github.com/BillyBones007/loyalty-service/internal/db/postgres"
	"github.com/BillyBones007/loyalty-service/internal/transport/router"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Config     *Config
	Storage    db.Store
	Routers    *chi.Mux
	HTTPServer *http.Server
}

func NewServer() *Server {
	server := Server{}
	server.Config = initConfig()
	server.Storage = postgres.InitStorage(server.Config.AddrDB)
	server.Routers = router.InitRouter(server.Storage)
	server.HTTPServer = initHTTPServer(server.Config.AddrServ, server.Routers)
	return &server
}

func initHTTPServer(addr string, router *chi.Mux) *http.Server {
	return &http.Server{Addr: addr, Handler: router}
}

func (s *Server) ShutdownServer() {
	// TODO: shutdown functions
}
