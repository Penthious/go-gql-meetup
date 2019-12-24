//go:generate go run github.com/99designs/gqlgen -v
package resolvers

import (
	"github.com/penthious/go-gql-meetup/graphql"
)

type Resolver struct{graphql.Domain}

type queryResolver struct{ *Resolver }


func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct { *Resolver }

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
