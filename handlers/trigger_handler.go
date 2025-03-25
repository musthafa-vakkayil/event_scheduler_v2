package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
	"github.com/musthafa-vakkayil/event_scheduler_v2/utils"
)

// ExecuteAPIEvent godoc
// @Summary Execute an API event
// @Description Execute an API event by making an API call
// @Tags Events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} models.EmptyResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /events/{id}/execute [get]
// @Security BearerAuth
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

	authPayload := ctx.MustGet(constants.AUTHORIZATION_PAYLOAD_KEY).(*token.Payload)

	logId, err := server.Repo.ExecuteEvent(authPayload.Username, event.ID, "API_EVENT", resp.Status, event.ApiRequestBody)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	err = server.TaskManager.EnqueueTask(ctx, int(logId), constants.TASK_ARCHIVE_LOG, constants.LOW_PRIORITY_QUEUE, server.Config.LogArchiveDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	// Invalidate cache
	cacheKey := "event-scheduler:log-*"       // Wildcard pattern
	server.Cache.DeletePattern(ctx, cacheKey) // New cache invalidation method

	ctx.JSON(http.StatusOK, json.RawMessage(body))
}
