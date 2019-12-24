package domain

import "github.com/penthious/go-gql-meetup/models"

//type HaveOwner interface {
//	IsOwner(user *User) bool
//}

type DB struct {
	UserRepo models.UserRepo
	MeetupRepo models.MeetupRepo
}


type Domain struct {
	DB DB
}
