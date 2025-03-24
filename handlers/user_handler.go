package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
	"github.com/musthafa-vakkayil/event_scheduler_v2/utils"
)

// @Summary Create a new user
// @Description Create a new user with username, password, and email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User data"
// @Success 200 {object} models.UserDto
// @Failure 400 {object} models.BadRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /users [post]
func (server *Server) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	user := models.User{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	usr, err := server.Repo.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	userDto := models.ConvertToUserDto(usr)

	ctx.JSON(http.StatusOK, userDto)
}

// @Summary Get user by username
// @Description Retrieve user details by username
// @Tags Users
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} models.UserDto
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 404 {object} models.NotFoundResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /users/{username} [get]
// @Security BearerAuth
func (server *Server) GetUser(ctx *gin.Context) {
	var req models.GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	user, err := server.Repo.GetUser(req.Username)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	userDto := models.ConvertToUserDto(user)

	ctx.JSON(http.StatusOK, userDto)
}

// @Summary Get Access Token
// @Description Authenticate and get access/refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 404 {object} models.NotFoundResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /login [post]
func (server *Server) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	userData, err := server.Repo.GetUser(req.Username)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, userData.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, constants.ErrorResponse(err))
		return
	}

	token, accessPayload, err := server.TokenMaker.CreateToken(req.Username, server.Config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.TokenMaker.CreateToken(req.Username, server.Config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	args := models.Session{
		ID:           refreshPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIP:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
		IsBlocked:    false,
	}

	session, err := server.Repo.CreateSession(args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	user := models.ConvertToUserDto(userData)

	response := models.LoginResponse{
		SessionId:             session.ID,
		AccessToken:           token,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  user,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Delete user by username
// @Description Remove a user by username
// @Tags Users
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {string} models.EmptyResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /users/{username} [delete]
// @Security BearerAuth
func (server *Server) DeleteUser(ctx *gin.Context) {
	var req models.GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(constants.AUTHORIZATION_PAYLOAD_KEY).(*token.Payload)

	if authPayload.Username != req.Username {
		err := errors.New("cannot delete other users")
		ctx.JSON(http.StatusUnauthorized, constants.ErrorResponse(err))
		return
	}

	err := server.Repo.DeleteUser(req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}
