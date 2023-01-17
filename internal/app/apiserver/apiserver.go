package apiserver

import (
	"encoding/json"
	"log"
	"net/http"
	"wb_l0/internal/app/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type APIServer struct {
	config *Config
	router *chi.Mux
	store *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		router: chi.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	s.configureRouter()

	if err := s.configureStore(); err != nil{
		log.Printf("[Error]: configure store error: %v", err)
		return err
	}
	log.Println("[Info]: starting api server")
	
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.Use(middleware.Logger)

	s.router.Route("/", func(r chi.Router) {
		r.Get("/get/{id}", s.GetOrderByID())
	})
}

func (s *APIServer) configureStore() error{
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	
	s.store = st
	
	return nil
}

func (s *APIServer) GetOrderByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		order, ok := s.store.GetOrderByID(id)
		if !ok {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		respBody, err := json.Marshal(order)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	}
}