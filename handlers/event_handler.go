package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	db "github.com/musthafa-vakkayil/event_scheduler_v2/db/sqlc"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
	"github.com/sqlc-dev/pqtype"
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

	arg := db.CreateEventParams{
		Name:        req.Name,
		Type:        db.EventTypes(req.Type),
		ApiEndPoint: req.ApiEndpoint,
		ApiMethod: sql.NullString{
			String: req.ApiMethod,
			Valid:  true,
		},
		CreatedBy: authPayload.Username,
	}

	// Validation for API Trigger
	if req.Type == "API" && req.ApiMethod != "GET" && req.ApiMethod != "DELETE" && req.ApiPayload == nil {
		err := errors.New("api payload is required for methods other than GET and delete")
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	// Handle ApiPayload safely
	if len(req.ApiPayload) > 0 {
		arg.ApiRequestBody = pqtype.NullRawMessage{
			RawMessage: []byte(req.ApiPayload),
			Valid:      true,
		}
	} else {
		arg.ApiRequestBody = pqtype.NullRawMessage{
			Valid: false,
		}
	}

	event, err := server.Store.CreateEvent(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, event)
}

type ListEventsRequest struct {
	PageNumber  int32 `form:"pageNumber" binding:"required,min=1"`
	PageSize    int32 `form:"pageSize" binding:"required,min=1"`
	CreatedByMe bool  `form:"createdByMe"`
}

func (server *Server) ListEvents(ctx *gin.Context) {
	var req ListEventsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(constants.AUTHORIZATION_PAYLOAD_KEY).(*token.Payload)

	var events []db.Event
	var err error

	if req.CreatedByMe {
		args := db.ListUserEventsParams{
			Limit:     req.PageSize,
			Offset:    (req.PageNumber - 1) * req.PageSize,
			CreatedBy: authPayload.Username,
		}
		events, err = server.Store.ListUserEvents(ctx, args)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
			return
		}
	} else {
		args := db.ListEventsParams{
			Limit:  req.PageSize,
			Offset: (req.PageNumber - 1) * req.PageSize,
		}
		events, err = server.Store.ListEvents(ctx, args)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, constants.ErrorResponse(err))
			return
		}
	}

	fmt.Println(events[0].ApiMethod)

	ctx.JSON(http.StatusOK, events)
}
