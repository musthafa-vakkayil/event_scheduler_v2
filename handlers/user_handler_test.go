package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/musthafa-vakkayil/event_scheduler_v2/constants"
	"github.com/musthafa-vakkayil/event_scheduler_v2/mocks"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/musthafa-vakkayil/event_scheduler_v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	server := SetupTestServer(t, mockRepo)

	t.Run("Success", func(t *testing.T) {
		// Define expected user
		reqBody := models.CreateUserRequest{
			Username: "testuser",
			Password: "password123",
			FullName: "Test User",
			Email:    "test@example.com",
		}

		hashedPassword := "$2a$12$abcdefghijklmnopqrstuv" // Fake hashed password

		mockRepo.On("CreateUser", mock.Anything).Return(models.User{
			Username:       reqBody.Username,
			HashedPassword: hashedPassword,
			FullName:       reqBody.FullName,
			Email:          reqBody.Email,
		}, nil).Once()

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Bad_Request", func(t *testing.T) {
		// Define expected user
		reqBody := models.CreateUserRequest{
			Username: "",
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal_Server_Error", func(t *testing.T) {
		// Define expected user
		reqBody := models.CreateUserRequest{
			Username: "testuser",
			Password: "password123",
			FullName: "Test User",
			Email:    "test@example.com",
		}

		mockRepo.On("CreateUser", mock.Anything).Return(models.User{}, errors.New("internal_error")).Once()

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	server := SetupTestServer(t, mockRepo)
	token, payload, err := server.TokenMaker.CreateToken("admin", server.Config.TokenDuration)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)
	authToken := fmt.Sprintf("%v %v", constants.AUTHORIZATION_TYPE_BEARER, token)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetUser", mock.Anything).Return(models.User{
			Username:       utils.RandomString(6),
			HashedPassword: "test123",
			FullName:       utils.RandomString(6),
			Email:          utils.RandomEmail(),
		}, nil).Once()

		req, _ := http.NewRequest("GET", "/users/testuser", nil)

		req.Header.Add(constants.AUTHORIZATION_HEADER_KEY, authToken)

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Bad Request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/12@", nil)

		req.Header.Add(constants.AUTHORIZATION_HEADER_KEY, authToken)

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("GetUser", mock.Anything).Return(models.User{}, errors.New("user not found")).Once()

		req, _ := http.NewRequest("GET", "/users/testuser", nil)

		req.Header.Add(constants.AUTHORIZATION_HEADER_KEY, authToken)

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Internal Server Error", func(t *testing.T) {
		// Define expected user
		reqBody := models.GetUserRequest{
			Username: "testuser",
		}

		mockRepo.On("GetUser", reqBody.Username).Return(models.User{}, errors.New("internal_error")).Once()

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("GET", "/users/testuser", bytes.NewBuffer(body))

		req.Header.Add(constants.AUTHORIZATION_HEADER_KEY, authToken)

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	server := SetupTestServer(t, mockRepo)

	t.Run("Success", func(t *testing.T) {
		// Define expected user
		password := utils.RandomString(6)
		hashedPassword, err := utils.HashPassword(password)
		assert.NoError(t, err)

		reqBody := models.LoginRequest{
			Username: "testuser",
			Password: password,
		}

		mockRepo.On("GetUser", mock.Anything).Return(models.User{
			Username:       reqBody.Username,
			HashedPassword: hashedPassword,
			FullName:       utils.RandomString(6),
			Email:          utils.RandomEmail(),
		}, nil).Once()

		mockRepo.On("CreateSession", mock.Anything).Return(models.Session{}, nil)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Invalid password", func(t *testing.T) {
		// Define expected user
		password := utils.RandomString(6)
		hashedPassword, err := utils.HashPassword(utils.RandomString(6))
		assert.NoError(t, err)

		reqBody := models.LoginRequest{
			Username: "testuser",
			Password: password,
		}

		mockRepo.On("GetUser", mock.Anything).Return(models.User{
			Username:       reqBody.Username,
			HashedPassword: hashedPassword,
			FullName:       utils.RandomString(6),
			Email:          utils.RandomEmail(),
		}, nil).Once()

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Bad Request", func(t *testing.T) {
		reqBody := models.LoginRequest{
			Username: "testuser",
		}

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockRepo.AssertExpectations(t)
	})
	t.Run("User Not found", func(t *testing.T) {
		password := utils.RandomString(6)

		reqBody := models.LoginRequest{
			Username: "testuser",
			Password: password,
		}

		mockRepo.On("GetUser", mock.Anything).Return(models.User{}, errors.New("user not found")).Once()

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))

		resp := httptest.NewRecorder()
		server.Router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockRepo.AssertExpectations(t)
	})
}
