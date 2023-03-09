package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseError struct {
	Error string `json:"error"`
}

func SendErr(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(ResponseError{Error: err.Error()})
	if err != nil {
		log.Println(err)
	}
}

func SendJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		SendErr(w, http.StatusInternalServerError, err)
		return
	}
}
