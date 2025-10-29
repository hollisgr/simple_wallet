package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Create() uuid.UUID
	Balance(uuid uuid.UUID) (float64, error)
	Deposit(uuid uuid.UUID, amount float64) (float64, error)
	Withdraw(uuid uuid.UUID, amount float64) (float64, error)
}

type storage struct {
	db *pgxpool.Pool
}

func New(p *pgxpool.Pool) Storage {
	return &storage{
		db: p,
	}
}

func (s *storage) Create() uuid.UUID {
	uuid := uuid.New()
	return uuid
}

func (s *storage) Balance(uuid uuid.UUID) (float64, error) {
	var res float64
	return res, nil
}

func (s *storage) Deposit(uuid uuid.UUID, amount float64) (float64, error) {
	var res float64
	return res, nil
}

func (s *storage) Withdraw(uuid uuid.UUID, amount float64) (float64, error) {
	var res float64
	return res, nil
}
