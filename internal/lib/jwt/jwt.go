package jwt

import (
	"BookingService/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user *entity.User, secret []byte, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["email"] = user.Contact.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
