package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/shirolimit/wallet-service/pkg/entities"
)

// Middleware is a service middleware type
type Middleware func(WalletService) WalletService

type loggingMiddleware struct {
	logger log.Logger
	next   WalletService
}

// LoggingMiddleware is a function that takes logger and produces
// a service Middleware used for logging
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next WalletService) WalletService {
		return &loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

// CreateAccount is a middleware function that prints information to log
// Named return parameter is used for defer
func (lmw loggingMiddleware) CreateAccount(ctx context.Context, acc entities.Account) (err error) {
	defer func(start time.Time) {
		lmw.logger.Log(
			"method", "CreateAccount",
			"account", acc,
			"error", err,
			"duration", time.Since(start),
		)
	}(time.Now())
	return lmw.next.CreateAccount(ctx, acc)
}

// ListAccounts is a middleware function that prints information to log
// Named return parameters are used for defer
func (lmw loggingMiddleware) ListAccounts(ctx context.Context) (accs []entities.AccountID, err error) {
	defer func(start time.Time) {
		lmw.logger.Log(
			"method", "ListAccounts",
			"accounts", accs,
			"error", err,
			"duration", time.Since(start),
		)
	}(time.Now())

	return lmw.next.ListAccounts(ctx)
}

// GetAccount is a middleware function that prints information to log
// Named return parameters are used for defer
func (lmw loggingMiddleware) GetAccount(ctx context.Context, id entities.AccountID) (acc entities.Account, err error) {
	defer func(start time.Time) {
		lmw.logger.Log(
			"method", "GetAccount",
			"id", id,
			"account", acc,
			"error", err,
			"duration", time.Since(start),
		)
	}(time.Now())

	return lmw.next.GetAccount(ctx, id)
}

// GetPayments is a middleware function that prints information to log
// Named return parameters are used for defer
func (lmw loggingMiddleware) GetPayments(ctx context.Context, id entities.AccountID) (payments []entities.Payment, err error) {
	defer func(start time.Time) {
		lmw.logger.Log(
			"method", "GetPayments",
			"id", id,
			"payments", payments,
			"error", err,
			"duration", time.Since(start),
		)
	}(time.Now())

	return lmw.next.GetPayments(ctx, id)
}

// MakePayment is a middleware function that prints information to log
// Named return parameter ise used for defer
func (lmw loggingMiddleware) MakePayment(ctx context.Context, payment entities.Payment) (err error) {
	defer func(start time.Time) {
		lmw.logger.Log(
			"method", "MakePayment",
			"payment", payment,
			"error", err,
			"duration", time.Since(start),
		)
	}(time.Now())

	return lmw.next.MakePayment(ctx, payment)
}
