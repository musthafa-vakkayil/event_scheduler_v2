package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(constants.AUTHORIZATION_HEADER_KEY)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is missing")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, constants.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, constants.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != constants.AUTHORIZATION_TYPE_BEARER {
			err := errors.New("invalid authorization type")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, constants.ErrorResponse(err))
			return
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, constants.ErrorResponse(err))
			return
		}

		ctx.Set(constants.AUTHORIZATION_PAYLOAD_KEY, payload)
		ctx.Next()
	}
}
