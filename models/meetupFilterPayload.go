package models

type MeetupFilterPayload struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}