//go:generate go run github.com/99designs/gqlgen -v
package go_gql_meetup

import (
	"context"
	"github.com/penthious/go-gql-meetup/models"
)

type Resolver struct{}

type queryResolver struct{ *Resolver }
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *queryResolver) Meetups(ctx context.Context) ([]*models.Meetup, error) {
	panic("not implemented")
}


type meetupResolver  struct{ *Resolver }
func (r *Resolver) Meetup() MeetupResolver {
	return &meetupResolver{r}
}
func (m *meetupResolver) User(ctx context.Context, obj *models.Meetup) (*models.User, error) {
	panic("implement me")
}

type userResolver  struct{ *Resolver }
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
func (u *userResolver) Meetups(ctx context.Context, obj *models.User) ([]*models.Meetup, error) {
	panic("implement me")
}

type mutationResolver struct { *Resolver }
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (g *mutationResolver) CreateMeetup(ctx context.Context, input NewMeetup) (*models.Meetup, error) {
	panic("implement me")
}



