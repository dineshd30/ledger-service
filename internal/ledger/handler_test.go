package ledger_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dineshd30/ledger-service/internal/ledger"
	internalMock "github.com/dineshd30/ledger-service/internal/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDoTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                    string
		ledgerId                string
		requestBody             string
		expectedStatus          int
		expectedResponseField   string
		expectedResponseMessage interface{}
		storeSetup              func() ledger.Store
	}{
		{
			name:                    "Missing ledgerId parameter",
			ledgerId:                "",
			requestBody:             `{"ledgerId": "ledger1", "ledgerType": "cash", "type": "credit", "description": "deposit", "amount": 100}`,
			expectedStatus:          http.StatusBadRequest,
			expectedResponseField:   "error",
			expectedResponseMessage: "failed get valid ledgerId",
			storeSetup: func() ledger.Store {
				return new(internalMock.Store)
			},
		},
		{
			name:                    "Invalid JSON payload",
			ledgerId:                "ledger1",
			requestBody:             `invalid json`,
			expectedStatus:          http.StatusBadRequest,
			expectedResponseField:   "error",
			expectedResponseMessage: "failed get valid request payload",
			storeSetup: func() ledger.Store {
				return new(internalMock.Store)
			},
		},
		{
			name:                    "Amount less than or equal to zero",
			ledgerId:                "ledger1",
			requestBody:             `{"ledgerId": "ledger1", "ledgerType": "cash", "type": "credit", "description": "deposit", "amount": 0}`,
			expectedStatus:          http.StatusBadRequest,
			expectedResponseField:   "error",
			expectedResponseMessage: "failed get amount greater than zero",
			storeSetup: func() ledger.Store {
				return new(internalMock.Store)
			},
		},
		{
			name:                    "Invalid transaction type",
			ledgerId:                "ledger1",
			requestBody:             `{"ledgerId": "ledger1", "ledgerType": "cash", "type": "invalid", "description": "deposit", "amount": 100}`,
			expectedStatus:          http.StatusBadRequest,
			expectedResponseField:   "error",
			expectedResponseMessage: "failed get transaction type either credit or debit",
			storeSetup: func() ledger.Store {
				return new(internalMock.Store)
			},
		},
		{
			name:                    "Successful credit transaction",
			ledgerId:                "ledger1",
			requestBody:             `{"ledgerId": "ledger1", "ledgerType": "cash", "type": "credit", "description": "deposit", "amount": 100}`,
			expectedStatus:          http.StatusOK,
			expectedResponseField:   "data",
			expectedResponseMessage: ledger.Transaction{ID: "tx-credit-1", Date: 1234567890, Type: ledger.Credit, Description: "deposit", Amount: 100, RunningBalance: 100},
			storeSetup: func() ledger.Store {
				mStore := new(internalMock.Store)
				reqDTO := ledger.TransactionRequestDTO{
					Type:        ledger.Credit,
					Description: "deposit",
					Amount:      100,
				}
				mStore.On("Credit", mock.Anything, "ledger1", reqDTO).Return(ledger.Transaction{
					ID:             "tx-credit-1",
					Date:           1234567890,
					Type:           ledger.Credit,
					Description:    "deposit",
					Amount:         100,
					RunningBalance: 100,
				}, nil)
				return mStore
			},
		},
		{
			name:                    "Successful debit transaction",
			ledgerId:                "ledger1",
			requestBody:             `{"ledgerId": "ledger1", "ledgerType": "cash", "type": "debit", "description": "withdrawal", "amount": 50}`,
			expectedStatus:          http.StatusOK,
			expectedResponseField:   "data",
			expectedResponseMessage: ledger.Transaction{ID: "tx-debit-1", Date: 1234567891, Type: ledger.Debit, Description: "withdrawal", Amount: 50, RunningBalance: 50},
			storeSetup: func() ledger.Store {
				mStore := new(internalMock.Store)
				reqDTO := ledger.TransactionRequestDTO{

					Type:        ledger.Debit,
					Description: "withdrawal",
					Amount:      50,
				}
				mStore.On("Debit", mock.Anything, "ledger1", reqDTO).Return(ledger.Transaction{
					ID:             "tx-debit-1",
					Date:           1234567891,
					Type:           ledger.Debit,
					Description:    "withdrawal",
					Amount:         50,
					RunningBalance: 50,
				}, nil)
				return mStore
			},
		},
		{
			name:                    "Store error during credit transaction",
			ledgerId:                "ledger1",
			requestBody:             `{"ledgerId": "ledger1", "ledgerType": "cash", "type": "credit", "description": "deposit", "amount": 100}`,
			expectedStatus:          http.StatusInternalServerError,
			expectedResponseField:   "error",
			expectedResponseMessage: "failed to perform transaction: 100.000000, got error: store error",
			storeSetup: func() ledger.Store {
				mStore := new(internalMock.Store)
				reqDTO := ledger.TransactionRequestDTO{
					Type:        ledger.Credit,
					Description: "deposit",
					Amount:      100,
				}
				mStore.On("Credit", mock.Anything, "ledger1", reqDTO).Return(ledger.Transaction{}, errors.New("store error"))
				return mStore
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := tc.storeSetup()

			req := httptest.NewRequest("POST", "/ledger/:ledgerId/transaction", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			if tc.ledgerId != "" {
				c.Params = []gin.Param{{Key: "ledgerId", Value: tc.ledgerId}}
			}
			c.Request = req

			handler := ledger.DoTransaction(store)
			handler(c)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tc.expectedResponseField == "error" {
				assert.Contains(t, resp, "error")
				assert.Equal(t, tc.expectedResponseMessage, resp["error"])
			} else {
				data, ok := resp["data"].(map[string]interface{})
				assert.True(t, ok)
				expectedTx := tc.expectedResponseMessage.(ledger.Transaction)
				assert.Equal(t, expectedTx.ID, data["id"])
				assert.Equal(t, string(expectedTx.Type), data["type"])
				assert.Equal(t, expectedTx.Description, data["description"])
				assert.Equal(t, expectedTx.Amount, data["amount"])
				assert.Equal(t, expectedTx.RunningBalance, data["runningBalance"])
			}
		})
	}
}
