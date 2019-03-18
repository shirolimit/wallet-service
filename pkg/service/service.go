package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shirolimit/wallet-service/pkg/db"
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
	storage db.Storage
}

var (
	nullUUID = uuid.UUID{}
)

// NewWalletService creates new instance of walletService
func NewWalletService(storage db.Storage) WalletService {
	return &walletService{
		storage: storage,
	}
}

func (ws *walletService) CreateAccount(ctx context.Context, acc entities.Account) error {
	if len(acc.ID) == 0 {
		return entities.ErrEmptyAccountID
	}

	if len(acc.Currency) == 0 {
		return entities.ErrEmptyAccountCurrency
	}

	if acc.Balance.IsNegative() {
		return entities.ErrNegativeBalance
	}

	return ws.storage.CreateAccount(ctx, acc)
}

func (ws *walletService) ListAccounts(ctx context.Context) ([]entities.AccountID, error) {
	accounts, err := ws.storage.ListAccounts(ctx)
	if accounts == nil {
		accounts = []entities.AccountID{}
	}
	return accounts, err
}

func (ws *walletService) GetAccount(ctx context.Context, id entities.AccountID) (entities.Account, error) {
	acc, err := ws.storage.GetAccount(ctx, id)
	if err != nil {
		return entities.Account{}, err
	}
	return *acc, nil
}

func (ws *walletService) GetPayments(ctx context.Context, id entities.AccountID) ([]entities.Payment, error) {
	payments, err := ws.storage.PaymentsByAccount(ctx, id)
	if payments == nil {
		payments = []entities.Payment{}
	}
	return payments, err
}

func (ws *walletService) MakePayment(ctx context.Context, payment entities.Payment) error {
	if payment.ID == nullUUID {
		return entities.ErrEmptyPaymentID
	}

	if len(payment.Account) == 0 {
		return entities.ErrEmptyAccountID
	}

	if payment.ToAccount == nil || len(*payment.ToAccount) == 0 {
		return entities.ErrEmptyPaymentDestination
	}

	if payment.Amount.IsNegative() || payment.Amount.IsZero() {
		return entities.ErrWrongPaymentAmount
	}

	if payment.Direction == entities.Incoming {
		return entities.ErrIncomingPaymentsNotAllowed
	}

	return ws.storage.CreatePayment(ctx, payment)
}
