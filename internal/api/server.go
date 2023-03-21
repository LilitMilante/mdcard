package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	r      *mux.Router
	srv    *http.Server
	ph     *PatientHandler
	authMw *AuthMiddleware
}

func NewServer(port string, ph *PatientHandler, authMw *AuthMiddleware) *Server {
	r := mux.NewRouter()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	return &Server{
		r:      r,
		srv:    srv,
		ph:     ph,
		authMw: authMw,
	}
}

func (s *Server) Start() error {
	p := s.r.PathPrefix("/patients").Subrouter()
	p.Use(s.authMw.Require)

	p.HandleFunc("", s.ph.AddPatient).Methods(http.MethodPost)
	p.HandleFunc("", s.ph.Patients).Methods(http.MethodGet)
	p.HandleFunc("/{passport_number}", s.ph.PatientByPassportNumber).Methods(http.MethodGet)
	p.HandleFunc("/{id}", s.ph.UpdatePatient).Methods(http.MethodPut)
	p.HandleFunc("/{id}", s.ph.DeletePatient).Methods(http.MethodDelete)

	p.HandleFunc("/cards", s.ph.AddCard).Methods(http.MethodPost)
	p.HandleFunc("/cards/{id}", s.ph.UpdateCard).Methods(http.MethodPut)

	s.r.HandleFunc("/sessions", s.ph.Login).Methods(http.MethodPost)

	return s.srv.ListenAndServe()
}
