package ledger

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
)

// TransactionType represents whether the transaction is a Credit or Debit
type TransactionType string

const (
	Credit TransactionType = "credit"
	Debit  TransactionType = "debit"
)

// Transaction represents a single ledger entry
type Transaction struct {
	ID             string          `json:"id"`
	Date           int64           `json:"date"`
	Type           TransactionType `json:"type"`
	Description    string          `json:"description"`
	Amount         float64         `json:"amount"`
	RunningBalance float64         `json:"runningBalance"`
}

// Ledger holds the ledger metadata and transaction history
type Ledger struct {
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	Transactions []Transaction `json:"transactions"`
}

// TransactionRequestDTO represents the request payload for deposit and withdraw operations
type TransactionRequestDTO struct {
	Type        TransactionType `json:"type"`
	Description string          `json:"description"`
	Amount      float64         `json:"amount"`
}

// Store represents the operations on the ledger
type Store interface {
	Credit(ctx context.Context, ledgerId string, trd TransactionRequestDTO) (Transaction, error)
	Debit(ctx context.Context, ledgerId string, trd TransactionRequestDTO) (Transaction, error)
	GetLastBalance(ctx context.Context, ledgerId string) (float64, error)
	GetTransactionHistory(ctx context.Context, ledgerId string) ([]Transaction, error)
}

// store is our in-memory implementation of Store
type store struct {
	uuid    UUIDGenerator
	ledgers map[string]*Ledger
}

// NewStore creates a new in-memory store instance
func NewStore(uuid UUIDGenerator, ledgers map[string]*Ledger) Store {
	return &store{
		uuid:    uuid,
		ledgers: ledgers,
	}
}

// Credit adds a credit transaction to the ledger
func (s *store) Credit(ctx context.Context, ledgerId string, trd TransactionRequestDTO) (Transaction, error) {
	ledger, lastBalance, err := s.getLedgerWithBalance(ledgerId)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to perform credit transaction, got error : %w", err)
	}

	newBalance := lastBalance + trd.Amount
	newTransaction := Transaction{
		ID:             s.uuid.Generate(),
		Date:           time.Now().UTC().UnixMilli(),
		Type:           trd.Type,
		Description:    trd.Description,
		Amount:         trd.Amount,
		RunningBalance: round(newBalance, 4),
	}
	ledger.Transactions = append(ledger.Transactions, newTransaction)
	zap.L().Info("credited the ledger", zap.String("ledgerId", ledgerId), zap.Float64("newBalance", newBalance))
	return newTransaction, nil
}

// Debit subtracts an amount from the ledger
func (s *store) Debit(ctx context.Context, ledgerId string, trd TransactionRequestDTO) (Transaction, error) {
	ledger, lastBalance, err := s.getLedgerWithBalance(ledgerId)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to perform debit transaction, got error : %w", err)
	}

	newBalance := lastBalance - trd.Amount
	if newBalance <= 0 {
		return Transaction{}, errors.New("failed to get new balance greater than or equal to 0")
	}

	newTransaction := Transaction{
		ID:             s.uuid.Generate(),
		Date:           time.Now().UTC().UnixMilli(),
		Type:           trd.Type,
		Description:    trd.Description,
		Amount:         trd.Amount,
		RunningBalance: round(newBalance, 4),
	}
	ledger.Transactions = append(ledger.Transactions, newTransaction)
	zap.L().Info("debited the ledger", zap.String("ledgerId", ledgerId), zap.Float64("newBalance", newBalance))
	return newTransaction, nil
}

// GetLastBalance returns the last balance for ledger
func (s *store) GetLastBalance(ctx context.Context, ledgerId string) (float64, error) {
	_, lastBalance, err := s.getLedgerWithBalance(ledgerId)
	if err != nil {
		return 0, fmt.Errorf("failed to get last balance, got error : %w", err)
	}

	zap.L().Info("got last ledger balance", zap.String("ledgerId", ledgerId), zap.Float64("lastBalance", lastBalance))
	return lastBalance, nil
}

// GetTransactionHistory returns the transaction history for ledger
func (s *store) GetTransactionHistory(ctx context.Context, ledgerId string) ([]Transaction, error) {
	ledger, _, err := s.getLedgerWithBalance(ledgerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction history, got error : %w", err)
	}

	zap.L().Info("got transaction history for ledger", zap.String("ledgerId", ledgerId))
	return ledger.Transactions, nil
}

// getLedgerWithBalance retrieves the ledger, last balance by ledgerId
func (s *store) getLedgerWithBalance(id string) (*Ledger, float64, error) {
	lastBalance := 0.0
	ledger, exists := s.ledgers[id]
	if !exists {
		return nil, lastBalance, fmt.Errorf("failed get ledger: %s", id)
	}

	if len(ledger.Transactions) > 0 {
		lastBalance = ledger.Transactions[len(ledger.Transactions)-1].RunningBalance
	}

	return ledger, lastBalance, nil
}

// round rounds float value to specific precision places
func round(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
