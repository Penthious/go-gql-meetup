package validators

import (
	"github.com/penthious/go-gql-meetup/domain/utils"
	"github.com/penthious/go-gql-meetup/models"
)

func  CreateMeetupValidator(meetup models.NewMeetup) (bool, map[string]string) {
	v := utils.NewValidator()

	if meetup.Name != "" {
		v.MustBeLongerThan("name", meetup.Name, 3)
	}

	return v.IsValid(), v.Errors
}
