package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/khorsl/simple_bank/db/sqlc"
	"github.com/khorsl/simple_bank/token"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !server.isUserAuthorizedToTransfer(ctx, req.FromAccountID, user.ID) {
		return
	}

	if !server.isValidTransfer(ctx, req) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) isValidAccountCurrency(ctx *gin.Context, accountID int64, currency string) bool {
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
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}

func (server *Server) isSufficientBalance(ctx *gin.Context, accountID int64, amount int64) bool {
	fromAccount, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		return false
	}

	if fromAccount.Balance < amount {
		err := fmt.Errorf("accountID %d does not have sufficient funds for %d transfer", accountID, amount)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}

func (server *Server) isValidTransfer(ctx *gin.Context, req transferRequest) bool {
	return server.isValidAccountCurrency(ctx, req.FromAccountID, req.Currency) &&
		server.isValidAccountCurrency(ctx, req.ToAccountID, req.Currency) &&
		server.isSufficientBalance(ctx, req.FromAccountID, req.Amount)
}

func (server *Server) isUserAuthorizedToTransfer(ctx *gin.Context, accountID int64, userId int64) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if userId != account.Owner {
		err := fmt.Errorf("user is unauthorized to transfer money from the account")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return false
	}

	return true
}
