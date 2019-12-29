package domain

import (
	"github.com/go-pg/pg/v9"
	"github.com/penthious/go-gql-meetup/domain/repos"
)

//type HaveOwner interface {
//	IsOwner(user *User) bool
//}

type DB struct {
	UserRepo repos.UserRepo
	MeetupRepo repos.MeetupRepo
	DB *pg.DB
}


type Domain struct {
	DB DB
}
