package apiserver

import (
	"fmt"
	"html/template"
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
		r.Get("/", s.MainPage())
		r.Get("/get", s.GetOrderByID())
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

type Page struct {
	UID string
	Data string
}

func (s *APIServer) GetOrderByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmplt, err := template.ParseFiles("../web/html/index.html")
		if err != nil {
			log.Printf("[Error]:template parsing error: %v", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		id := r.FormValue("id")
		fmt.Println(id)
		order, ok := s.Store.GetOrderByID(id)
		
		event := Page{
			UID:  "ID: " + id,
			Data: string(order),
		}

		if !ok {
			log.Println("get order by id !ok")
			event = Page{
				Data: "id \"" + id + "\" not found",
			}
		}

		err = tmplt.Execute(w, event)
		if err != nil {
			log.Printf("[Error]:template execution error: %v", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
	}
}


func (s *APIServer) MainPage() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		tmplt, err := template.ParseFiles("../web/html/index.html")
		if err != nil {
			log.Printf("[Error]:template parsing error: %v", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		event := Page{}
		err = tmplt.Execute(w, event)
		if err != nil {
			log.Printf("[Error]:template execution error: %v", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
	}
}