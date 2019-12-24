package middleware
//
//import (
//	"context"
//	"github.com/penthious/go-gql-meetup/graphql"
//	"github.com/penthious/go-gql-meetup/models"
//	"net/http"
//)
//
//// A private key for context that only this package can access. This is important
//// to prevent collisions between different context uses
//var userCtxKey = &contextKey{"user"}
//type contextKey struct {
//	name string
//}
//
//// Middleware decodes the share session cookie and packs the session into context
//func Middleware(d *graphql.Domain) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			c, err := r.Cookie("auth-cookie")
//
//			// Allow unauthenticated users in
//			if err != nil || c == nil {
//				next.ServeHTTP(w, r)
//				return
//			}
//
//			userId, err := d.
//			if err != nil {
//				http.Error(w, "Invalid cookie", http.StatusForbidden)
//				return
//			}
//
//			// get the user from the database
//			user := getUserByID(db, userId)
//
//			// put it in context
//			ctx := context.WithValue(r.Context(), userCtxKey, user)
//
//			// and call the next with our new context
//			r = r.WithContext(ctx)
//			next.ServeHTTP(w, r)
//		})
//	}
//}
//
//// ForContext finds the user from the context. REQUIRES Middleware to have run.
//func ForContext(ctx context.Context) *models.User {
//	raw, _ := ctx.Value(userCtxKey).(*models.User)
//	return raw
//}