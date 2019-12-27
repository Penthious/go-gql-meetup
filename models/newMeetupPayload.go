package models

type NewMeetupPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
