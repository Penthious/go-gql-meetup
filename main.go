package main

import (
	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
	"github.com/penthious/go-gql-meetup/database"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/server"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load("./environments/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DB := database.New(&pg.Options{
		User: os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE"),
	})

	defer DB.Close()

	DB.AddQueryHook(database.DBLogger{})
	graphqlDB := domain.DB{
		UserRepo:   database.NewUserRepo(DB),
		MeetupRepo: database.NewMeetupRepo(DB),
		DB: DB,
	}
	d := &domain.Domain{DB: graphqlDB}
	router := server.SetupRouter(d)

	err = http.ListenAndServe(":"+ os.Getenv("SITE_PORT"), router)
	if err != nil {
		panic(err)
	}
}
