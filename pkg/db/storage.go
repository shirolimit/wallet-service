package db

import (
	"context"

	"github.com/shirolimit/wallet-service/pkg/entities"
)

// Storage is an interface to some persistent storage, for example relational database
type Storage interface {
	CreateAccount(context.Context, entities.Account) error
	GetAccount(context.Context, entities.AccountID) (*entities.Account, error)
	ListAccounts(context.Context) ([]entities.AccountID, error)

	PaymentsByAccount(context.Context, entities.AccountID) ([]entities.Payment, error)
	CreatePayment(context.Context, entities.Payment) error
}
