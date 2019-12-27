package middleware

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/domain/sevices"
	"github.com/penthious/go-gql-meetup/models"
	"net/http"
	"os"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
type ContextKey string
var userCtxKey = ContextKey("user")
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type authResponseWriter struct {
	http.ResponseWriter
	userIDToResolver string
	userIDFromCookie string
	domain.Domain
}

func (w *authResponseWriter) Write(b []byte) (int, error) {
	if w.userIDToResolver != w.userIDFromCookie {
		token, _ := sevices.GenToken(w.userIDToResolver, w.Domain)

		http.SetCookie(w, &http.Cookie{
			Name:     "auth",
			Value:    token.AccessToken,
		})
	}
return w.ResponseWriter.Write(b)
}
// Middleware decodes the share session cookie and packs the session into context
func Auth(d *domain.Domain) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			arw := authResponseWriter{w, "", "", *d}
			userIDContextKey := ContextKey("userID")
			c, err := r.Cookie("auth")
			if err != nil || c == nil {
				// Allow unauthenticated users in to unauthenticated routes
				ctx := context.WithValue(r.Context(), userIDContextKey, &arw.userIDToResolver)
				r = r.WithContext(ctx)
				next.ServeHTTP(&arw, r)
				return
			}
			claims := &sevices.Claims{}
			jwtToken, err := jwt.ParseWithClaims(c.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !jwtToken.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			user, err := d.DB.UserRepo.GetByKey("id", claims.User.ID)
			if err != nil || user == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			arw.userIDFromCookie = user.ID
			arw.userIDToResolver = user.ID

			ctx := context.WithValue(r.Context(), userIDContextKey, &arw.userIDToResolver)
			ctx = context.WithValue(ctx, userCtxKey, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(&arw, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *models.User {
	raw, _ := ctx.Value(userCtxKey).(*models.User)
	return raw
}