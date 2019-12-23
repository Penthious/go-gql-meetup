package dataloaders

import (
	"context"
	"github.com/penthious/go-gql-meetup/graphql"
	"net/http"
)

const USER_LOADER_KEY = "userLoader"

func DataloaderMiddleware(d *graphql.Domain, next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLoader := CreateUserLoader(d)

		ctx := context.WithValue(r.Context(), USER_LOADER_KEY, &userLoader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func GetUserLoader(ctx context.Context) *UserLoader {
	return ctx.Value(USER_LOADER_KEY).(*UserLoader )
}
