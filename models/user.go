package models

import (
	"context"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"`

	//CreatedAt time.Time   `json:"createdAt"`
	//UpdatedAt time.Time   `json:"updatedAt"`
	//DeletedAt pg.NullTime `json:"deletedAt" pg:",soft_delete"`
}

func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	password, err := hashPassword(u.Password)
	if err != nil {
		return ctx, err
	}

	u.Password = password

	return ctx, nil
}

func hashPassword(password string) (string, error){
	passwordByte := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)

	return string(passwordHash), err
}
