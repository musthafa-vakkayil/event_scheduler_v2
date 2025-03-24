package models

import "time"

type Log struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	EventId    int64     `json:"event_id" gorm:"not null"`
	ExecutedOn time.Time `json:"executed_on"`
	Status     string    `json:"status" gorm:"not null"`
	IsArchived bool      `json:"is_archived"`
	LogType    string    `json:"log_type" gorm:"not null"`
}

type LogsResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CreatedBy  string    `json:"created_by"`
	ExecutedOn time.Time `json:"executed_on"`
	Status     string    `json:"status"`
	IsArchived bool      `json:"is_archived"`
	LogType    string    `json:"log_type"`
}

type ListLogsRequest struct {
	PageNumber   int  `form:"pageNumber" binding:"required,min=1"`
	PageSize     int  `form:"pageSize" binding:"required,min=1"`
	OnlyActive   bool `form:"onlyActive"`
	OnlyArchived bool `form:"onlyArchived"`
}
