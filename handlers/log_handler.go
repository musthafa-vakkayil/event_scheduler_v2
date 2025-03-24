package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
)

// @Summary List Triggers
// @Description List all Executed Events
// @Tags Logs
// @Produce json
// @Param pageNumber query int true "Page Number" minimum(1)
// @Param pageSize query int true "Page Size" minimum(1)
// @Param onlyActive query bool false "Only Active Events"
// @Param onlyArchived query bool false "Only Archived Events"
// @Success 200 {object} models.SwaggerLogsResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedRequestResponse
// @Failure 500 {object} models.InternalServerErrorResponse
// @Router /logs [get]
// @Security BearerAuth
func (server *Server) ListLogs(ctx *gin.Context) {
	var req models.ListLogsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// Try to get data from cache first
	logs, err := server.Cache.Get(ctx, req.OnlyArchived, req.OnlyActive, req.PageNumber, req.PageSize)
	if err != nil {
		log.Println("❌ Redis cache error, falling back to DB:", err)
	}

	if logs != nil {
		log.Println("✅ Cache hit")
		ctx.JSON(http.StatusOK, logs)
		return
	}

	log.Println("❌ Cache miss. Fetching from DB...")

	logs, err = server.Repo.ListLogs(req.OnlyActive, req.OnlyArchived, req.PageSize, (req.PageNumber-1)*req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	// Set the result in cache
	_ = server.Cache.Set(ctx, logs, req.OnlyArchived, req.OnlyActive, req.PageNumber, req.PageSize, server.Config.CacheExpiration)

	ctx.JSON(http.StatusOK, logs)
}
