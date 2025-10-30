package db

import (
	"context"
	"fmt"

	"cmd/app/main.go/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Create(ctx context.Context, uuid uuid.UUID) error
	Balance(ctx context.Context, uuid uuid.UUID) (model.Wallet, error)
	Deposit(ctx context.Context, uuid uuid.UUID, amount float64) (model.Wallet, error)
	Withdraw(ctx context.Context, uuid uuid.UUID, amount float64) (model.Wallet, error)
}

type storage struct {
	db *pgxpool.Pool
}

func New(p *pgxpool.Pool) Storage {
	return &storage{
		db: p,
	}
}

// Create inserts a new wallet record into the database and returns any encountered errors.
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

// Balance retrieves balance information from the database for a specified wallet UUID.
func (s *storage) Balance(ctx context.Context, uuid uuid.UUID) (model.Wallet, error) {
	var res model.Wallet
	query := `
		SELECT 
			uuid,
			balance
		FROM
			wallets
		WHERE
			uuid = @uuid
	`
	args := pgx.NamedArgs{
		"uuid": uuid,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Wallet])

	if err != nil {
		return res, err
	}

	return res, nil
}

// Deposit updates the balance of a wallet by adding a specified amount and returns updated wallet data.
func (s *storage) Deposit(ctx context.Context, uuid uuid.UUID, amount float64) (model.Wallet, error) {
	var res model.Wallet
	query := `
		UPDATE 
			wallets
		SET
			balance = balance + @amount
		WHERE
			uuid = @uuid
		RETURNING balance
	`
	args := pgx.NamedArgs{
		"uuid":   uuid,
		"amount": amount,
	}
	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&res.Balance)

	if err != nil {
		return res, err
	}

	res.UUID = uuid

	return res, nil
}

// Withdraw subtracts a specified amount from the wallet's balance and returns updated wallet data.
func (s *storage) Withdraw(ctx context.Context, uuid uuid.UUID, amount float64) (model.Wallet, error) {
	var res model.Wallet
	query := `
		UPDATE 
			wallets
		SET
			balance = balance - @amount
		WHERE
			uuid = @uuid
		RETURNING balance
	`
	args := pgx.NamedArgs{
		"uuid":   uuid,
		"amount": amount,
	}
	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&res.Balance)

	if err != nil {
		return res, err
	}

	res.UUID = uuid

	return res, nil
}
