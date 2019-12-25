package utils

import (
	"github.com/go-chi/chi"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/domain/middleware"
	"github.com/rs/cors"
)


func SetupRouter(d *domain.Domain) *chi.Mux {
	r := chi.NewRouter()
	middleware.SetupMiddleware(r)
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)
	//r.Use()

	return r
}