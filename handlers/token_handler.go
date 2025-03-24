package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
)

// @Summary Renew Access Token
// @Description Get a new access token using the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RenewAccessTokenRequest true "Token data"
// @Success 200 {object} models.RenewAccessTokenResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /token/renew [post]
func (server *Server) RenewAccessToken(ctx *gin.Context) {
	var req models.RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	refreshPayload, err := server.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, constants.ErrorResponse(err))
		return
	}

	session, err := server.Repo.GetSession(refreshPayload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	if session.IsBlocked {
		err := errors.New("blocked session")
		ctx.JSON(http.StatusUnauthorized, constants.ErrorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := errors.New("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, constants.ErrorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := errors.New("invalid refresh token")
		ctx.JSON(http.StatusUnauthorized, constants.ErrorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := errors.New("expired refresh token")
		ctx.JSON(http.StatusUnauthorized, constants.ErrorResponse(err))
		return
	}

	token, accessPayload, err := server.TokenMaker.CreateToken(refreshPayload.Username, server.Config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	response := models.RenewAccessTokenResponse{
		AccessToken:          token,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}

	ctx.JSON(http.StatusOK, response)
}
