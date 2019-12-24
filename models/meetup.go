package models

type Meetup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	
	UserID string `json:"user_id"`
}

type MeetupRepo interface {
	GetByIDs(id []string) ([]*Meetup, error)
	Create(meetup *Meetup) (*Meetup, error)
	All() ([]*Meetup, error)
	GetByKey(key, value string) (*Meetup, error)
	Update(meetup *Meetup) (*Meetup, error)
	Delete(meetup *Meetup) error
	GetMeetupsForUser(id string) ([]*Meetup, error)
}
