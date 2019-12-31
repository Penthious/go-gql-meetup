package sevices

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/models"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type JWTToken struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}
type Claims struct {
	User *models.User
	jwt.StandardClaims
}


func GenToken(userID string, d domain.Domain) (*JWTToken, error) {
	jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))
	expiresAt := time.Now().Add(time.Hour * 24 * 7) // 1 week
	user, _:= d.DB.UserRepo.GetByKey("id", userID)

	jwtToken.Claims = Claims{
		User:           user,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expiresAt.Unix()},
	}

	accessToken, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, err
	}

	return &JWTToken{AccessToken: accessToken, ExpiresAt: expiresAt}, nil
}
func CheckPassword(password string, u *models.User) error {
	passwordByte, passwordHashedByte := []byte(password), []byte(u.Password)
	return bcrypt.CompareHashAndPassword(passwordHashedByte, passwordByte)
}