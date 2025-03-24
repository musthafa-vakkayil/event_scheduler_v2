package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
)

// @Summary Create API Event
// @Description Create a new API Event
// @Tags Events
// @Accept json
// @Produce json
// @Param request body models.CreateAPIEventRequestSwagger true "API Event data"
// @Success 200 {object} models.SwaggerEventDto
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /events/api [post]
// @Security BearerAuth
func (server *Server) CreateAPIEvent(ctx *gin.Context) {
	var req models.CreateAPIEventRequest
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

	event, err := server.Repo.CreateEvent(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, event)
}

// @Summary List Events
// @Description List all events
// @Tags Events
// @Produce json
// @Param pageNumber query int true "Page Number" minimum(1)
// @Param pageSize query int true "Page Size" minimum(1)
// @Success 200 {object} models.SwaggerListEventResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /events [get]
// @Security BearerAuth
func (server *Server) ListEvents(ctx *gin.Context) {
	var req models.ListEventsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	events, err := server.Repo.ListEvents(req.PageSize, (req.PageNumber-1)*req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func (server *Server) GetEvent(ctx *gin.Context) {
	var req models.GetEventRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	event, err := server.Repo.GetEvent(req.ID)
	if err != nil {
		if err.Error() == "event not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (server *Server) DeleteEvent(ctx *gin.Context) {
	var req models.GetEventRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	err := server.Repo.DeleteEvent(int64(req.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}
