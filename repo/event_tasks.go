package repo

import (
	"fmt"

	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"gorm.io/gorm"
)

func (r *Repo) StoreTaskID(et models.EventTask) error {
	return r.DB.Create(&et).Error
}

// Cancel tasks for the given event
func (r *Repo) CancelTasksForEvent(eventID int) error {
	return r.DB.Model(&models.EventTask{}).
		Where("event_id = ?", eventID).
		Update("is_canceled", true).Error
}

// Check if a task is canceled
func (r *Repo) IsTaskCanceled(taskID string) (bool, error) {
	var task models.EventTask
	err := r.DB.Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return task.IsCanceled, nil
}

func (r *Repo) DeleteTaskEvent(taskID string) error {
	// Delete the task entry by task_id
	if err := r.DB.Table("event_tasks").
		Where("task_id = ?", taskID).
		Delete(&models.EventTask{}).Error; err != nil {
		return fmt.Errorf("failed to delete task %s: %v", taskID, err)
	}

	return nil
}
