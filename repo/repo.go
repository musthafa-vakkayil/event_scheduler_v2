package repo

import (
	"github.com/google/uuid"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user models.User) (models.User, error)
	GetUser(username string) (models.User, error)
	CreateEvent(event models.Event) (models.Event, error)
	GetEvent(id int) (models.Event, error)
	ListEvents(limit, offset int) ([]models.Event, error)
	ExecuteEvent(username string, eventID int64, logType string, status string, apiPayload datatypes.JSON) (int64, error)
	DeleteEvent(eventID int) error
	DeleteUser(username string) error
	ListLogs(onlyActive bool, onlyArchived bool, limit int, offset int) ([]models.LogsResponse, error)
	CreateSession(session models.Session) (models.Session, error)
	GetSession(id uuid.UUID) (models.Session, error)
	MarkLogAsArchived(logID int) error
	DeleteLog(logID int) error
	CreateLog(log models.Log) (models.Log, error)
	StoreTaskID(et models.EventTask) error
	IsTaskCanceled(taskID string) (bool, error)
	UpdateEvent(event models.Event) (models.Event, error)
	CancelTasksForEvent(eventID int) error
	DeleteTaskEvent(taskID string) error
}

// Repository struct holds the GORM database instance
type Repo struct {
	DB *gorm.DB
}

// NewRepository initializes and returns a Repository instance
func NewRepository(db *gorm.DB) Repository {
	return &Repo{DB: db}
}
