package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/lib/pq"
	"github.com/shirolimit/wallet-service/pkg/entities"
)

// pgStorage is a Storage implementation that uses Postgres
type pgStorage struct {
	db *sql.DB
}

// pgAccount is a helper struct for working with Account entity
// it contains serial id field
type pgAccount struct {
	account    entities.Account
	internalID int64
}

// balanceUpdateHelper is a tiny struct to simplify balance updates
type balanceUpdateHelper struct {
	internalAccountID int64
	diff              decimal.Decimal
}

type getPaymentsHelper struct {
	id            uuid.UUID
	source        entities.AccountID
	destination   entities.AccountID
	sourceID      int64
	destinationID int64
	amount        decimal.Decimal
}

// NewPgStorage creates new Postgres storage with specified connection string
func NewPgStorage(connectionString string) Storage {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	return &pgStorage{
		db: db,
	}
}

// PgStorageFromHandle creates new Postgres storage with specified sql.DB instance
// It is mostly used in tests
func PgStorageFromHandle(db *sql.DB) Storage {
	return &pgStorage{
		db: db,
	}
}

func (ps *pgStorage) CreateAccount(ctx context.Context, acc entities.Account) error {
	_, err := ps.db.ExecContext(
		ctx,
		"insert into accounts (account_id, currency, balance) values ($1, $2, $3);",
		acc.ID, acc.Currency, acc.Balance,
	)

	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if !ok {
			return err
		}

		if pgErr.Code == pq.ErrorCode("23505") {
			return entities.ErrAccountAlreadyExists
		}
		return err
	}
	return nil
}

func (ps *pgStorage) GetAccount(ctx context.Context, id entities.AccountID) (*entities.Account, error) {
	pgAcc, err := ps.selectAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pgAcc.account, nil
}

func (ps *pgStorage) ListAccounts(ctx context.Context) ([]entities.AccountID, error) {
	rows, err := ps.db.QueryContext(ctx, "select account_id from accounts;")
	if err != nil {
		return nil, err
	}

	accounts := make([]entities.AccountID, 0)
	for rows.Next() {
		var account string
		err = rows.Scan(&account)
		if err != nil {
			break
		}
		accounts = append(accounts, entities.AccountID(account))
	}
	return accounts, nil
}

func (ps *pgStorage) PaymentsByAccount(ctx context.Context, id entities.AccountID) ([]entities.Payment, error) {
	pgAcc, err := ps.selectAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	rows, err := ps.db.QueryContext(
		ctx,
		`select p.id, p.source_id, p.destination_id, a1.account_id as source, a2.account_id as destination, p.amount 
		from payments as p
			join accounts as a1 on source_id = a1.id
			join accounts as a2 on destination_id = a2.id
		where p.source_id = $1 or p.destination_id = $1`,
		pgAcc.internalID,
	)
	if err != nil {
		return nil, err
	}

	payments := make([]entities.Payment, 0)
	for rows.Next() {
		var helper getPaymentsHelper
		err = rows.Scan(&helper.id, &helper.sourceID, &helper.destinationID,
			&helper.source, &helper.destination, &helper.amount)

		if err != nil {
			break
		}

		payment := entities.Payment{
			Account: pgAcc.account.ID,
			Amount:  helper.amount,
			ID:      helper.id,
		}
		if helper.sourceID == pgAcc.internalID {
			payment.Direction = entities.Outgoing
			payment.ToAccount = &helper.destination
		} else {
			payment.Direction = entities.Incoming
			payment.FromAccount = &helper.source
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (ps *pgStorage) CreatePayment(ctx context.Context, payment entities.Payment) error {
	var sourceAccount *pgAccount
	var destinationAccount *pgAccount
	var err error

	if payment.Direction == entities.Outgoing {
		sourceAccount, err = ps.selectAccount(ctx, payment.Account)
		destinationAccount, err = ps.selectAccount(ctx, *payment.ToAccount)
	} else {
		sourceAccount, err = ps.selectAccount(ctx, *payment.FromAccount)
		destinationAccount, err = ps.selectAccount(ctx, payment.Account)
	}

	if sourceAccount == nil {
		return entities.ErrPaymentSourceNotFound
	}

	if destinationAccount == nil {
		return entities.ErrPaymentDestinationNotFound
	}

	if sourceAccount.account.Currency != destinationAccount.account.Currency {
		return entities.ErrDifferentCurrencies
	}

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// insert payment
	_, err = tx.Exec(
		"insert into payments (id, source_id, destination_id, amount) values ($1, $2, $3, $4);",
		payment.ID,
		sourceAccount.internalID,
		destinationAccount.internalID,
		payment.Amount,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// explicitly close query to avoid errors
	//r.Close()

	// update balances
	updates := []balanceUpdateHelper{
		{internalAccountID: sourceAccount.internalID, diff: payment.Amount.Neg()},
		{internalAccountID: destinationAccount.internalID, diff: payment.Amount},
	}

	// make sure that we update accounts in the same order to avoid deadlocks
	if sourceAccount.internalID > destinationAccount.internalID {
		updates[0], updates[1] = updates[1], updates[0]
	}

	for _, u := range updates {
		_, err = tx.Exec(
			"update accounts set balance = balance + $1 where id = $2;",
			u.diff,
			u.internalAccountID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}

		// explicitly close query to avoid errors
		//r.Close()
	}

	return tx.Commit()
}

func (ps *pgStorage) selectAccount(ctx context.Context, id entities.AccountID) (*pgAccount, error) {
	var acc pgAccount
	err := ps.db.QueryRowContext(
		ctx,
		"select id, account_id, currency, balance from accounts where account_id = $1;",
		id,
	).Scan(&acc.internalID, &acc.account.ID, &acc.account.Currency, &acc.account.Balance)

	if err != nil {
		return nil, err
	}

	return &acc, nil
}
