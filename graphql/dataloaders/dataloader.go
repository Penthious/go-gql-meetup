package dataloaders

import (
	"context"
	"github.com/penthious/go-gql-meetup/graphql"
	"github.com/penthious/go-gql-meetup/models"
	"net/http"
	"time"
)

const USER_LOADER_KEY = "userLoader"

func DataloaderMiddleware(d *graphql.Domain, next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLoader := UserLoader{
			maxBatch: 100,
			wait: 1 * time.Millisecond,
			fetch: func(ids []string) ([]*models.User, []error) {
				users, err := d.DB.UserRepo.GetByIDs(ids)

				return users, []error{err}
			},
		}

		ctx := context.WithValue(r.Context(), USER_LOADER_KEY, &userLoader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func GetUserLoader(ctx context.Context) *UserLoader {
	return ctx.Value(USER_LOADER_KEY).(*UserLoader )
}
