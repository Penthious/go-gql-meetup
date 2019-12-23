package resolvers

import (
	"context"
	"github.com/penthious/go-gql-meetup/models"
)

type userResolver  struct{ *Resolver }
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
func (u *userResolver) Meetups(ctx context.Context, obj *models.User) ([]*models.Meetup, error) {
	return u.Domain.DB.MeetupRepo.All() // @todo refactor to only get current users
}
