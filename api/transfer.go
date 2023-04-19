package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/emmyvera/simplebank/db/sqlc"
	"github.com/emmyvera/simplebank/token"
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

	fromAccount, valid := server.validAccount(ctx, req.FromAcctID, req.Currency)

	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != fromAccount.Owner {
		err := errors.New("Account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAcctID, req.Currency)

	if !valid {
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

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("Account [%d] currency mismatch %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
