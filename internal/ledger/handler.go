package ledger

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DoTransaction performs credit or debit operation
func DoTransaction(store Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zap.L().Info("called deposit handler")

		reqCtx := ctx.Request.Context()
		defer reqCtx.Done()

		ledgerId := ctx.Param("ledgerId")
		if ledgerId == "" {
			ErrorHandler(ctx, http.StatusBadRequest, errors.New("failed get valid ledgerId"))
			return
		}

		var req TransactionRequestDTO
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ErrorHandler(ctx, http.StatusBadRequest, errors.New("failed get valid request payload"))
			return
		}

		if req.Amount <= 0 {
			ErrorHandler(ctx, http.StatusBadRequest, errors.New("failed get amount greater than zero"))
			return
		}

		if !(req.Type == Credit || req.Type == Debit) {
			ErrorHandler(ctx, http.StatusBadRequest, errors.New("failed get transaction type either credit or debit"))
			return
		}

		var res Transaction
		var err error
		if req.Type == Credit {
			res, err = store.Credit(ctx, ledgerId, req)
		} else {
			res, err = store.Debit(ctx, ledgerId, req)
		}

		if err != nil {
			ErrorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("failed to perform transaction: %f, got error: %w", req.Amount, err))
			return
		}

		SuccessHandler(ctx, http.StatusOK, res)
	}
}

// ViewBalance performs view balance operation
func ViewBalance(store Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zap.L().Info("called view balance handler")

		ledgerId := ctx.Param("ledgerId")
		if ledgerId == "" {
			ErrorHandler(ctx, http.StatusBadRequest, errors.New("failed get valid ledgerId"))
			return
		}

		balance, err := store.GetLastBalance(context.Background(), ledgerId)
		if err != nil {
			ErrorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("failed to perform view balance, got error: %w", err))
			return
		}

		SuccessHandler(ctx, http.StatusOK, balance)
	}
}

// ViewTransactionHistory performs view transaction history
func ViewTransactionHistory(store Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zap.L().Info("called view transaction history handler")

		ledgerId := ctx.Param("ledgerId")
		if ledgerId == "" {
			ErrorHandler(ctx, http.StatusBadRequest, errors.New("failed get valid ledgerId"))
			return
		}

		transactions, err := store.GetTransactionHistory(context.Background(), ledgerId)
		if err != nil {
			ErrorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("failed to perform view transaction history, got error: %w", err))
			return
		}

		SuccessHandler(ctx, http.StatusOK, transactions)
	}
}

// ErrorHandler is a function to handle errors
func ErrorHandler(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

// SuccessHandler is a function to handle success
func SuccessHandler(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"data": data})
}
