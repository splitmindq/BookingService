package storage

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	GetPool() *pgxpool.Pool

	Close()
}
