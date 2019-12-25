package utils
//
//import (
//	"context"
//	"fmt"
//	"github.com/99designs/gqlgen/graphql"
//	"github.com/penthious/go-gql-meetup/graphql/resolvers"
//)
//
//func SetDirectives(c *resolvers.Config)  {
//	c.Directives.Length = func(ctx context.Context, obj interface{}, next graphql.Resolver, min *int, max *int) (interface{}, error) {
//		fmt.Println("\n====================\n")
//		fmt.Println(min)
//		fmt.Println("\n====================\n")
//		//if  {
//		fmt.Println(max)
//		fmt.Println("\n====================\n")
//		fmt.Println(obj)
//		fmt.Println("\n====================\n")
//
//		//
//		//	// block calling the next resolver
//		//	return nil, fmt.Errorf("Access denied")
//		//}
//
//		// or let it pass through
//		return next(ctx)
//	}
//}