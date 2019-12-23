package postgres

import (
	"github.com/go-pg/pg/v9"
	"github.com/penthious/go-gql-meetup/models"
)

type UserRepo struct {
	DB *pg.DB
}

func (u *UserRepo) GetByID(id string) (*models.User, error)  {
	var user models.User

	err  := u.DB.Model(&user).Where("id = ?", id).First()

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (u *UserRepo) GetByIDs(ids []string) ([]*models.User, error)  {
	var users []*models.User

	err := u.DB.Model(&users).Where("id in (?)", pg.In(ids)).Select()

	if err != nil {
		return nil, err
	}

	return users, nil

}

func NewUserRepo(DB *pg.DB) *UserRepo {
	return &UserRepo{DB: DB}
}
