package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
)

func (server *Server) ListLogs(ctx *gin.Context) {
	var req models.ListLogsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	logs, err := server.Repo.ListLogs(req.OnlyActive, req.OnlyArchived, req.PageSize, (req.PageNumber-1)*req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, logs)
}
