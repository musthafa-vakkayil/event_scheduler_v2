package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/repo"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
	"gorm.io/datatypes"
)

type CreateEventRequest struct {
	Name        string         `json:"name" binding:"required"`
	Type        string         `json:"type" binding:"required,oneof=API"`
	ApiEndpoint string         `json:"api_endpoint" binding:"required"`
	ApiMethod   string         `json:"api_method" binding:"required,oneof=GET POST PUT PATCH DELETE"`
	ApiPayload  datatypes.JSON `json:"api_payload"`
}

func (server *Server) CreateEvent(ctx *gin.Context) {
	var req CreateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(constants.AUTHORIZATION_PAYLOAD_KEY).(*token.Payload)

	// Validation for API Trigger
	if req.Type == "API" && req.ApiMethod != "GET" && req.ApiMethod != "DELETE" && req.ApiPayload == nil {
		err := errors.New("api payload is required for methods other than GET and delete")
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	arg := models.Event{
		Name:           req.Name,
		Type:           req.Type,
		ApiEndPoint:    req.ApiEndpoint,
		ApiMethod:      req.ApiMethod,
		CreatedBy:      authPayload.Username,
		ApiRequestBody: req.ApiPayload,
	}

	event, err := repo.CreateEvent(server.GormDB, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, event)
}

type ListEventsRequest struct {
	PageNumber  int  `form:"pageNumber" binding:"required,min=1"`
	PageSize    int  `form:"pageSize" binding:"required,min=1"`
	CreatedByMe bool `form:"createdByMe"`
}

func (server *Server) ListEvents(ctx *gin.Context) {
	var req ListEventsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	events, err := repo.ListEvents(server.GormDB, req.PageSize, (req.PageNumber-1)*req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, events)
}
