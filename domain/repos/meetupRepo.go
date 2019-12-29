package repos

import "github.com/penthious/go-gql-meetup/models"

type MeetupRepo interface {
	All() ([]*models.Meetup, error)
	Create(meetup *models.Meetup) (*models.Meetup, error)
	Delete(meetup *models.Meetup) error
	GetByIDs(id []string) ([]*models.Meetup, error)
	GetByKey(key, value string) (*models.Meetup, error)
	GetMeetupsByFilter(filter *models.MeetupFilterPayload) ([]*models.Meetup, error)
	GetMeetupsForUser(id string) ([]*models.Meetup, error)
	Update(meetup *models.Meetup) (*models.Meetup, error)
}
