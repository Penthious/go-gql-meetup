package database

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/penthious/go-gql-meetup/domain/utils"
	"github.com/penthious/go-gql-meetup/models"
)

type UserRepo struct {
	DB *pg.DB
}

func (u *UserRepo) All() ([]*models.User, error) {
	var users []*models.User

	err := u.DB.Model(&users).Select()

	return users, err
}

func (u *UserRepo) GetByIDs(ids []string) ([]*models.User, error)  {
	var users []*models.User

	err := u.DB.Model(&users).Where("id in (?)", pg.In(ids)).Select()

	return users, err
}

func (u *UserRepo) GetByKey(key, value string) (*models.User, error) {
	user := new(models.User)
	condition := fmt.Sprintf("%v = '%v'", key, value)
	err := u.DB.Model(user).Where(condition).First()

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, utils.ErrNoResult
		}
	}

	return user, err
}

func (u *UserRepo) Create(user *models.User) error {
	_, err := u.DB.Model(user).Returning("*").Insert()

	return err
}

func (u *UserRepo) Update(user *models.User) (*models.User, error) {
	_, err := u.DB.Model(user).Where("id = ?", user.ID).Update()

	return user, err
}

func NewUserRepo(DB *pg.DB) *UserRepo {
	return &UserRepo{DB: DB}
}
