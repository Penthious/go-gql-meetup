//package main
//
//import (
//	"github.com/99designs/gqlgen/graphql/handler"
//	"github.com/99designs/gqlgen/graphql/playground"
//	"github.com/go-chi/chi"
//	"github.com/go-pg/pg/v9"
//	"github.com/penthious/go-gql-meetup/graphql"
//	"github.com/penthious/go-gql-meetup/graphql/dataloaders"
//	"github.com/penthious/go-gql-meetup/database"
//	"log"
//	"net/http"
//	"os"
//
//	go_gql_meetup "github.com/penthious/go-gql-meetup/graphql/resolvers"
//)
//
//const defaultPort = "8080"
//
//func main() {
//	DB := database.New(&pg.Options{
//		User: "tleffew",
//		Password:"database",
//		Database:"meetup_dev",
//	})
//
//	defer DB.Close()
//	DB.AddQueryHook(database.DBLogger{})
//	graphqlDB := graphql.DB{
//		UserRepo: database.NewUserRepo(DB),
//		MeetupRepo: database.NewMeetupRepo(DB),
//	}
//	g := &graphql.Domain{DB: graphqlDB}
//
//	//r := utils.SetupRouter(g)
//	r := chi.NewRouter()
//
//
//	port := os.Getenv("PORT")
//	if port == "" {
//		port = defaultPort
//	}
//	c := go_gql_meetup.Config{
//		Resolvers: &go_gql_meetup.Resolver{Domain: *g},
//	}
//
//	srv := handler.NewDefaultServer(go_gql_meetup.NewExecutableSchema(c))
//
//	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
//	r.Handle("/query", dataloaders.DataloaderMiddleware(g, srv))
//
//	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
//	log.Fatal(http.ListenAndServe(":"+port, nil))
//}
package main

import (
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-pg/pg/v9"
	"github.com/gorilla/websocket"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/domain/utils"
	"github.com/penthious/go-gql-meetup/graphql/dataloaders"
	"github.com/penthious/go-gql-meetup/database"
	"net/http"

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
