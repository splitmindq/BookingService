package jwt

import (
	"BookingService/internal/entity"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user *entity.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["email"] = user.Contact.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseTokenAndGetUID(tokenString string, secret string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if uid, ok := claims["user_id"]; ok {
			return int64(uid.(float64)), nil
		}
	}
	return 0, fmt.Errorf("invalid token claims")
}
