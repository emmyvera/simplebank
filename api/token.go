package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenUserRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenUserResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessTokenUser(ctx *gin.Context) {
	var req renewAccessTokenUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshToken, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshToken.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check If It's Blocked
	if session.IsBlocked {
		err := fmt.Errorf("Session is Blocked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check User
	if session.Username != refreshToken.Username {
		err := fmt.Errorf("Unauthorized User")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check token
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("Mismatch Token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check Token Expire Date
	if time.Now().After(session.ExpiredAt) {
		err := fmt.Errorf("Session Expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshToken.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenUserResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
