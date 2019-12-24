package resolvers

import (
	"context"
	"github.com/penthious/go-gql-meetup/graphql/dataloaders"
	"github.com/penthious/go-gql-meetup/utils"

	//"github.com/penthious/go-gql-meetup/graphql/dataloaders"
	"github.com/penthious/go-gql-meetup/models"
)

func (r *queryResolver) Meetups(ctx context.Context) ([]*models.Meetup, error) {
	return r.Domain.DB.MeetupRepo.All()
}

type meetupResolver  struct{ *Resolver }

func (r *Resolver) Meetup() MeetupResolver {
	return &meetupResolver{r}
}
func (m *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {
	//return m.UsersRepo.GetByID(obj.UserID)
	return dataloaders.GetUserLoader(ctx).Load(obj.UserID)
}
//
//func  CreateMeetupIsValid(meetup models.NewMeetup) (bool, map[string]string) {
//	v := utils.NewValidator()
//
//	if meetup.Name != "" {
//		v.MustBeLongerThan("name", meetup.Name, 3)
//	}
//
//	return v.IsValid(), v.Errors
//}
func (m *mutationResolver) CreateMeetup(ctx context.Context, input models.NewMeetup) (*models.Meetup, error) {

	// @todo add validators in
	//if ok, err := CreateMeetupIsValid(input); ok != true {
	//	return nil,
	//}

	meetup := &models.Meetup{
		Name:        input.Name,
		Description: input.Description,
		UserID:      "1", //@todo refactor with jwt
	}

	return m.Domain.DB.MeetupRepo.Create(meetup)

}

func (m *mutationResolver) UpdateMeetup(ctx context.Context, id string, input models.UpdateMeetup) (*models.Meetup, error) {
	didUpdate := false
	meetup, err := m.DB.MeetupRepo.GetByKey("id", id)
	if err != nil || meetup == nil {
		return nil, utils.ErrNoResult
	}

	if *input.Name != ""{
		meetup.Name = *input.Name
		didUpdate = true
	}
	if *input.Description != ""{
		meetup.Description = *input.Description
		didUpdate = true
	}

	if didUpdate {
		meetup, err = m.DB.MeetupRepo.Update(meetup)
	} else {
		return nil, utils.ErrDidNotUpdate
	}

	if err != nil {
		return nil, utils.ErrUpdateError{Err: err}
	}

	return meetup, nil
}

func (m *mutationResolver) DeleteMeetup(ctx context.Context, id string) (bool, error) {

	meetup, err := m.DB.MeetupRepo.GetByKey("id", id)
	err = m.DB.MeetupRepo.Delete(meetup)

	if err != nil {
		return false, utils.ErrNoResult
	}

	return true, nil
}
