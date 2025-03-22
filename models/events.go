package models

import (
	"time"

	"gorm.io/datatypes"
)

type Event struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	Name           string         `json:"name" gorm:"not null"`
	Type           string         `json:"type" gorm:"not null"`
	ApiEndPoint    string         `json:"api_endpoint"`
	ApiMethod      string         `json:"api_method"`
	ApiRequestBody datatypes.JSON `json:"api_request_body"`
	CreatedBy      string         `json:"created_by"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;"`
	ExecutedAt     time.Time      `gorm:"autoUpdateTime;"`
}

type CreateEventRequest struct {
	Name        string         `json:"name" binding:"required"`
	Type        string         `json:"type" binding:"required,oneof=API"`
	ApiEndpoint string         `json:"api_endpoint" binding:"required"`
	ApiMethod   string         `json:"api_method" binding:"required,oneof=GET POST PUT PATCH DELETE"`
	ApiPayload  datatypes.JSON `json:"api_payload"`
}

type ListEventsRequest struct {
	PageNumber  int  `form:"pageNumber" binding:"required,min=1"`
	PageSize    int  `form:"pageSize" binding:"required,min=1"`
	CreatedByMe bool `form:"createdByMe"`
}

type GetEventRequest struct {
	ID int `uri:"id" binding:"required"`
}
