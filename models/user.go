package models

type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"`

	//CreatedAt time.Time   `json:"createdAt"`
	//UpdatedAt time.Time   `json:"updatedAt"`
	//DeletedAt pg.NullTime `json:"deletedAt" pg:",soft_delete"`
}

