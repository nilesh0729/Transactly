package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	Anuskh "github.com/nilesh0729/Transactly/internal/db/Result"
	"github.com/nilesh0729/Transactly/internal/token"
)

type listEntryRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) ListEntry(ctx *gin.Context) {
	var req listEntryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accountIDURI := struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}{}
	if err := ctx.ShouldBindUri(&accountIDURI); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify account ownership
	account, err := server.store.GetAccounts(ctx, accountIDURI.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(sql.ErrNoRows)) // Don't leak existence
		return
	}

	arg := Anuskh.ListEntriesParams{
		AccountID: accountIDURI.ID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}

	entries, err := server.store.ListEntries(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

type listTransferRequest struct {
	AccountID int64 `form:"account_id" binding:"required,min=1"`
	PageID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) ListTransfer(ctx *gin.Context) {
	var req listTransferRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Verify account ownership
	account, err := server.store.GetAccounts(ctx, req.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(sql.ErrNoRows))
		return
	}

	arg := Anuskh.ListTransfersParams{
		FromAccountID: req.AccountID,
		ToAccountID:   req.AccountID,
		Limit:         req.PageSize,
		Offset:        (req.PageID - 1) * req.PageSize,
	}

	transfers, err := server.store.ListTransfers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}
