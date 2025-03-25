package models

import (
	"time"

	"gorm.io/datatypes"
)

// @Description Event object used for output
type Event struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	Name           string         `json:"name" gorm:"not null"`
	Type           string         `json:"type" gorm:"not null"`
	ApiEndPoint    string         `json:"api_endpoint"`
	ApiMethod      string         `json:"api_method"`
	ApiRequestBody datatypes.JSON `json:"api_request_body"`
	RunAt          *time.Time     `json:"run_at"`
	AfterXMins     int            `json:"after_x_mins"`
	Interval       int            `json:"interval"`
	IsRecurring    bool           `json:"is_recurring"`
	CreatedBy      string         `json:"created_by"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;"`
	ExecutedAt     *time.Time     `json:"executed_at"`
}
type CreateAPIEventRequest struct {
	Name        string         `json:"name" binding:"required" example:"API Event 1"`
	Type        string         `json:"type" binding:"required,oneof=API" example:"API"`
	ApiEndpoint string         `json:"api_endpoint" binding:"required" example:"https://api.example.com/v1/users"`
	ApiMethod   string         `json:"api_method" binding:"required,oneof=GET POST PUT PATCH DELETE" example:"GET"`
	ApiPayload  datatypes.JSON `json:"api_payload" example:"{\"key\": \"value\"}"`
}

type ListEventsRequest struct {
	PageNumber int `form:"pageNumber" binding:"required,min=1"`
	PageSize   int `form:"pageSize" binding:"required,min=1"`
}

// @Description ListEventsResponse object used for output
type ListEventsResponse []Event

type GetEventRequest struct {
	ID int `uri:"id" binding:"required"`
}

// @Description CreateScheduledEventRequest object used for input
type CreateScheduledEventRequest struct {
	Name         string `json:"name" binding:"required" example:"Scheduled Event 1"`
	Type         string `json:"type" binding:"required,oneof=SCHEDULED" example:"SCHEDULED"`
	RunAtDate    string `json:"run_at_this_date" example:"2025-03-25T12:15:00+05:30"`
	RunAfterMins int    `json:"run_after_x_mins" example:"5"`
	IntervalMins int    `json:"repeat_after_x_mins" example:"5"`
	IsRecurring  bool   `json:"is_recurring" example:"false"`
}

// @Description CreateTestScheduledEventRequest object used for input
type CreateTestScheduledEventRequest struct {
	RunAtDate time.Time `json:"run_at_this_date" binding:"required" example:"2021-08-01T00:00:00Z"`
}

type EventTask struct {
	EventID    int    `json:"event_id"`
	TaskID     string `json:"task_id"`
	IsCanceled bool   `json:"is_canceled"`
}
