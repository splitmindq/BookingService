package pgx

import (
	"BookingService/internal/config"
	"BookingService/internal/storage"
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

var _ storage.Storage = (*Storage)(nil)

//go:embed sql/tables.sql
var tableSchema string

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	
	poolConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	poolConfig.ConnConfig.User = cfg.DBConfig.Username
	poolConfig.ConnConfig.Password = cfg.DBConfig.Password
	poolConfig.ConnConfig.Host = cfg.DBConfig.Host
	poolConfig.ConnConfig.Port = uint16(cfg.DBConfig.Port)
	poolConfig.ConnConfig.Database = cfg.DBConfig.Name
	poolConfig.ConnConfig.TLSConfig = nil

	poolConfig.MaxConns = cfg.DBConfig.MaxConnections
	poolConfig.MinConns = cfg.DBConfig.MinConnections
	poolConfig.MaxConnLifetime = cfg.DBConfig.MaxConnectionLifetime
	poolConfig.MaxConnIdleTime = cfg.DBConfig.MaxConnectionIdleTime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	if err := createTable(ctx, db); err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}
	return &Storage{db: db}, nil

}

func (s *Storage) GetPool() *pgxpool.Pool {
	return s.db
}

func createTable(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, tableSchema)
	if err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}
	log.Println("Database schema initialized successfully")
	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}
