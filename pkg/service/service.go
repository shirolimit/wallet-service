package service

import "github.com/shirolimit/wallet-service/pkg/entities"

// WalletService is the main service interface
type WalletService interface {
	CreateAccount(entities.Account) error
	ListAccounts() ([]string, error)
	GetAccount(entities.AccountID) (entities.Account, error)

	GetPayments(entities.AccountID) ([]entities.Payment, error)
	MakePayment(entities.Payment) error
}

type walletService struct {
}

func NewWalletService() WalletService {
	return &walletService{}
}

func (ws *walletService) CreateAccount(acc entities.Account) error {
	return nil
}
