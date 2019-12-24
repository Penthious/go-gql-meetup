package resolvers

import (
	"context"
	"github.com/penthious/go-gql-meetup/models"
)

type userResolver  struct{ *Resolver }
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
func (q *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	return q.Domain.DB.UserRepo.All()
}
func (q *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	return q.Domain.DB.UserRepo.GetByKey("id", id)
}

func (u *userResolver) Meetups(ctx context.Context, obj *models.User) ([]*models.Meetup, error) {
	return u.Domain.DB.MeetupRepo.GetMeetupsForUser(obj.ID)// @todo refactor to only get current users
}

