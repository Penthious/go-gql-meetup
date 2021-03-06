package database

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/penthious/go-gql-meetup/domain/utils"
	"github.com/penthious/go-gql-meetup/models"
)

type MeetupRepo struct {
	DB *pg.DB
}

func (m *MeetupRepo) GetMeetupsByFilter(filter *models.MeetupFilterPayload) ([]*models.Meetup, error) {
	var meetups []*models.Meetup

	query := m.DB.Model(&meetups).Order("id")

	if filter != nil {
		if filter.Name != nil {
			query.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", *filter.Name))
		}
		if filter.Description != nil {
			query.Where("description ILIKE ?", fmt.Sprintf("%%%s%%", *filter.Description))
		}
	}

	err := query.Select()

	return meetups, err
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

	return meetup, err
}

func (m *MeetupRepo) Create(meetup *models.Meetup) (*models.Meetup, error) {
	_, err := m.DB.Model(meetup).Returning("*").Insert(meetup)

	return meetup, err
}

func (m *MeetupRepo) All() ([]*models.Meetup, error) {
	var meetups []*models.Meetup

	err := m.DB.Model(&meetups).Order("id").Select()

	return meetups, err
}

func (m *MeetupRepo) GetByIDs(ids []string) ([]*models.Meetup, error)  {
	var meetups []*models.Meetup

	err  := m.DB.Model(&meetups).Where("id in (?)", pg.In(ids)).Select()

	return meetups, err

}

func (m *MeetupRepo) GetByKey(key, value string) (*models.Meetup, error) {
	meetup := new(models.Meetup)

	condition := fmt.Sprintf("%v = '%v'", key, value)

	err := m.DB.Model(meetup).Where(condition).First()

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, utils.ErrNoResult
		}
	}

	return meetup, err
}

func NewMeetupRepo(DB *pg.DB) *MeetupRepo {
	return &MeetupRepo{DB: DB}
}
