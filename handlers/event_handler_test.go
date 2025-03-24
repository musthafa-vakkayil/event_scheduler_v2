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
)

func TestGetEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.Repository)
	mockCache := new(mocks.Cache)
	mockManager := new(mocks.TaskManager)

	t.Run("Event found", func(t *testing.T) {
		mockRepo.On("GetEvent", 1).Return(models.Event{ID: 1}, nil).Once()

		server := SetupTestServer(t, mockRepo, mockCache, mockManager)

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/events/1", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Add the URI parameter manually
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		server.GetEvent(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "1")
		mockCache.AssertExpectations(t)
	})
	t.Run("bad request", func(t *testing.T) {

		server := SetupTestServer(t, mockRepo, mockCache, mockManager)

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/events/0", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Add the URI parameter manually
		ctx.Params = gin.Params{
			{Key: "id", Value: "0"},
		}

		server.GetEvent(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockCache.AssertExpectations(t)
	})
	t.Run("Not found", func(t *testing.T) {
		mockRepo.On("GetEvent", 1).Return(models.Event{}, errors.New("event not found")).Once()

		server := SetupTestServer(t, mockRepo, mockCache, mockManager)

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/events/1", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Add the URI parameter manually
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		server.GetEvent(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockCache.AssertExpectations(t)
	})
	t.Run("Internal Error", func(t *testing.T) {
		mockRepo.On("GetEvent", 1).Return(models.Event{}, errors.New("internal error")).Once()

		server := SetupTestServer(t, mockRepo, mockCache, mockManager)

		// Perform the request
		req, _ := http.NewRequest(http.MethodGet, "/events/1", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Add the URI parameter manually
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		server.GetEvent(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockCache.AssertExpectations(t)
	})
}

func TestDeleteEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(mocks.Repository)
	mockCache := new(mocks.Cache)
	mockManager := new(mocks.TaskManager)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("DeleteEvent", 1).Return(nil).Once()

		server := SetupTestServer(t, mockRepo, mockCache, mockManager)

		// Perform the request
		req, _ := http.NewRequest(http.MethodDelete, "/events/1", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Add the URI parameter manually
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		server.DeleteEvent(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "OK")
		mockCache.AssertExpectations(t)
	})
	t.Run("bad request", func(t *testing.T) {

		server := SetupTestServer(t, mockRepo, mockCache, mockManager)

		// Perform the request
		req, _ := http.NewRequest(http.MethodDelete, "/events/0", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Add the URI parameter manually
		ctx.Params = gin.Params{
			{Key: "id", Value: "0"},
		}

		server.DeleteEvent(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockCache.AssertExpectations(t)
	})
	t.Run("Internal Error", func(t *testing.T) {
		mockRepo.On("DeleteEvent", 1).Return(errors.New("internal error")).Once()

		server := SetupTestServer(t, mockRepo, mockCache, mockManager)

		// Perform the request
		req, _ := http.NewRequest(http.MethodDelete, "/events/1", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Add the URI parameter manually
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		server.DeleteEvent(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockCache.AssertExpectations(t)
	})
}
