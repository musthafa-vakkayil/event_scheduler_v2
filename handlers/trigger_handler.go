package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/tasks"
	"github.com/musthafa-vakkayil/event_scheduler_v2/utils"
)

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

	// Use the utility function to make the API call
	payload := []byte(event.ApiRequestBody)
	resp, err := utils.MakeAPICall(event.ApiMethod, event.ApiEndPoint, payload)
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

	logId, err := server.Repo.ExecuteEvent(event.ID, resp.Status)
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

	ctx.JSON(http.StatusOK, json.RawMessage(body))
}
