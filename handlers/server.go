package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	db "github.com/musthafa-vakkayil/event_scheduler_v2/db/sqlc"
	"github.com/musthafa-vakkayil/event_scheduler_v2/middleware"
	"github.com/musthafa-vakkayil/event_scheduler_v2/token"
)

// Server serves HTTP requests for our banking service
type Server struct {
	Config     config.Config
	Store      db.Store
	TokenMaker token.Maker
	Router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config config.Config, store db.Store) (*Server, error) {
	maker, err := token.NewJWTMaker(config.JWTSecretKey)
	if err != nil {
		return nil, fmt.Errorf("unable create token maker %w", err)
	}
	server := &Server{Store: store, TokenMaker: maker, Config: config}

	server.SetupRoutes()

	return server, nil
}

func (server *Server) SetupRoutes() {
	router := gin.Default()

	router.POST("/login", server.Login)
	router.POST("/users", server.CreateUser)

	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(server.TokenMaker))

	authRoutes.GET("/users/:username", server.GetUser)
	authRoutes.POST("/events", server.CreateEvent)
	authRoutes.GET("/events", server.ListEvents)

	server.Router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
