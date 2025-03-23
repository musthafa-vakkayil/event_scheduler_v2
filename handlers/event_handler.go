package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/tasks"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
)

func (server *Server) CreateEvent(ctx *gin.Context) {
	var req models.CreateEventRequest
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

func (server *Server) ExecuteAPIEvent(ctx *gin.Context) {
	var req models.GetEventRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	event, err := server.Repo.GetEvent(req.ID)
	if err != nil {
		if err.Error() == "event not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	if event.Type != "API" {
		err := errors.New("cannot excute events other than API")
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	logId, err := server.Repo.ExecuteEvent(event.ID, "SUCCESS")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	// Create an archive task
	archiveTask, err := tasks.NewArchiveLogTask(int(logId))
	if err != nil {
		log.Fatal("Failed to create archive task:", err)
	}

	// Enqueue the archive task with a 2-minute delay (for testing)
	_, err = server.QueueClient.Enqueue(archiveTask, asynq.Queue("low"), asynq.ProcessIn(server.Config.LogArchiveDuration))
	if err != nil {
		log.Fatal("Failed to enqueue archive task:", err)
	}

	// Invalidate cache
	cacheKey := "event-scheduler:log-*"       // Wildcard pattern
	server.Cache.DeletePattern(ctx, cacheKey) // New cache invalidation method

	ctx.JSON(http.StatusOK, "OK")
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
