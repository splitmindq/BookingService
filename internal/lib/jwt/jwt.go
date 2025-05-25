package jwt

import (
	"BookingService/internal/entity"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user *entity.User, secret string, duration time.Duration) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"role":    user.Role,
		"exp":     time.Now().Add(duration).Unix(),
	})

	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseTokenAndGetUID(tokenString string, secret string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid, ok := claims["user_id"].(float64)
		if !ok {
			return 0, "", fmt.Errorf("invalid user_id in token claims")
		}
		role, _ := claims["role"].(string)
		if role == "" {
			role = "user"
		}
		return int64(uid), role, nil
	}
	return 0, "", fmt.Errorf("invalid token claims")
}
