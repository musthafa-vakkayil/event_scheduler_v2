package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/middleware"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(constants.AUTHORIZATION_HEADER_KEY, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test token maker
	secretKey := "Af3jTWzfhGnzJvZPe63Mc8DAVNTlwkXG"
	tokenMaker, err := token.NewJWTMaker(secretKey)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, tokenMaker, constants.AUTHORIZATION_TYPE_BEARER, "user", 2*time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {}, // No auth added
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, tokenMaker, "invalid", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, tokenMaker, constants.AUTHORIZATION_TYPE_BEARER, "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.Default()

			router.GET(
				"/auth",
				middleware.AuthMiddleware(tokenMaker), // Use middleware directly
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{"message": "success"})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)

			tc.setupAuth(t, request)
			router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}
