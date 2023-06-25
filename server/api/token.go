package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type renewAccessTokenResponse struct {
	AcessToken          string    `json:"access_token"`
	AcessTokenExpiredAt time.Time `json:"access_token_expired_at"`
}

// renewAccessToken handles creating a new access token with a given refresh token
func (server *Server) renewAccessToken(ctx *gin.Context) {

	refreshToken, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	refreshTokenPayload, err := server.tokenMaker.VerifyToken(refreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshTokenPayload.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("session not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("block session")))
		return
	}

	if session.Username != refreshTokenPayload.Username {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("incorrect sessions user")))
		return
	}

	if session.RefreshToken != refreshToken {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("mismatch refresh token")))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(refreshTokenPayload.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}
	resp := renewAccessTokenResponse{AcessToken: accessToken, AcessTokenExpiredAt: accessTokenPayload.ExpiredAt}

	ctx.JSON(http.StatusOK, gin.H{"ok": true, "payload": resp})
}
