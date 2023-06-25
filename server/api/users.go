package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/eyal-solomon1/ctoviot/db/sqlc"
	"github.com/eyal-solomon1/ctoviot/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var invalidRequesParams error = errors.New("invalid request params")

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string
	FullName          string
	Email             string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
}

// convertToCreateUserResponse converts a db.User object to a userResponse object.
func convertToCreateUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// CreateUser handles the creation of a new user.
func (server *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(invalidRequesParams))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
		Balance:        50,
	})

	if err != nil {
		server.errLogger(UserAPILogGroup, err)
		if pgError, ok := err.(*pq.Error); ok {
			switch pgError.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	responseUser := convertToCreateUserResponse(user)
	server.infoLogger(UserAPILogGroup, fmt.Sprintf("created %v user", user.Username))
	ctx.JSON(http.StatusOK, gin.H{"ok": true, "payload": responseUser})
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	SessionsId            uuid.UUID    `json:"sessions_id"`
	AcessToken            string       `json:"access_token"`
	AcessTokenExpiredAt   time.Time    `json:"access_token_expired_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time    `json:"refresh_token_expired_at"`
	User                  userResponse `json:"user"`
}

// LoginUser handles the login of a user.
func (server *Server) LoginUser(ctx *gin.Context) {

	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(invalidRequesParams))
		return

	}

	user, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return

	}

	accessToken, accessPyaload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userRespons := convertToCreateUserResponse(user)
	rsps := loginUserResponse{
		SessionsId:            session.ID,
		AcessToken:            accessToken,
		AcessTokenExpiredAt:   accessPyaload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshPayload.ExpiredAt,
		User:                  userRespons,
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("refresh_token", refreshToken, int(time.Until(refreshPayload.ExpiredAt).Seconds()), "/", server.config.CookieDomainEndpoint, true, true)

	ctx.JSON(http.StatusOK, gin.H{"ok": true, "payload": rsps})

}

// LogoutUser handles the logout of a user.
func (server *Server) LogoutUser(ctx *gin.Context) {
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

	if session.Username != refreshTokenPayload.Username {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("incorrect sessions user")))
		return
	}

	if session.RefreshToken != refreshToken {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("mismatch refresh token")))
		return
	}

	_, err = server.store.BlockSession(ctx, session.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusNoContent, gin.H{"ok": true})
}
