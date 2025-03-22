package models

import "time"

type Log struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	EventId    int64     `json:"event_id" gorm:"not null"`
	ExecutedOn time.Time `json:"executed_on"`
	Status     string    `json:"status" gorm:"not null"`
	IsArchived bool      `json:"is_archived"`
}
