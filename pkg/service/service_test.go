package service_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"

	"github.com/shirolimit/wallet-service/pkg/service"

	"github.com/golang/mock/gomock"
	"github.com/shirolimit/wallet-service/pkg/db"
	"github.com/shirolimit/wallet-service/pkg/entities"
)

func Test_walletService_CreateAccount(t *testing.T) {
	type args struct {
		acc          entities.Account
		storageError error
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantCall bool
	}{
		{
			"error_on_empty_account_id",
			args{acc: entities.Account{ID: "", Currency: "USD", Balance: decimal.New(100, 0)}},
			true,
			false,
		},
		{
			"error_on_empty_currency",
			args{acc: entities.Account{ID: "alice", Currency: "", Balance: decimal.New(100, 0)}},
			true,
			false,
		},
		{
			"error_on_negative_balance",
			args{acc: entities.Account{ID: "alice", Currency: "USD", Balance: decimal.New(-100, 0)}},
			true,
			false,
		},
		{
			"calls_storage",
			args{acc: entities.Account{ID: "alice", Currency: "USD", Balance: decimal.New(100, 0)}},
			false,
			true,
		},
		{
			"error_on_storage_error",
			args{
				acc:          entities.Account{ID: "alice", Currency: "USD", Balance: decimal.New(100, 0)},
				storageError: entities.ErrAccountAlreadyExists,
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := db.NewMockStorage(ctrl)
			svc := service.NewWalletService(mockStorage)
			if tt.wantCall {
				mockStorage.EXPECT().CreateAccount(context.TODO(), tt.args.acc).Return(tt.args.storageError)
			}
			if err := svc.CreateAccount(context.TODO(), tt.args.acc); (err != nil) != tt.wantErr {
				t.Errorf("walletService.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_walletService_ListAccounts(t *testing.T) {
	type args struct {
		storageData  []entities.AccountID
		storageError error
	}
	tests := []struct {
		name    string
		args    args
		want    []entities.AccountID
		wantErr bool
	}{
		{
			"calls_storage",
			args{storageData: []entities.AccountID{"alice", "bob"}, storageError: nil},
			[]entities.AccountID{"alice", "bob"},
			false,
		},
		{
			"error_on_storage_error",
			args{storageData: nil, storageError: entities.ErrDatabaseConnection},
			[]entities.AccountID{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := db.NewMockStorage(ctrl)
			svc := service.NewWalletService(mockStorage)

			mockStorage.EXPECT().ListAccounts(context.TODO()).Return(tt.args.storageData, tt.args.storageError)
			got, err := svc.ListAccounts(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("walletService.ListAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("walletService.ListAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_walletService_GetAccount(t *testing.T) {
	type args struct {
		id           entities.AccountID
		storageData  *entities.Account
		storageError error
	}
	tests := []struct {
		name     string
		args     args
		want     entities.Account
		wantErr  bool
		wantCall bool
	}{
		{
			"calls_storage",
			args{
				storageData:  &entities.Account{ID: "alice", Balance: decimal.New(100, 0), Currency: "USD"},
				storageError: nil,
				id:           entities.AccountID("alice"),
			},
			entities.Account{ID: "alice", Balance: decimal.New(100, 0), Currency: "USD"},
			false,
			true,
		},
		{
			"error_on_storage_error",
			args{
				storageData:  nil,
				storageError: entities.ErrDatabaseConnection,
				id:           entities.AccountID("alice"),
			},
			entities.Account{},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := db.NewMockStorage(ctrl)
			svc := service.NewWalletService(mockStorage)

			if tt.wantCall {
				mockStorage.EXPECT().GetAccount(context.TODO(), tt.args.id).
					Return(tt.args.storageData, tt.args.storageError)
			}
			got, err := svc.GetAccount(context.TODO(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("walletService.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("walletService.GetAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func accountIDRef(id string) *entities.AccountID {
	accId := entities.AccountID(id)
	return &accId
}
func Test_walletService_GetPayments(t *testing.T) {
	type args struct {
		id           entities.AccountID
		storageData  []entities.Payment
		storageError error
	}
	tests := []struct {
		name    string
		args    args
		want    []entities.Payment
		wantErr bool
	}{
		{
			"calls_storage",
			args{
				id: "alice",
				storageData: []entities.Payment{
					{
						Account:   "alice",
						Amount:    decimal.New(100, 0),
						Direction: entities.Outgoing,
						ToAccount: accountIDRef("bob"),
					},
				},
				storageError: nil,
			},
			[]entities.Payment{
				{
					Account:   "alice",
					Amount:    decimal.New(100, 0),
					Direction: entities.Outgoing,
					ToAccount: accountIDRef("bob"),
				},
			},
			false,
		},
		{
			"error_on_storage_error",
			args{
				id:           "alice",
				storageData:  []entities.Payment{},
				storageError: entities.ErrDatabaseConnection,
			},
			[]entities.Payment{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := db.NewMockStorage(ctrl)
			svc := service.NewWalletService(mockStorage)

			mockStorage.EXPECT().PaymentsByAccount(context.TODO(), tt.args.id).Return(tt.args.storageData, tt.args.storageError)
			got, err := svc.GetPayments(context.TODO(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("walletService.GetPayments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("walletService.GetPayments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_walletService_MakePayment(t *testing.T) {
	type args struct {
		payment      entities.Payment
		storageError error
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantCall bool
	}{
		{
			"error_on_empty_account",
			args{
				payment: entities.Payment{
					ID:        uuid.New(),
					Account:   "",
					ToAccount: accountIDRef("bob"),
					Amount:    decimal.New(100, 0),
					Direction: entities.Outgoing,
				},
			},
			true,
			false,
		},
		{
			"error_on_empty_to_account",
			args{
				payment: entities.Payment{
					ID:        uuid.New(),
					Account:   "alice",
					ToAccount: accountIDRef(""),
					Amount:    decimal.New(100, 0),
					Direction: entities.Outgoing,
				},
			},
			true,
			false,
		},
		{
			"error_on_nil_to_account",
			args{
				payment: entities.Payment{
					ID:        uuid.New(),
					Account:   "alice",
					ToAccount: nil,
					Amount:    decimal.New(100, 0),
					Direction: entities.Outgoing,
				},
			},
			true,
			false,
		},
		{
			"error_on_negative_amount",
			args{
				payment: entities.Payment{
					ID:        uuid.New(),
					Account:   "alice",
					ToAccount: accountIDRef("bob"),
					Amount:    decimal.New(-100, 0),
					Direction: entities.Outgoing,
				},
			},
			true,
			false,
		},
		{
			"error_on_zero_amount",
			args{
				payment: entities.Payment{
					ID:        uuid.New(),
					Account:   "alice",
					ToAccount: accountIDRef("bob"),
					Amount:    decimal.New(0, 0),
					Direction: entities.Outgoing,
				},
			},
			true,
			false,
		},
		{
			"error_on_same_amount",
			args{
				payment: entities.Payment{
					ID:        uuid.New(),
					Account:   "alice",
					ToAccount: accountIDRef("alice"),
					Amount:    decimal.New(100, 0),
					Direction: entities.Outgoing,
				},
			},
			true,
			false,
		},
		{
			"error_on_empty_id",
			args{
				payment: entities.Payment{
					Account:   "alice",
					ToAccount: accountIDRef("bob"),
					Amount:    decimal.New(100, 0),
					Direction: entities.Outgoing,
				},
			},
			true,
			false,
		},
		{
			"error_on_incoming_payment",
			args{
				payment: entities.Payment{
					ID:        uuid.New(),
					Account:   "alice",
					ToAccount: accountIDRef("bob"),
					Amount:    decimal.New(100, 0),
					Direction: entities.Incoming,
				},
			},
			true,
			false,
		},
		{
			"storage_called",
			args{
				payment: entities.Payment{
					ID:        uuid.New(),
					Account:   "alice",
					ToAccount: accountIDRef("bob"),
					Amount:    decimal.New(100, 0),
					Direction: entities.Outgoing,
				},
			},
			false,
			true,
		},
		{
			"error_on_storage_error",
			args{
				payment: entities.Payment{
					ID:        uuid.New(),
					Account:   "alice",
					ToAccount: accountIDRef("bob"),
					Amount:    decimal.New(100, 0),
					Direction: entities.Outgoing,
				},
				storageError: entities.ErrDatabaseConnection,
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := db.NewMockStorage(ctrl)
			svc := service.NewWalletService(mockStorage)

			if tt.wantCall {
				mockStorage.EXPECT().CreatePayment(context.TODO(), tt.args.payment).
					Return(tt.args.storageError)
			}
			if err := svc.MakePayment(context.TODO(), tt.args.payment); (err != nil) != tt.wantErr {
				t.Errorf("walletService.MakePayment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
