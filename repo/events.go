package repo

import (
	"fmt"
	"time"

	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"gorm.io/datatypes"
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
	var events models.ListEventsResponse

	// Use GORM query with limit and offset
	err := r.DB.Order("id").Limit(limit).Offset(offset).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *Repo) ExecuteEvent(username string, eventID int64, logType string, status string, apiPayload datatypes.JSON) (int64, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// Update executed_at in events table
	if err := tx.Model(&models.Event{}).Where("id = ?", eventID).Update("executed_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Create log entry
	log := models.Log{
		EventId:     eventID,
		TriggeredOn: time.Now(),
		Status:      status,
		IsArchived:  false,
		LogType:     logType,
		ExecutedBy:  username,
		ApiPayload:  apiPayload,
	}

	if err := tx.Create(&log).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return log.ID, tx.Commit().Error
}

func (r *Repo) DeleteEvent(eventID int) error {
	// Start a transaction
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete the event
	if err := tx.Where("id = ?", eventID).Delete(&models.Event{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete event: %v", err)
	}

	// Mark related tasks as canceled
	if err := tx.Table("event_tasks").Where("event_id = ?", eventID).Update("is_canceled", true).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to mark tasks as canceled: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *Repo) UpdateEvent(event models.Event) (models.Event, error) {
	err := r.DB.Model(&models.Event{}).
		Where("id = ?", event.ID).
		Updates(event).Error
	if err != nil {
		return models.Event{}, err
	}

	return event, nil
}
