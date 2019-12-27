package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-pg/pg/v9"
	"github.com/gorilla/websocket"
	"github.com/penthious/go-gql-meetup/database"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/graphql/dataloaders"
	"github.com/penthious/go-gql-meetup/graphql/resolvers"
	"github.com/penthious/go-gql-meetup/server"
	"github.com/penthious/go-gql-meetup/server/directives"
	"net/http"
)

func main() {
	DB := database.New(&pg.Options{
		User: "tleffew",
		Password:"database",
		Database:"meetup_dev",
	})

	defer DB.Close()

	DB.AddQueryHook(database.DBLogger{})
	graphqlDB := domain.DB{
		UserRepo:   database.NewUserRepo(DB),
		MeetupRepo: database.NewMeetupRepo(DB),
	}
	g := &domain.Domain{DB: graphqlDB}
	router := server.SetupRouter(g)

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing


	c := resolvers.Config{
		Resolvers: &resolvers.Resolver{Domain: *g},
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

	router.Handle("/", playground.Handler("Meetups gql", "/query"))
	//router.Handle("/query", srv)
	router.Handle("/query", dataloaders.DataloaderMiddleware(g, srv))

	err := http.ListenAndServe(":8081", router)
	if err != nil {
		panic(err)
	}
}
