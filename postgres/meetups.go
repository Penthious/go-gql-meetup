package postgres

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/penthious/go-gql-meetup/models"
	"github.com/penthious/go-gql-meetup/utils"
)

type MeetupRepo struct {
	DB *pg.DB
}

func (m *MeetupRepo) GetMeetupsForUser(id string) ([]*models.Meetup, error) {
	var meetups []*models.Meetup
	err :=  m.DB.Model(&meetups).Where("user_id = ?", id).Order("id").Select()

	return meetups, err
}

func (m *MeetupRepo) Delete(meetup *models.Meetup) error {
	return m.DB.Delete(meetup)
}

func (m *MeetupRepo) Update(meetup *models.Meetup) (*models.Meetup, error) {
	_, err := m.DB.Model(meetup).Where("id = ?", meetup.ID).Update()

	if err != nil {
		return nil, err
	}

	return meetup, nil
}

func (m *MeetupRepo) Create(meetup *models.Meetup) (*models.Meetup, error) {
	_, err := m.DB.Model(meetup).Returning("*").Insert(meetup)

	if err != nil {
		return nil, err
	}

	return meetup, err
}

func (m *MeetupRepo) All() ([]*models.Meetup, error) {
	var meetups []*models.Meetup

	err := m.DB.Model(&meetups).Order("id").Select()

	if err != nil {
		return nil, err
	}

	return meetups, nil
}

func (m *MeetupRepo) GetByIDs(ids []string) ([]*models.Meetup, error)  {
	var meetups []*models.Meetup

	err  := m.DB.Model(&meetups).Where("id = (?)", pg.In(ids)).Select()

	if err != nil {
		return nil, err
	}

	return meetups, nil

}

func (m *MeetupRepo) GetByKey(key, value string) (*models.Meetup, error) {
	meetup := new(models.Meetup)

	condition := fmt.Sprintf("%v = %v", key, value)

	err := m.DB.Model(meetup).Where(condition).First()

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, utils.ErrNoResult
		}
		return nil, err
	}

	return meetup, nil
}

func NewMeetupRepo(DB *pg.DB) *MeetupRepo {
	return &MeetupRepo{DB: DB}
}
