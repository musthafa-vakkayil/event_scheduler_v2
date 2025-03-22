package repo

import (
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"gorm.io/gorm"
)

func DeleteUser(db *gorm.DB, username string) error {
	err := db.Where("username = ?", username).Delete(&models.User{}).Error
	return err
}

func CreateUser(db *gorm.DB, user models.User) (models.User, error) {
	err := db.Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUser(db *gorm.DB, username string) (models.User, error) {
	var user models.User
	err := db.Table("users").Where("username=?", username).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
