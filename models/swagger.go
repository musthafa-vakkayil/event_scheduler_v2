package models

import (
	"time"
)

// ErrorResponse represents a standard error response
// @Description Standard error response object
type ErrorResponse struct {
	Error string `json:"error"`
}

// InternalServerErrorResponse represents a internal error response
// @Description Standard error response object
type InternalServerErrorResponse struct {
	Error string `json:"internal server error"`
}

// BadRequestResponse represents a internal error response
// @Description Standard error response object
type BadRequestResponse struct {
	Error string `json:"bad request"`
}

// UnauthorizedRequestResponse represents a internal error response
// @Description Standard error response object
type UnauthorizedRequestResponse struct {
	Error string `json:"unauthorized request"`
}

// EmptyResponse represents a internal error response
// @Description Standard error response object
type EmptyResponse struct {
	Error string `json:"ok"`
}

// NotFoundResponse represents a internal error response
// @Description Standard error response object
type NotFoundResponse struct {
	Error string `json:"not found"`
}

// @Description CreateAPIEventRequestSwagger used only for Swagger docs
type CreateAPIEventRequestSwagger struct {
	Name        string            `json:"name" binding:"required" example:"API Event 1"`
	Type        string            `json:"type" binding:"required,oneof=API" example:"API"`
	ApiEndpoint string            `json:"api_endpoint" binding:"required" example:"https://api.example.com/v1/users"`
	ApiMethod   string            `json:"api_method" binding:"required,oneof=GET POST PUT PATCH DELETE" example:"GET"`
	ApiPayload  map[string]string `json:"api_payload"`
}

// @Description SwaggerEventDto used only for Swagger docs
type SwaggerEventDto struct {
	ID             int64             `json:"id"`
	Name           string            `json:"name"`
	Type           string            `json:"type"`
	ApiEndPoint    string            `json:"api_endpoint"`
	ApiMethod      string            `json:"api_method"`
	ApiRequestBody map[string]string `json:"api_request_body"`
	RunAt          time.Time         `json:"run_at"`
	AfterXMins     int               `json:"after_x_mins"`
	Interval       int               `json:"interval"`
	IsRecurring    bool              `json:"is_recurring"`
	CreatedBy      string            `json:"created_by"`
	CreatedAt      time.Time         `gorm:"autoCreateTime;"`
	ExecutedAt     time.Time         `json:"executed_at"`
}

// @Description SwaggerListEventResponse used only for Swagger docs
type SwaggerListEventResponse []SwaggerEventDto

type SwaggerLogs struct {
	ID          int64             `json:"id"`
	ExecutedBy  string            `json:"executed_by"`
	TriggeredOn time.Time         `json:"triggered_on"`
	Status      string            `json:"status"`
	IsArchived  bool              `json:"is_archived"`
	LogType     string            `json:"log_type"`
	ApiPayload  map[string]string `json:"api_payload"`
}

// @Description SwaggerLogsResponse used only for Swagger docs
type SwaggerLogsResponse []SwaggerLogs

// @Description CreateTestEventRequest object used for swagger
type CreateTestEventRequest struct {
	Name         string     `json:"name" binding:"required" example:"Test Event 1"`
	Type         string     `json:"type" binding:"required,oneof=SCHEDULED" example:"SCHEDULED"`
	RunAtDate    *time.Time `json:"run_at_this_date" example:"2021-08-01T00:00:00Z"`
	RunAfterMins int        `json:"run_after_x_mins" example:"5"`
}
