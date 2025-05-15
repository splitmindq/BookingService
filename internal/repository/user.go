package repository

import (
	"BookingService/internal/entity"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}
