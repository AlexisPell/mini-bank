package api

import (
	"database/sql"
	"errors"
	db "github.com/alexispell/minibank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, "Couldn't parse body"))
		return
	}

	createAccountArgs := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	acc, err := s.store.CreateAccount(ctx, createAccountArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}

type getAccountRequest struct {
	ID int64 `json:"id" uri:"id" binding:"required,min=1"`
}

//type getAccountResponse struct {
//	ID        int64     `json:"id"`
//	Owner     string    `json:"owner"`
//	Balance   int64     `json:"balance"`
//	Currency  string    `json:"currency"`
//	CreatedAt time.Time `json:"created_at"`
//}

func (s *Server) getAccountById(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err, "ID is required"))
		return
	}

	acc, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err, "Not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "Internal error"))
		return
	}
	ctx.JSON(http.StatusOK, acc)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	listAccountArgs := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := s.store.ListAccounts(ctx, listAccountArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err, "Internal error"))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
