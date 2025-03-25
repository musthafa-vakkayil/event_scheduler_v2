package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
	"github.com/musthafa-vakkayil/event_scheduler_v2/utils"
	"github.com/musthafa-vakkayil/event_scheduler_v2/validator"
)

// @Summary Create Test API Event
// @Description Create a new API Event without saving to DB
// @Tags Test Events
// @Accept json
// @Produce json
// @Param request body models.CreateAPIEventRequestSwagger true "API Event data"
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /test/events/api [post]
// @Security BearerAuth
func (server *Server) CreateTestAPIEvent(ctx *gin.Context) {
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

	// Use the utility function to make the API call
	payload := []byte(req.ApiPayload)
	resp, err := utils.MakeAPICall(req.ApiMethod, req.ApiEndpoint, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err := errors.New("failed to read response body")
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	args := models.Log{
		TriggeredOn: time.Now(),
		Status:      resp.Status,
		IsArchived:  false,
		ExecutedBy:  authPayload.Username,
		ApiPayload:  req.ApiPayload,
		LogType:     "TEST_EVENT",
	}

	logData, err := server.Repo.CreateLog(args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
	}

	err = server.TaskManager.EnqueueTask(ctx, int(logData.ID), constants.TASK_ARCHIVE_LOG, constants.LOW_PRIORITY_QUEUE, server.Config.LogArchiveDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	// Invalidate cache
	cacheKey := "event-scheduler:log-*"       // Wildcard pattern
	server.Cache.DeletePattern(ctx, cacheKey) // New cache invalidation method

	ctx.JSON(http.StatusOK, json.RawMessage(body))
}

// @Summary Create Test Scheduled Event
// @Description Schedule a one time event
// @Tags Test Events
// @Accept json
// @Produce json
// @Param request body models.CreateTestEventRequest true "Test Event data"
// @Success 200 {object} models.EmptyResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /test/events/schedule [post]
// @Security BearerAuth
func (server *Server) CreateTestScheduledEvent(ctx *gin.Context) {
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

	randomId := utils.RandomInt(100, 1000)
	// Enqueue the task based on RunAt or RunAfter
	var enqueueErr error

	// Parse time using the shared function
	runTime, err := utils.ParseAndValidateTime(req.RunAtDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// If `RunAt` is provided → Enqueue with `ProcessAt`
	if runTime != nil {
		enqueueErr = server.TaskManager.EnqueueTaskAt(ctx, randomId, constants.TASK_TEST_EVENT, constants.DEFAULT_PRIORITY_QUEUE, *runTime)
	} else {
		// If `RunAfterMins` is provided → Enqueue with `ProcessIn`
		delay := time.Duration(req.RunAfterMins) * time.Minute
		enqueueErr = server.TaskManager.EnqueueTaskIn(ctx, randomId, constants.TASK_TEST_EVENT, constants.DEFAULT_PRIORITY_QUEUE, delay)
	}

	if enqueueErr != nil {
		err := fmt.Errorf("failed to enqueue event: %v", enqueueErr)
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	response := models.Event{
		ID:         randomId,
		RunAt:      nil,
		AfterXMins: req.RunAfterMins,
	}

	ctx.JSON(http.StatusOK, response)
}
