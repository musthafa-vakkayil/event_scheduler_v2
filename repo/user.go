package repo

import (
	"fmt"

	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"gorm.io/gorm"
)

// CreateUser inserts a new user into the database
func (r *Repo) CreateUser(user models.User) (models.User, error) {
	err := r.DB.Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// GetUser retrieves a user by username
func (r *Repo) GetUser(username string) (models.User, error) {
	var user models.User
	err := r.DB.Table("users").Where("username = ?", username).Take(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *Repo) DeleteUser(username string) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete sessions created by user
	if err := tx.Where("username = ?", username).Delete(&models.Session{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete events created by user
	if err := tx.Where("created_by = ?", username).Delete(&models.Event{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete user
	if err := tx.Where("username = ?", username).Delete(&models.User{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
