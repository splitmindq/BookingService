package service

import (
	"BookingService/internal/entity"
	"BookingService/internal/storage/pgx"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *pgx.UserRepo
}

func NewAuthService(userRepo *pgx.UserRepo) *AuthService {
	return &AuthService{userRepo}
}

func (s *AuthService) SignUp(ctx context.Context, input entity.SignUpInput) (*entity.User, error) {

	existingUser, err := s.userRepo.FindByEmail(ctx, input.Contact.Email)
	if err == nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
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

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) SignIn(ctx context.Context, input entity.SignInInput) (*entity.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
