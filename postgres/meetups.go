package postgres

import (
	"github.com/go-pg/pg/v9"
	"github.com/penthious/go-gql-meetup/models"
)

type MeetupRepo struct {
	DB *pg.DB
}

func (m *MeetupRepo) GetMeetups() ([]*models.Meetup, error)  {
	var meetups []*models.Meetup

	err := m.DB.Model(&meetups).Select()

	if err != nil {
		return nil, err
	}

	return meetups, nil
	
}

func (m *MeetupRepo) GetByIDs(ids []string) ([]*models.Meetup, error)  {
	var meetups []*models.Meetup

	err  := m.DB.Model(meetups).Where("id = (?)", pg.In(ids)).Select()

	if err != nil {
		return nil, err
	}

	return meetups, nil

}

func (m *MeetupRepo) CreateMeetup(meetup *models.Meetup) (*models.Meetup, error) {
	_, err := m.DB.Model(meetup).Returning("*").Insert(meetup)

	if err != nil {
		return nil, err
	}

	return meetup, err
}
func NewMeetupRepo(DB *pg.DB) *MeetupRepo {
	return &MeetupRepo{DB: DB}
}
