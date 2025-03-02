package ledger_test

import (
	"context"
	"testing"

	"github.com/dineshd30/ledger-service/internal/ledger"
	internalMock "github.com/dineshd30/ledger-service/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestStoreCredit(t *testing.T) {
	tests := []struct {
		name            string
		ledgerId        string
		initialLedger   *ledger.Ledger
		creditRequest   ledger.TransactionRequestDTO
		expectedBalance float64
		expectError     bool
	}{
		{
			name:     "Credit transaction on new ledger",
			ledgerId: "ledger1",
			initialLedger: &ledger.Ledger{
				ID:   "ledger1",
				Type: "cash",
			},
			creditRequest: ledger.TransactionRequestDTO{
				Type:        ledger.Credit,
				Description: "initial deposit",
				Amount:      100,
			},
			expectedBalance: 100,
			expectError:     false,
		},
		{
			name:     "Credit transaction on existing ledger",
			ledgerId: "ledger2",
			initialLedger: &ledger.Ledger{
				ID:   "ledger2",
				Type: "cash",
				Transactions: []ledger.Transaction{
					{
						ID:             "tx-old",
						Date:           1234567890,
						Type:           ledger.Credit,
						Description:    "previous deposit",
						Amount:         50,
						RunningBalance: 50,
					},
				},
			},
			creditRequest: ledger.TransactionRequestDTO{
				Type:        ledger.Credit,
				Description: "additional deposit",
				Amount:      75,
			},
			expectedBalance: 125,
			expectError:     false,
		},
		{
			name:          "Credit transaction with empty ledgerId returns error",
			ledgerId:      "",
			initialLedger: nil,
			creditRequest: ledger.TransactionRequestDTO{
				Type:        ledger.Credit,
				Description: "deposit with empty ledgerId",
				Amount:      100,
			},
			expectedBalance: 0,
			expectError:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ledgers := make(map[string]*ledger.Ledger)
			if tc.initialLedger != nil {
				ledgers[tc.ledgerId] = tc.initialLedger
			}
			uuid := internalMock.UUIDGenerator{}
			uuid.On("Generate").Return("123")
			storeInstance := ledger.NewStore(&uuid, ledgers)

			tx, err := storeInstance.Credit(context.Background(), tc.ledgerId, tc.creditRequest)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBalance, tx.RunningBalance)

				l, ok := ledgers[tc.ledgerId]
				assert.True(t, ok)
				assert.NotEmpty(t, l.Transactions)
				lastTx := l.Transactions[len(l.Transactions)-1]
				assert.Equal(t, tc.expectedBalance, lastTx.RunningBalance)
			}
		})
	}
}

func TestStoreDebit(t *testing.T) {
	tests := []struct {
		name            string
		ledgerId        string
		initialLedger   *ledger.Ledger
		debitRequest    ledger.TransactionRequestDTO
		expectedBalance float64
		expectError     bool
	}{
		{
			name:     "Debit transaction on existing ledger",
			ledgerId: "ledger3",
			initialLedger: &ledger.Ledger{
				ID:   "ledger3",
				Type: "cash",
				Transactions: []ledger.Transaction{
					{
						ID:             "tx-old",
						Date:           1234567890,
						Type:           ledger.Credit,
						Description:    "previous deposit",
						Amount:         200,
						RunningBalance: 200,
					},
				},
			},
			debitRequest: ledger.TransactionRequestDTO{
				Type:        ledger.Debit,
				Description: "withdrawal",
				Amount:      75,
			},
			expectedBalance: 125,
			expectError:     false,
		},
		{
			name:          "Debit transaction on non-existent ledger returns error",
			ledgerId:      "ledger4",
			initialLedger: nil,
			debitRequest: ledger.TransactionRequestDTO{
				Type:        ledger.Debit,
				Description: "withdrawal",
				Amount:      50,
			},
			expectedBalance: 0,
			expectError:     true,
		},
		{
			name:     "Debit transaction with insufficient funds",
			ledgerId: "ledger5",
			initialLedger: &ledger.Ledger{
				ID:   "ledger5",
				Type: "cash",
				Transactions: []ledger.Transaction{
					{
						ID:             "tx-old",
						Date:           1234567890,
						Type:           ledger.Credit,
						Description:    "previous deposit",
						Amount:         50,
						RunningBalance: 50,
					},
				},
			},
			debitRequest: ledger.TransactionRequestDTO{
				Type:        ledger.Debit,
				Description: "withdrawal exceeding funds",
				Amount:      75,
			},
			expectedBalance: 0,
			expectError:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ledgers := make(map[string]*ledger.Ledger)
			if tc.initialLedger != nil {
				ledgers[tc.ledgerId] = tc.initialLedger
			}
			uuid := internalMock.UUIDGenerator{}
			uuid.On("Generate").Return("123")
			storeInstance := ledger.NewStore(&uuid, ledgers)

			tx, err := storeInstance.Debit(context.Background(), tc.ledgerId, tc.debitRequest)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBalance, tx.RunningBalance)
				l, ok := ledgers[tc.ledgerId]
				assert.True(t, ok)
				assert.NotEmpty(t, l.Transactions)
				lastTx := l.Transactions[len(l.Transactions)-1]
				assert.Equal(t, tc.expectedBalance, lastTx.RunningBalance)
			}
		})
	}
}
