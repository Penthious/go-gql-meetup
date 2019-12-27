package models

type Meetup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	
	UserID string `json:"user_id"`
	//CreatedAt time.Time   `json:"createdAt"`
	//UpdatedAt time.Time   `json:"updatedAt"`
	//DeletedAt pg.NullTime `json:"deletedAt" pg:",soft_delete"`
}

