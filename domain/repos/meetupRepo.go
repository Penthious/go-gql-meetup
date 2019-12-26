package repos

import "github.com/penthious/go-gql-meetup/models"

type MeetupRepo interface {
	GetByIDs(id []string) ([]*models.Meetup, error)
	Create(meetup *models.Meetup) (*models.Meetup, error)
	All() ([]*models.Meetup, error)
	GetByKey(key, value string) (*models.Meetup, error)
	Update(meetup *models.Meetup) (*models.Meetup, error)
	Delete(meetup *models.Meetup) error
	GetMeetupsForUser(id string) ([]*models.Meetup, error)
	GetMeetupsByFilter(filter *models.MeetupFilter) ([]*models.Meetup, error)
}
