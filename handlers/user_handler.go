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

	token, err := server.TokenMaker.CreateToken(req.Username, server.Config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	user := models.ConvertToUserDto(userData)

	response := models.LoginResponse{
		User:  user,
		Token: token,
	}

	ctx.JSON(http.StatusOK, response)
}

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
