package api

import (
	"net/http"

	"medical-card/internal/service"
)

type AuthMiddleware struct {
	srv Service
}

func NewAuthMiddleware(srv Service) *AuthMiddleware {
	return &AuthMiddleware{
		srv: srv,
	}
}

func (a *AuthMiddleware) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("ssid")
		if err != nil {
			SendErr(w, http.StatusUnauthorized, service.ErrUnauthorized)
			return
		}

		_, err = a.srv.PatientBySessionID(r.Context(), cookie.Value)
		if err != nil {
			SendErr(w, http.StatusUnauthorized, service.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
