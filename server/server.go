package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-pg/pg/v9"
	"github.com/gorilla/websocket"
	"github.com/penthious/go-gql-meetup/database"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/domain/utils"
	"github.com/penthious/go-gql-meetup/graphql/dataloaders"
	"github.com/vektah/gqlparser/gqlerror"
	"net/http"
	"reflect"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	go_gql_meetup "github.com/penthious/go-gql-meetup/graphql/resolvers"
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
	router := utils.SetupRouter(g)

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing


	c := go_gql_meetup.Config{
		Resolvers: &go_gql_meetup.Resolver{Domain: *g},
	}
	c.Directives.Length = func(ctx context.Context, obj interface{}, next graphql.Resolver, min *int, max *int) (interface{}, error) {


		// @todo: https://github.com/99designs/gqlgen/issues/887 figure out how to get path set correctly in error
		v, _ := next(ctx)
		if reflect.TypeOf(v).String() == "string" {
			// Creates
			if len(v.(string)) < *min {
				return nil, gqlerror.Errorf("format")
			}

		} else if reflect.TypeOf(v).String() == "*string" {
			// Updates
			if len(*v.(*string)) < *min {
				return nil, gqlerror.Errorf("format")
			}
		}

		return v, nil
	}


	//utils.SetDirectives(c)
	srv := handler.NewDefaultServer(go_gql_meetup.NewExecutableSchema(c))
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

	router.Handle("/", playground.Handler("Starwars", "/query"))
	//router.Handle("/query", srv)
	router.Handle("/query", dataloaders.DataloaderMiddleware(g, srv))

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
