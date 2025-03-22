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
