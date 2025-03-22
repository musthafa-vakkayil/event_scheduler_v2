package repo

import (
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user models.User) (models.User, error)
	GetUser(username string) (models.User, error)
	CreateEvent(event models.Event) (models.Event, error)
	GetEvent(id int) (models.Event, error)
	ListEvents(limit, offset int) ([]models.Event, error)
	ExecuteEvent(eventID int64, status string) error
	DeleteEvent(eventID int64) error
	DeleteUser(username string) error
}

// Repository struct holds the GORM database instance
type Repo struct {
	DB *gorm.DB
}

// NewRepository initializes and returns a Repository instance
func NewRepository(db *gorm.DB) Repository {
	return &Repo{DB: db}
}
