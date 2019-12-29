package repos

import "github.com/penthious/go-gql-meetup/models"

type UserRepo interface {
	All() ([]*models.User, error)
	Create(user *models.User) error
	GetByIDs(id []string) ([]*models.User, error)
	GetByKey(key, value string) (*models.User, error)
	Update(user *models.User) (*models.User, error)
}
