package apiserver

import (
	"log"
	"net/http"
	"wb_l0/internal/app/store"
	

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type APIServer struct {
	config *Config
	router *chi.Mux
	Store *store.Store
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

	sub := NewSubscriber(s)
	go sub.Subscribe()

	log.Println("[Info]: starting api server")
	
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Route("/", func(r chi.Router) {
		r.Get("/get/{id}", s.GetOrderByID())
	})
}

func (s *APIServer) configureStore() error{
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	
	s.Store = st

	err := s.Store.GetOrdersAll()
	if err != nil {
		return err
	}
	
	return nil
}

func (s *APIServer) GetOrderByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		order, ok := s.Store.GetOrderByID(id)
		if !ok {
			log.Println("get order by id !ok")
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(order)
	}
}