package resolvers

import (
	"context"
	"errors"
	"github.com/penthious/go-gql-meetup/graphql/dataloaders"

	//"github.com/penthious/go-gql-meetup/graphql/dataloaders"
	"github.com/penthious/go-gql-meetup/models"
)

func (r *queryResolver) Meetups(ctx context.Context) ([]*models.Meetup, error) {
	return r.MeetupsRepo.GetMeetups()
}

type meetupResolver  struct{ *Resolver }

func (r *Resolver) Meetup() MeetupResolver {
	return &meetupResolver{r}
}
func (m *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {
	//return m.UsersRepo.GetByID(obj.UserID)
	return dataloaders.GetUserLoader(ctx).Load(obj.UserID)
}

func (m *mutationResolver) CreateMeetup(ctx context.Context, input models.NewMeetup) (*models.Meetup, error) {
	if len(input.Name) < 3 {
		return nil, errors.New("name not long enough")
	}
	if len(input.Description) < 0 {
		return nil, errors.New("description is required")
	}

	meetup := &models.Meetup{
		Name:        input.Name,
		Description: input.Description,
		UserID:      "1", //@todo refactor with jwt
	}

	return m.MeetupsRepo.CreateMeetup(meetup)
}
