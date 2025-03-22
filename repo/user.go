package repo

import (
	"fmt"

	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"gorm.io/gorm"
)

// DeleteUser deletes a user by username
func (r *Repo) DeleteUser(username string) error {
	return r.DB.Where("username = ?", username).Delete(&models.User{}).Error
}

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
