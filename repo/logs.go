package repo

import (
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
)

func (r *Repo) ListLogs(onlyActive bool, onlyArchived bool, limit int, offset int) ([]models.LogsResponse, error) {
	var logs []models.LogsResponse

	query := r.DB.Table("logs").
		Select("logs.id, logs.executed_by, logs.triggered_on, logs.status, logs.is_archived, logs.log_type, logs.api_payload")

	if onlyActive {
		query = query.Where("logs.is_archived = false")
	} else if onlyArchived {
		query = query.Where("logs.is_archived = true")
	}

	query = query.Order("logs.id DESC").Limit(limit).Offset(offset)

	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *Repo) MarkLogAsArchived(logID int) error {
	return r.DB.Model(&models.Log{}).Where("id = ?", logID).Update("is_archived", true).Error
}

func (r *Repo) DeleteLog(logID int) error {
	return r.DB.Delete(&models.Log{}, logID).Error
}

func (r *Repo) CreateLog(log models.Log) (models.Log, error) {
	if err := r.DB.Create(&log).Error; err != nil {
		return models.Log{}, err
	}

	return log, nil
}
