package postgres

import (
	"context"
	"log"

	"github.com/BillyBones007/loyalty-service/internal/db"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool            *pgxpool.Pool
	ConfigDB        *pgxpool.Config
	UserRepository  *UserRepository
	OrderRepository *OrderRepository
}

// Return struct with connections pool and config Postgres
func InitStorage(dst string) *Storage {
	config, err := pgxpool.ParseConfig(dst)
	if err != nil {
		log.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	store := Storage{Pool: pool, ConfigDB: config}
	store.CreateTables()
	return &store
}

// Migrations. Creates tables in the database if they do not exist.
func (s *Storage) CreateTables() error {
	m, err := migrate.New("file://migrations/postgres", s.ConfigDB.ConnString())
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}
	return nil
}

// Close the connections pool
func (s *Storage) Close() {
	s.Pool.Close()
}

// Get access to UserReposytory
func (s *Storage) User() db.UniversalUserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}
	s.UserRepository = &UserRepository{store: s}
	return s.UserRepository
}

// Get access to OrderRepository
func (s *Storage) Order() db.UniversalOrderRepository {
	if s.OrderRepository != nil {
		return s.OrderRepository
	}
	s.OrderRepository = &OrderRepository{store: s}
	return s.OrderRepository
}
