package dataloaders

import (
	"github.com/penthious/go-gql-meetup/graphql"
	"github.com/penthious/go-gql-meetup/models"
	"time"
)

func CreateUserLoader(d *graphql.Domain) UserLoader {
	return UserLoader{
		maxBatch: 100,
		wait:     1 * time.Millisecond,
		fetch: func(ids []string) ([]*models.User, []error) {
			users, err := d.DB.UserRepo.GetByIDs(ids)

			if err != nil {
				return nil, []error{err}
			}

			u := make(map[string]*models.User, len(users))
			for _, user := range users {
				u[user.ID] = user
			}

			result := make([]*models.User, len(ids))

			for i, id := range ids {
				result[i] = u[id]
			}

			return result, []error{err}
		},
	}
}
