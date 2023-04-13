package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/emmyvera/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAcctID int64  `json:"from_acct_id" binding:"required,min=1"`
	ToAcctID   int64  `json:"to_acct_id" binding:"required,min=1"`
	Amount     int64  `json:"amount" binding:"required,gt=1"`
	Currency   string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(ctx, req.FromAcctID, req.Currency) {
		return
	}

	if !server.validAccount(ctx, req.ToAcctID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAcctID: req.FromAcctID,
		ToAcctID:   req.ToAcctID,
		Amount:     req.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("Account [%d] currency mismatch %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
