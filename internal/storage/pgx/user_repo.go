package pgx

import (
	"BookingService/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {

	query := `INSERT INTO booking_service.users (phone,email,password_hash,first_name,last_name,created_at)
						VALUES ($1,$2,$3,$4,$5,$6) RETURNING user_id,created_at,updated_at`

	err := r.db.QueryRow(ctx, query,
		user.Contact.Phone,
		user.Contact.Email,
		user.PasswordHash,
		user.Name,
		user.Surname,
		user.CreatedAt).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && (pgErr.ConstraintName == "users_email_key" || pgErr.ConstraintName == "users_phone_key") {
				return fmt.Errorf("%w: %v", ErrUserExists, pgErr)
			}
		}
		return fmt.Errorf("%w: %v", ErrUserCreate, err)
	}

	return nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `SELECT * FROM booking_service.users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Contact.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Surname,
		&user.Contact.Phone,
		&user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &user, err

}
