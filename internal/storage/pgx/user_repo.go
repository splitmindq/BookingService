package pgx

import (
	"BookingService/internal/entity"
	"BookingService/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
	"log/slog"
)

var _ repository.UserRepository = (*UserRepo)(nil)

type UserRepo struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewUserRepo(db *pgxpool.Pool, log *slog.Logger) *UserRepo {
	return &UserRepo{db: db, log: log}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `INSERT INTO booking_service.users (phone, email, password_hash, first_name, last_name)
              VALUES ($1, $2, $3, $4, $5)
              RETURNING user_id, created_at, updated_at`

	log.Info("INFO OF USER BEFORE EXECUTING QUERY", user)

	err := r.db.QueryRow(ctx, query,
		user.Contact.Phone,
		user.Contact.Email,
		user.PasswordHash,
		user.Name,
		user.Surname,
	).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_email_key" {
				return nil, fmt.Errorf("%w: email already exists", ErrUserExists)
			}
			if pgErr.ConstraintName == "users_phone_key" {
				return nil, fmt.Errorf("%w: phone already exists", ErrUserExists)
			}
		}
		return nil, fmt.Errorf("%w: %v", ErrUserCreate, err)
	}

	user.Role = "user"

	log.Info("INFO OF USER AFTER EXECUTING QUERY", user)

	return user, nil
}
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	query := `SELECT u.user_id, u.email, u.password_hash, u.first_name, u.last_name, u.phone, u.created_at, u.updated_at, ur.role
              FROM booking_service.users u
              LEFT JOIN booking_service.user_roles ur ON u.user_id = ur.user_id
              WHERE u.email = $1`

	var role sql.NullString
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Contact.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Surname,
		&user.Contact.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		r.log.Error("failed to find user by email", "error", err)
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	// Устанавливаем роль "user" по умолчанию, если role NULL
	user.Role = role.String
	if !role.Valid {
		user.Role = "user"
	}

	r.log.Info("Found user",
		slog.Int64("user_id", user.Id),
		slog.String("email", user.Contact.Email),
		slog.String("role", user.Role))

	return &user, nil
}
func (r *UserRepo) IsAdmin(ctx context.Context, userId int64) (bool, error) {

	query := `SELECT EXISTS (
					    SELECT 1
						from booking_service.user_roles
						WHERE user_id = $1 AND role = 'admin')`

	err := r.db.QueryRow(ctx, query, userId).Scan(userId)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, errors.New("user is not an admin")
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if user is an admin: %w", err)
	}

	return userId > 0, nil

}

//
//func (r *UserRepo) GetRole(ctx context.Context, userId int64) (string, error) {
//	query := `SELECT role FROM booking_service.user_roles WHERE user_id = $1`
//
//	var role string
//
//	err := r.db.QueryRow(ctx, query, userId).Scan(&role)
//	if err != nil {
//		return "", fmt.Errorf("failed to find user by id: %w", err)
//	}
//
//	return role, nil
//
//}
