package repo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"gorm.io/gorm"
)

func (r *Repo) CreateSession(session models.Session) (models.Session, error) {
	err := r.DB.Create(&session).Error
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}

func (r *Repo) GetSession(id uuid.UUID) (models.Session, error) {
	var session models.Session
	err := r.DB.Table("sessions").Where("id = ?", id).Take(&session).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Session{}, fmt.Errorf("session not found")
		}
		return models.Session{}, err
	}

	return session, nil
}
