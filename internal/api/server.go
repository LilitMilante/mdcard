package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	r   *mux.Router
	srv *http.Server
	ph  *PatientHandler
}

func NewServer(port string, ph *PatientHandler) *Server {
	r := mux.NewRouter()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	return &Server{
		r:   r,
		srv: srv,
		ph:  ph,
	}
}

func (s *Server) Start() error {
	s.r.HandleFunc("/patients", s.ph.AddPatient).Methods(http.MethodPost)
	s.r.HandleFunc("/patients", s.ph.Patients).Methods(http.MethodGet)
	s.r.HandleFunc("/patients/{passport_number}", s.ph.PatientByPassportNumber).Methods(http.MethodGet)
	s.r.HandleFunc("/patients/{id}", s.ph.UpdatePatient).Methods(http.MethodPut)
	s.r.HandleFunc("/patients/{id}", s.ph.DeletePatient).Methods(http.MethodDelete)

	s.r.HandleFunc("/patients/cards", s.ph.AddCard).Methods(http.MethodPost)
	s.r.HandleFunc("/patients/cards/{id}", s.ph.UpdateCard).Methods(http.MethodPut)

	s.r.HandleFunc("/sessions", s.ph.Login).Methods(http.MethodPost)

	return s.srv.ListenAndServe()
}
