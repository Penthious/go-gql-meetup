package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/graphql/dataloaders"
	"github.com/penthious/go-gql-meetup/graphql/resolvers"
	"github.com/penthious/go-gql-meetup/server/directives"
	"github.com/penthious/go-gql-meetup/server/middleware"
	"github.com/rs/cors"
	"net/http"
)


func SetupRouter(d *domain.Domain) *chi.Mux {
	r := chi.NewRouter()
	middleware.SetupMiddleware(r, d)
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)
	//r.Use()

	c := resolvers.Config{
		Resolvers: &resolvers.Resolver{Domain: *d},
	}
	directives.SetDirectives(&c)
	srv := handler.NewDefaultServer(resolvers.NewExecutableSchema(c))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "localhost:8080"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	r.Handle("/", playground.Handler("Meetups gql", "/query"))
	r.Handle("/query", dataloaders.DataloaderMiddleware(d, srv))
	return r
}