package service

import (
	"context"

	"github.com/shirolimit/wallet-service/pkg/entities"
)

// WalletService is the main service interface
type WalletService interface {
	CreateAccount(ctx context.Context, account entities.Account) error
	ListAccounts(ctx context.Context) ([]entities.AccountID, error)
	GetAccount(ctx context.Context, id entities.AccountID) (entities.Account, error)

	GetPayments(ctx context.Context, id entities.AccountID) ([]entities.Payment, error)
	MakePayment(ctx context.Context, payment entities.Payment) error
}

type walletService struct {
}

// NewWalletService creates new instance of walletService
func NewWalletService() WalletService {
	return &walletService{}
}

func (ws *walletService) CreateAccount(ctx context.Context, acc entities.Account) error {
	return entities.ErrAccountAlreadyExists
}

func (ws *walletService) ListAccounts(ctx context.Context) ([]entities.AccountID, error) {
	return nil, nil
}

func (ws *walletService) GetAccount(ctx context.Context, id entities.AccountID) (entities.Account, error) {
	return entities.Account{}, nil
}

func (ws *walletService) GetPayments(ctx context.Context, id entities.AccountID) ([]entities.Payment, error) {
	return nil, nil
}

func (ws *walletService) MakePayment(ctx context.Context, payment entities.Payment) error {
	return nil
}
