package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/mocks"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListLogsCacheMiss(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.Repository)
	mockCache := new(mocks.Cache)

	// Mock cache miss
	mockCache.On("Get", mock.Anything, false, true, 1, 10).Return(nil, nil)

	// Mock DB fallback
	mockRepo.On("ListLogs", true, false, 10, 0).Return([]models.LogsResponse{
		{ID: 2, ExecutedBy: "DB Log"},
	}, nil)

	mockCache.On("Set", mock.Anything, mock.Anything, false, true, 1, 10, mock.Anything).Return(nil)

	server := SetupTestServer(t, mockRepo, mockCache)

	// Perform the request
	req, _ := http.NewRequest(http.MethodGet, "/logs?pageNumber=1&pageSize=10&onlyActive=true", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	server.ListLogs(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "DB Log")

	// Ensure DB and cache were called
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestListLogsCacheHit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.Repository)
	mockCache := new(mocks.Cache)

	// Mock cache hit
	mockCache.On("Get", mock.Anything, false, false, 1, 10).Return([]models.LogsResponse{
		{ID: 1, ExecutedBy: "Cache Hit Log"},
	}, nil)

	server := SetupTestServer(t, mockRepo, mockCache)

	// Perform the request
	req, _ := http.NewRequest(http.MethodGet, "/logs?pageNumber=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	server.ListLogs(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Cache Hit Log")

	// Ensure cache was called
	mockCache.AssertExpectations(t)
}

func TestListLogsCacheError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.Repository)
	mockCache := new(mocks.Cache)

	// Mock cache error
	mockCache.On("Get", mock.Anything, false, true, 1, 10).Return(nil, errors.New("redis error"))

	// Mock DB fallback
	mockRepo.On("ListLogs", true, false, 10, 0).Return([]models.LogsResponse{
		{ID: 3, ExecutedBy: "DB Fallback Log"},
	}, nil)

	// Cache set after DB fallback
	mockCache.On("Set", mock.Anything, mock.Anything, false, true, 1, 10, mock.Anything).Return(nil)

	server := SetupTestServer(t, mockRepo, mockCache)

	// Perform the request
	req, _ := http.NewRequest(http.MethodGet, "/logs?pageNumber=1&pageSize=10&onlyActive=true", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	server.ListLogs(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "DB Fallback Log")

	// Ensure DB was called after cache error
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestListLogsInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.Repository)
	mockCache := new(mocks.Cache)

	// Mock cache miss
	mockCache.On("Get", mock.Anything, false, true, 1, 10).Return(nil, nil)

	// Mock DB fallback
	mockRepo.On("ListLogs", true, false, 10, 0).Return([]models.LogsResponse{}, errors.New("internal_error"))

	server := SetupTestServer(t, mockRepo, mockCache)

	// Perform the request
	req, _ := http.NewRequest(http.MethodGet, "/logs?pageNumber=1&pageSize=10&onlyActive=true", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	server.ListLogs(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockCache.AssertExpectations(t)
}

func TestListLogsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.Repository)
	mockCache := new(mocks.Cache)

	server := SetupTestServer(t, mockRepo, mockCache)

	// Perform the request
	req, _ := http.NewRequest(http.MethodGet, "/logs", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	server.ListLogs(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockCache.AssertExpectations(t)
}
