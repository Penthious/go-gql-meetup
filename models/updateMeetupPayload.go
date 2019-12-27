package models

type UpdateMeetupPayload struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

