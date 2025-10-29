package service

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Wallet interface {
	Deposit(uuid uuid.UUID, amount float64) (float64, error)
	Withdraw(uuid uuid.UUID, amount float64) (float64, error)
	Balance(uuid uuid.UUID) (float64, error)
}

type wallet struct {
	storage *pgxpool.Pool
}

func New(s *pgxpool.Pool) Wallet {
	return &wallet{
		storage: s,
	}
}

func (ws *wallet) Deposit(uuid uuid.UUID, amount float64) (float64, error) {
	var res float64
	return res, nil
}
func (ws *wallet) Withdraw(uuid uuid.UUID, amount float64) (float64, error) {
	var res float64
	return res, nil
}
func (ws *wallet) Balance(uuid uuid.UUID) (float64, error) {
	var res float64
	return res, nil
}
