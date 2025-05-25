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
	"log"
	"log/slog"
)

type AuthService struct {
	userRepo repository.UserRepository
	Cfg      *config.Config
	log      *slog.Logger
}

func NewAuthService(userRepo *pgx.UserRepo, cfg *config.Config, log *slog.Logger) *AuthService {
	return &AuthService{userRepo, cfg, log}
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

	s.log.Info("SignUp user ", user)

	user, err = s.userRepo.Create(ctx, user)

	s.log.Info("CREATED user ", user)

	if err != nil {
		log.Printf("Error creating user: %v", err)
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

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := jwt.NewToken(user, s.Cfg.HTTPServer.JwtSecret, s.Cfg.HTTPServer.JwtExpire)

	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}
