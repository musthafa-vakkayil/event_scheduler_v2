package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
	"github.com/musthafa-vakkayil/event_scheduler_v2/utils"
	"github.com/musthafa-vakkayil/event_scheduler_v2/validator"
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
	err := validator.ValidateAPIEventRequest(req)
	if err != nil {
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

	err := server.Repo.DeleteEvent(req.ID)
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

	err := validator.ValidateScheduleEventRequest(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// If `IsRecurring` is true, `IntervalMins` must be greater than 0
	if req.IsRecurring && req.IntervalMins <= 0 {
		err := errors.New("repeat_after_x_mins must be greater than 0")
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// Get the authenticated user from the token
	authPayload := ctx.MustGet(constants.AUTHORIZATION_PAYLOAD_KEY).(*token.Payload)

	// Parse time using the shared function
	runTime, err := utils.ParseAndValidateTime(req.RunAtDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// Map the request to the Event model
	arg := models.Event{
		Name:        req.Name,
		Type:        req.Type,
		RunAt:       runTime,
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

	// Enqueue the task based on RunAt or RunAfter
	var taskId string
	var enqueueErr error

	// If `RunAt` is provided → Enqueue with `ProcessAt`
	if event.RunAt != nil {
		taskId, enqueueErr = server.TaskManager.EnqueueTaskAt(ctx, event.ID, constants.TASK_SCHEDULE_EVENT, constants.CRITICAL_PRIORITY_QUEUE, *event.RunAt)
	} else {
		// If `RunAfterMins` is provided → Enqueue with `ProcessIn`
		delay := time.Duration(event.AfterXMins) * time.Minute
		taskId, enqueueErr = server.TaskManager.EnqueueTaskIn(ctx, event.ID, constants.TASK_SCHEDULE_EVENT, constants.CRITICAL_PRIORITY_QUEUE, delay)
	}

	if enqueueErr != nil {
		err := fmt.Errorf("failed to enqueue event: %v", enqueueErr)
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	eventTask := models.EventTask{
		EventID:    int(event.ID),
		TaskID:     taskId,
		IsCanceled: false,
	}

	err = server.Repo.StoreTaskID(eventTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, event)
}

// @Summary Edit Scheduled Event
// @Description Edit a Scheduled Event Details
// @Tags Events
// @Accept json
// @Produce json
// @Param request body models.CreateScheduledEventRequest true "Schedule Event data"
// @Success 200 {object} models.SwaggerEventDto
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /events/schedule/{id} [put]
// @Security BearerAuth
func (server *Server) UpdateScheduledEvent(ctx *gin.Context) {
	// Parse URI for event ID
	var mod models.GetEventRequest
	if err := ctx.ShouldBindUri(&mod); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// Parse request body
	var req models.CreateScheduledEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// Validate the request
	err := validator.ValidateScheduleEventRequest(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// Fetch the existing event
	existingEvent, err := server.Repo.GetEvent(mod.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, constants.ErrorResponse(err))
		return
	}

	// Parse new time
	var runTime *time.Time
	if req.RunAtDate != "" {
		parsedTime, err := utils.ParseAndValidateTime(req.RunAtDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
			return
		}
		runTime = parsedTime
	}

	// Update event details
	updatedEvent := models.Event{
		ID:          int64(mod.ID),
		Name:        req.Name,
		Type:        req.Type,
		RunAt:       runTime,
		AfterXMins:  req.RunAfterMins,
		Interval:    req.IntervalMins,
		IsRecurring: req.IsRecurring,
	}

	// Update event in the database
	event, err := server.Repo.UpdateEvent(updatedEvent)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	// Cancel the old task if the time or recurrence settings are changing
	timeChanged := (runTime != nil && existingEvent.RunAt != nil && !runTime.Equal(*existingEvent.RunAt)) ||
		(req.RunAfterMins != existingEvent.AfterXMins)

	if timeChanged {
		// Cancel old tasks
		if err := server.Repo.CancelTasksForEvent(mod.ID); err != nil {
			ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
			return
		}

		// Enqueue the task with new time settings
		var taskID string
		var enqueueErr error

		if runTime != nil {
			taskID, enqueueErr = server.TaskManager.EnqueueTaskAt(ctx, event.ID, constants.TASK_SCHEDULE_EVENT, constants.CRITICAL_PRIORITY_QUEUE, *runTime)
		} else {
			delay := time.Duration(event.AfterXMins) * time.Minute
			taskID, enqueueErr = server.TaskManager.EnqueueTaskIn(ctx, event.ID, constants.TASK_SCHEDULE_EVENT, constants.CRITICAL_PRIORITY_QUEUE, delay)
		}

		if enqueueErr != nil {
			ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(fmt.Errorf("failed to enqueue event: %v", enqueueErr)))
			return
		}

		// Store the new task ID in the database
		err = server.Repo.StoreTaskID(models.EventTask{
			EventID:    int(event.ID),
			TaskID:     taskID,
			IsCanceled: false,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, event)
}

// @Summary Edit API Event
// @Description Edit an API Event
// @Tags Events
// @Accept json
// @Produce json
// @Param request body models.CreateAPIEventRequestSwagger true "API Event data"
// @Success 200 {object} models.SwaggerEventDto
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /events/api/{id} [put]
// @Security BearerAuth
func (server *Server) EditAPIEvent(ctx *gin.Context) {
	// Parse URI for event ID
	var mod models.GetEventRequest
	if err := ctx.ShouldBindUri(&mod); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	var req models.CreateAPIEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(constants.AUTHORIZATION_PAYLOAD_KEY).(*token.Payload)

	// Validation for API Trigger
	err := validator.ValidateAPIEventRequest(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// Fetch the existing event
	_, err = server.Repo.GetEvent(mod.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, constants.ErrorResponse(err))
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

	event, err := server.Repo.UpdateEvent(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, event)
}
