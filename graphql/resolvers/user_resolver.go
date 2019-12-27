package resolvers

import (
	"context"
	"github.com/penthious/go-gql-meetup/domain/sevices"
	"github.com/penthious/go-gql-meetup/domain/utils"
	"github.com/penthious/go-gql-meetup/models"
	"github.com/penthious/go-gql-meetup/server/middleware"
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

func (m *mutationResolver) Register(ctx context.Context, input models.RegisterPayload) (*models.User, error) {
	userExist, _ := m.DB.UserRepo.GetByKey("email", input.Email)

	if userExist != nil {
		return nil, utils.ErrUserWithEmailAlreadyExist
	}

	password, err := sevices.SetPassword(input.Password)

	if err != nil {
		return nil, err
	}

	data := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: *password,
	}

	user, err := m.DB.UserRepo.Create(data)

	if err != nil {
		return nil, err
	}

	authPointer := ctx.Value(middleware.ContextKey("userID")).(*string)
	*authPointer = user.ID

	return user, nil
}

func (m *mutationResolver) Login(ctx context.Context, input models.LoginPayload) (*models.User, error) {
	user, err := m.DB.UserRepo.GetByKey("email", input.Email)
	if user == nil || err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	err = sevices.CheckPassword(input.Password, user)

	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	authPointer := ctx.Value(middleware.ContextKey("userID")).(*string)
	*authPointer = user.ID

	return user, nil
}
