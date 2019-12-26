package resolvers

import (
	"context"
	"fmt"
	"github.com/penthious/go-gql-meetup/domain/utils"
	"github.com/penthious/go-gql-meetup/graphql/dataloaders"

	//"github.com/penthious/go-gql-meetup/graphql/dataloaders"
	"github.com/penthious/go-gql-meetup/models"
)

func (q *queryResolver) Meetups(ctx context.Context, filter *models.MeetupFilter) ([]*models.Meetup, error) {
	return q.Domain.DB.MeetupRepo.GetMeetupsByFilter(filter)
}


type meetupResolver  struct{ *Resolver }

func (r *Resolver) Meetup() MeetupResolver {
	return &meetupResolver{r}
}
func (m *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {
	return dataloaders.GetUserLoader(ctx).Load(obj.UserID)
}

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

	fmt.Println("WE GOT HERE")
	if input.Name != nil {
		if *input.Name != meetup.Name {
			meetup.Name = *input.Name
			didUpdate = true
		}
	}
	fmt.Println("WE geot HERE")
	if input.Description != nil {
		if *input.Description != meetup.Description {
			meetup.Description = *input.Description
			didUpdate = true
		}
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

func (m *mutationResolver) FilterMeetups(ctx context.Context, filter *models.MeetupFilter) {

}
