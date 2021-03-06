package middleware

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/penthious/go-gql-meetup/domain"
	"time"
)

func SetupMiddleware(r *chi.Mux, d *domain.Domain) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(6, "application/json"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(Auth(d))
}
