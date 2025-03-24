package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/tasks"
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

// @Summary Get Event by ID
// @Description Get an event by ID
// @Tags Events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} models.SwaggerEventDto
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /events/{id} [get]
// @Security BearerAuth
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

// @Summary Delete Event
// @Description Delete an Event
// @Tags Events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {string} models.EmptyResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /events/{id} [delete]
// @Security BearerAuth
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

// @Summary Create Scheduled Event
// @Description Schedule a new event for future
// @Tags Events
// @Accept json
// @Produce json
// @Param request body models.CreateScheduledEventRequest true "Schedule Event data"
// @Success 200 {object} models.SwaggerEventDto
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /events/schedule [post]
// @Security BearerAuth
func (server *Server) CreateScheduledEvent(ctx *gin.Context) {
	var req models.CreateScheduledEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// Either `RunAtDate` or `RunAfterMins` should be provided
	if (req.RunAtDate == nil && req.RunAfterMins <= 0) || (req.RunAtDate != nil && req.RunAfterMins > 0) {
		err := errors.New("provide either 'run_at_this_date' or 'run_after_x_mins', but not both")
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// If `RunAtDate` is provided, ensure it's in the future
	if req.RunAtDate != nil {
		if req.RunAtDate.Before(time.Now()) {
			err := errors.New("run_at_this_date cannot be in the past")
			ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
			return
		}
	}

	// If `IsRecurring` is true, `IntervalMins` must be greater than 0
	if req.IsRecurring && req.IntervalMins <= 0 {
		err := errors.New("repeat_after_x_mins must be greater than 0")
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// Get the authenticated user from the token
	authPayload := ctx.MustGet(constants.AUTHORIZATION_PAYLOAD_KEY).(*token.Payload)

	// Map the request to the Event model
	arg := models.Event{
		Name:        req.Name,
		Type:        req.Type,
		RunAt:       req.RunAtDate,
		AfterXMins:  req.RunAfterMins,
		Interval:    req.IntervalMins,
		IsRecurring: req.IsRecurring,
		CreatedBy:   authPayload.Username,
	}

	// Store the event in the database
	event, err := server.Repo.CreateEvent(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	// Create an archive task
	scheduleTask, err := tasks.NewScheduleEventTask(event.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
	}

	// Enqueue the task based on RunAt or RunAfter
	var enqueueErr error

	// If `RunAt` is provided → Enqueue with `ProcessAt`
	if event.RunAt != nil {
		enqueueTime := event.RunAt.UTC()
		_, enqueueErr = server.QueueClient.Enqueue(scheduleTask, asynq.Queue("critical"), asynq.ProcessAt(enqueueTime))
	} else {
		// If `RunAfterMins` is provided → Enqueue with `ProcessIn`
		delay := time.Duration(event.AfterXMins) * time.Minute
		_, enqueueErr = server.QueueClient.Enqueue(scheduleTask, asynq.Queue("critical"), asynq.ProcessIn(delay))
	}

	if enqueueErr != nil {
		err := fmt.Errorf("failed to enqueue event: %v", enqueueErr)
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, event)
}
