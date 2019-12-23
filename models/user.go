package models

type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserRepo interface {
	GetByIDs(id []string) ([]*User, error)
	//GetByID(id int64) (*models.User, error)
	//GetByEmail(email string) (*models.User, error)
	//GetByUsername(username string) (*models.User, error)
	//Create(user *models.User) (*models.User, error)
}

