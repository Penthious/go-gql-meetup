package models

type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserRepo interface {
	GetByIDs(id []string) ([]*User, error)
	GetByKey(key, value string) (*User, error)
	All() ([]*User, error)
	//Create(user *models.User) (*models.User, error)
}

