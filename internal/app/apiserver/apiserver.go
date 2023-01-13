package apiserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type APIServer struct {
	config *Config
	router *chi.Mux
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		router: chi.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	s.configureRouter()
	log.Println("[Info]: starting api server")
	
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.Use(middleware.Logger)

	s.router.Route("/", func(r chi.Router) {
		r.Get("/get/{id}", GetOrderByID())
	})
}