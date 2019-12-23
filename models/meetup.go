package models

type Meetup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	
	UserID string `json:"user_id"`
}

type MeetupRepo interface {
	GetByIDs(id []string) ([]*Meetup, error)
	//GetByID(id int64) (*models.User, error)
	//GetByEmail(email string) (*models.User, error)
	//GetByUsername(username string) (*models.User, error)
	//Create(user *models.User) (*models.User, error)
}
