package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	"github.com/musthafa-vakkayil/event_scheduler_v2/middleware"
	"github.com/musthafa-vakkayil/event_scheduler_v2/repo"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Server serves HTTP requests for our banking service
type Server struct {
	Config     config.Config
	TokenMaker token.Maker
	Router     *gin.Engine
	Repo       repo.Repository
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config config.Config, repo repo.Repository) (*Server, error) {
	maker, err := token.NewJWTMaker(config.JWTSecretKey)
	if err != nil {
		return nil, fmt.Errorf("unable create token maker %w", err)
	}

	server := &Server{TokenMaker: maker, Config: config, Repo: repo}

	server.SetupRoutes()

	return server, nil
}

func (server *Server) SetupRoutes() {
	server.Router = gin.Default()
	router := server.Router

	router.POST("/login", server.Login)
	router.POST("/users", server.CreateUser)

	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(server.TokenMaker))

	authRoutes.GET("/users/:username", server.GetUser)
	authRoutes.DELETE("/users/:username", server.DeleteUser)
	authRoutes.POST("/events", server.CreateEvent)
	authRoutes.GET("/events", server.ListEvents)
	authRoutes.GET("/events/:id", server.GetEvent)
	authRoutes.GET("/events/:id/execute", server.ExecuteAPIEvent)
	authRoutes.DELETE("/events/:id", server.DeleteEvent)
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}

// ConnectGORM initializes a GORM DB connection
func ConnectGORM(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		cfg.DBHost,     // e.g., "localhost"
		cfg.DBUser,     // e.g., "postgres"
		cfg.DBPassword, // e.g., "yourpassword"
		cfg.DBName,     // e.g., "mydb"
		cfg.DBPort,     // e.g., "5432"
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GORM DB: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get DB from GORM: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("✅ GORM DB connected successfully")
	return db, nil
}
