package main

import (
	"github.com/go-pg/pg/v9"
	"github.com/penthious/go-gql-meetup/graphql"
	"github.com/penthious/go-gql-meetup/graphql/dataloaders"
	"github.com/penthious/go-gql-meetup/postgres"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	go_gql_meetup "github.com/penthious/go-gql-meetup/graphql/resolvers"
)

const defaultPort = "8080"

func main() {
	DB := postgres.New(&pg.Options{
		User: "tleffew",
		Password:"postgres",
		Database:"meetup_dev",
	})

	defer DB.Close()
	DB.AddQueryHook(postgres.DBLogger{})
	graphqlDB := graphql.DB{
		UserRepo: postgres.NewUserRepo(DB),
		MeetupRepo: postgres.NewMeetupRepo(DB),
	}
	g := &graphql.Domain{DB: graphqlDB}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	c := go_gql_meetup.Config{
		Resolvers: &go_gql_meetup.Resolver{Domain: *g},
	}

	queryHandler := handler.GraphQL(go_gql_meetup.NewExecutableSchema(c))

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", dataloaders.DataloaderMiddleware(g, queryHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
