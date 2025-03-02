package mock

import (
	"context"
	"fmt"

	"github.com/dineshd30/ledger-service/internal/ledger"
	"github.com/stretchr/testify/mock"
)

type Store struct {
	mock.Mock
}

func (s *Store) Credit(ctx context.Context, ledgerId string, trd ledger.TransactionRequestDTO) (ledger.Transaction, error) {
	fmt.Println("Called mocked Credit function")
	args := s.Called(ctx, ledgerId, trd)
	return args.Get(0).(ledger.Transaction), args.Error(1)
}

func (s *Store) Debit(ctx context.Context, ledgerId string, trd ledger.TransactionRequestDTO) (ledger.Transaction, error) {
	fmt.Println("Called mocked Debit function")
	args := s.Called(ctx, ledgerId, trd)
	return args.Get(0).(ledger.Transaction), args.Error(1)
}

func (s *Store) GetLastBalance(ctx context.Context, ledgerId string) (float64, error) {
	fmt.Println("Called mocked GetLastBalance function")
	args := s.Called(ctx, ledgerId)
	return args.Get(0).(float64), args.Error(1)
}

func (s *Store) GetTransactionHistory(ctx context.Context, ledgerId string) ([]ledger.Transaction, error) {
	fmt.Println("Called mocked GetTransactionHistory function")
	args := s.Called(ctx, ledgerId)
	return args.Get(0).([]ledger.Transaction), args.Error(1)
}
