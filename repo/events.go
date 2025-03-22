package repo

import (
	"fmt"

	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"gorm.io/gorm"
)

func CreateEvent(db *gorm.DB, event models.Event) (models.Event, error) {
	err := db.Create(&event).Error
	if err != nil {
		return models.Event{}, err
	}
	return event, nil
}

func GetEvent(db *gorm.DB, id int64) (models.Event, error) {
	var event models.Event

	// Use `Take()` or `First()` to return a single record and handle not found error
	err := db.Table("events").Where("id = ?", id).Take(&event).Error

	// Handle the "record not found" error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Event{}, fmt.Errorf("event not found")
		}
		return models.Event{}, err
	}

	return event, nil
}

func ListEvents(db *gorm.DB, limit, offset int) ([]models.Event, error) {
	var events []models.Event

	// Use GORM query with limit and offset
	err := db.Order("id").Limit(limit).Offset(offset).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}
