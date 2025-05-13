package service

import (
	"BookingService/internal/config"
	"BookingService/internal/entity"
	"BookingService/internal/lib/jwt"
	"BookingService/internal/repository"
	"BookingService/internal/storage/pgx"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *pgx.UserRepo, cfg *config.Config) *AuthService {
	return &AuthService{userRepo, cfg}
}

func (s *AuthService) SignUp(ctx context.Context, input entity.SignUpInput) (int64, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &entity.User{
		Contact: entity.Contact{
			Email: input.Contact.Email,
			Phone: input.Contact.Phone,
		},
		PasswordHash: string(hashedPassword),
		Name:         input.FirstName,
		Surname:      input.LastName,
	}

	user, err = s.userRepo.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return user.Id, nil
}

func (s *AuthService) SignIn(ctx context.Context, input entity.SignInInput) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := jwt.NewToken(user, s.cfg.JwtSecret, s.cfg.JwtExpire)

	return token, nil
}

//func (s *AuthService) IsAdmin(ctx context.Context, userId int64) error {
//
//}
