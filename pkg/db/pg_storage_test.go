package db_test

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/shirolimit/wallet-service/pkg/entities"
	"github.com/shopspring/decimal"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	mydb "github.com/shirolimit/wallet-service/pkg/db"
)

func Test_PgStorage_ListAccounts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' while opening a mock database connection", err)
	}
	defer db.Close()

	expected := []entities.AccountID{"alice", "bob"}

	mock.ExpectQuery("select account_id from accounts;").
		WillReturnRows(sqlmock.NewRows([]string{"account_id"}).AddRow("alice").AddRow("bob"))

	storage := mydb.PgStorageFromHandle(db)
	accounts, storageErr := storage.ListAccounts(context.TODO())
	if storageErr != nil {
		t.Errorf("Error while quering accounts list: %v", err)
	}
	if !reflect.DeepEqual(accounts, expected) {
		t.Errorf("Expectation failed. Expected accounts = %v, actual = %v", expected, accounts)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func Test_PgStorage_ListAccountsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' while opening a mock database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("select account_id from accounts;").WillReturnError(sql.ErrConnDone)

	storage := mydb.PgStorageFromHandle(db)
	_, storageErr := storage.ListAccounts(context.TODO())
	if storageErr == nil {
		t.Errorf("Error expectation failed")
	}
}

func Test_PgStorage_CreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' while opening a mock database connection", err)
	}
	defer db.Close()

	acc := entities.Account{
		ID:       "alice",
		Balance:  decimal.New(100, 0),
		Currency: "USD",
	}

	mock.ExpectExec("insert into accounts").
		WithArgs(acc.ID, acc.Currency, acc.Balance).
		WillReturnResult(sqlmock.NewResult(0, 1))

	storage := mydb.PgStorageFromHandle(db)
	storageErr := storage.CreateAccount(context.TODO(), acc)
	if storageErr != nil {
		t.Errorf("Error while creating account: %v", storageErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func Test_PgStorage_CreatePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' while opening a mock database connection", err)
	}
	defer db.Close()

	toAccount := entities.AccountID("bob")
	payment := entities.Payment{
		Account:   "alice",
		Amount:    decimal.New(100, 0),
		ToAccount: &toAccount,
		Direction: entities.Outgoing,
	}

	mock.ExpectQuery("select id, account_id, currency, balance from accounts").
		WithArgs(payment.Account).
		WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "currency", "balance"}).
			AddRow(1, "alice", "USD", decimal.New(100, 0)))

	mock.ExpectQuery("select id, account_id, currency, balance from accounts").
		WithArgs(*payment.ToAccount).
		WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "currency", "balance"}).
			AddRow(2, "bob", "USD", decimal.New(200, 0)))

	mock.ExpectBegin()
	mock.ExpectExec("insert into payments").
		WithArgs(payment.ID, 1, 2, payment.Amount).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("update accounts").
		WithArgs(payment.Amount.Neg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("update accounts").
		WithArgs(payment.Amount, 2).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	storage := mydb.PgStorageFromHandle(db)
	storageErr := storage.CreatePayment(context.TODO(), payment)
	if storageErr != nil {
		t.Errorf("Error while creating payment: %v", storageErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
