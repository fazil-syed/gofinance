package middleware

import (
	"crypto/subtle"
	"fazil-syed/gofinance/internal/config"
	"net/http"
)

type AuthMiddleWare struct {
	config *config.Config
}

func (m *AuthMiddleWare) BasicAuthMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		userMatch := subtle.ConstantTimeCompare([]byte(username), []byte(m.config.Auth.UserName)) == 1
		passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(m.config.Auth.PassWord)) == 1
		if !ok || !userMatch || !passwordMatch {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func NewAuthMiddleWare(cfg *config.Config) *AuthMiddleWare {
	m := &AuthMiddleWare{
		config: cfg,
	}
	return m
}
