package domain

import (
	"github.com/penthious/go-gql-meetup/domain/repos"
)

//type HaveOwner interface {
//	IsOwner(user *User) bool
//}

type DB struct {
	UserRepo repos.UserRepo
	MeetupRepo repos.MeetupRepo
}


type Domain struct {
	DB DB
}
