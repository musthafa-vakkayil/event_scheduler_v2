package repo

import (
	"fmt"

	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"gorm.io/gorm"
)

func (r *Repo) CreateEvent(event models.Event) (models.Event, error) {
	err := r.DB.Create(&event).Error
	if err != nil {
		return models.Event{}, err
	}
	return event, nil
}

func (r *Repo) GetEvent(id int) (models.Event, error) {
	var event models.Event

	// Use `Take()` or `First()` to return a single record and handle not found error
	err := r.DB.Table("events").Where("id = ?", id).Take(&event).Error

	// Handle the "record not found" error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Event{}, fmt.Errorf("event not found")
		}
		return models.Event{}, err
	}

	return event, nil
}

func (r *Repo) ListEvents(limit, offset int) ([]models.Event, error) {
	var events []models.Event

	// Use GORM query with limit and offset
	err := r.DB.Order("id").Limit(limit).Offset(offset).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}
