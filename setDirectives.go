package go_gql_meetup
//
import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/penthious/go-gql-meetup/graphql/resolvers"
	"github.com/vektah/gqlparser/gqlerror"
	"reflect"
)

func SetDirectives(c *resolvers.Config)  {
	length(c)
}

func length(c *resolvers.Config) {
	c.Directives.Length = func(ctx context.Context, obj interface{}, next graphql.Resolver, min *int, max *int) (interface{}, error) {
		// @todo: https://github.com/99designs/gqlgen/issues/887 figure out how to get path set correctly in error
		v, _ := next(ctx)
		if reflect.TypeOf(v).String() == "string" {
			// Creates
			if len(v.(string)) < *min {
				return nil, gqlerror.Errorf("format")
			}

		} else if reflect.TypeOf(v).String() == "*string" {
			// Updates
			if len(*v.(*string)) < *min {
				return nil, gqlerror.Errorf("format")
			}
		}

		return v, nil
	}
}