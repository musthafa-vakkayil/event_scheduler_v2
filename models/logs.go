package models

import (
	"time"

	"gorm.io/datatypes"
)

type Log struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	EventId     int64          `json:"event_id"`
	TriggeredOn time.Time      `json:"triggered_on"`
	Status      string         `json:"status" gorm:"not null"`
	IsArchived  bool           `json:"is_archived"`
	LogType     string         `json:"log_type" gorm:"not null"`
	ExecutedBy  string         `json:"executed_by"`
	ApiPayload  datatypes.JSON `json:"api_payload"`
}

type LogsResponse struct {
	ID          int64          `json:"id"`
	ExecutedBy  string         `json:"executed_by"`
	TriggeredOn time.Time      `json:"triggered_on"`
	Status      string         `json:"status"`
	IsArchived  bool           `json:"is_archived"`
	LogType     string         `json:"log_type"`
	ApiPayload  datatypes.JSON `json:"api_payload"`
}

type ListLogsRequest struct {
	PageNumber   int  `form:"pageNumber" binding:"required,min=1"`
	PageSize     int  `form:"pageSize" binding:"required,min=1"`
	OnlyActive   bool `form:"onlyActive"`
	OnlyArchived bool `form:"onlyArchived"`
}
