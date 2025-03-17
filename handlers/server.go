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
	config     config.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config config.Config, store db.Store) (*Server, error) {
	maker, err := token.NewJWTMaker(config.JWTSecretKey)
	if err != nil {
		return nil, fmt.Errorf("unable create token maker %w", err)
	}
	server := &Server{store: store, tokenMaker: maker, config: config}

	server.SetupRoutes()

	return server, nil
}

func (server *Server) SetupRoutes() {
	router := gin.Default()

	router.POST("/login", server.Login)
	router.POST("/users", server.createUser)

	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(server.tokenMaker))

	authRoutes.GET("/users", server.getUser)

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
