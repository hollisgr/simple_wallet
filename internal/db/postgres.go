package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Create(ctx context.Context, uuid uuid.UUID) error
	Balance(ctx context.Context, uuid uuid.UUID) (float64, error)
	Deposit(ctx context.Context, uuid uuid.UUID, amount float64) (float64, error)
	Withdraw(ctx context.Context, uuid uuid.UUID, amount float64) (float64, error)
}

type storage struct {
	db *pgxpool.Pool
}

func New(p *pgxpool.Pool) Storage {
	return &storage{
		db: p,
	}
}

func (s *storage) Create(ctx context.Context, uuid uuid.UUID) error {
	query := `
		INSERT INTO
			wallets (uuid)
		VALUES
			(@uuid)
		RETURNING
			uuid
	`
	args := pgx.NamedArgs{
		"uuid": uuid,
	}
	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&uuid)
	if err != nil {
		return fmt.Errorf("db create wallet error: %v", err)
	}
	return nil
}

func (s *storage) Balance(ctx context.Context, uuid uuid.UUID) (float64, error) {
	var res float64
	return res, nil
}

func (s *storage) Deposit(ctx context.Context, uuid uuid.UUID, amount float64) (float64, error) {
	var res float64
	return res, nil
}

func (s *storage) Withdraw(ctx context.Context, uuid uuid.UUID, amount float64) (float64, error) {
	var res float64
	return res, nil
}
