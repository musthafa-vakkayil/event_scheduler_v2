package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/cache"
	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	"github.com/musthafa-vakkayil/event_scheduler_v2/middleware"
	"github.com/musthafa-vakkayil/event_scheduler_v2/repo"
	"github.com/musthafa-vakkayil/event_scheduler_v2/tasks"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// Swagger
	_ "github.com/musthafa-vakkayil/event_scheduler_v2/docs"
)

// Server serves HTTP requests for our banking service
type Server struct {
	Config      config.Config
	TokenMaker  token.Maker
	Router      *gin.Engine
	Repo        repo.Repository
	Cache       cache.Cache
	TaskManager tasks.TaskManager
}

func NewServer(config config.Config, repo repo.Repository, cache cache.Cache, tm tasks.TaskManager) (*Server, error) {
	maker, err := token.NewJWTMaker(config.JWTSecretKey)
	if err != nil {
		return nil, fmt.Errorf("unable create token maker %w", err)
	}

	server := &Server{TokenMaker: maker, Config: config, Repo: repo, Cache: cache, TaskManager: tm}

	server.SetupRoutes()

	return server, nil
}

func (server *Server) SetupRoutes() {
	server.Router = gin.Default()

	// Swagger routes
	server.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router := server.Router

	// Enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins (or specify allowed domains)
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Open routes
	router.GET("/health", server.HealthCheckHandler)
	router.POST("/login", server.Login)
	router.POST("/users", server.CreateUser)
	router.POST("/token/renew", server.RenewAccessToken)

	// Authenticated routes with JWT
	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(server.TokenMaker))

	// User Routes
	authRoutes.GET("/users/:username", server.GetUser)
	authRoutes.DELETE("/users/:username", server.DeleteUser)

	// Event Routes
	authRoutes.POST("/events/api", server.CreateAPIEvent)
	authRoutes.POST("/events/schedule", server.CreateScheduledEvent)
	authRoutes.GET("/events", server.ListEvents)
	authRoutes.GET("/events/:id", server.GetEvent)
	authRoutes.GET("/events/:id/execute", server.ExecuteAPIEvent)
	authRoutes.DELETE("/events/:id", server.DeleteEvent)
	authRoutes.PUT("/events/schedule/:id", server.UpdateScheduledEvent)
	authRoutes.PUT("/events/api/:id", server.EditAPIEvent)

	// Log Routes
	authRoutes.GET("/logs", server.ListLogs)

	// Temporary Event routes
	authRoutes.POST("/test/events/api", server.CreateTestAPIEvent)
	authRoutes.POST("/test/events/schedule", server.CreateTestScheduledEvent)
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
