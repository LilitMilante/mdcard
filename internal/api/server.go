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

	return s.srv.ListenAndServe()
}
