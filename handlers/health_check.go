package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Check the health of the server
// @Description Returns the server's health status
// @Tags Health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (server *Server) HealthCheckHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Server is Up and Running",
	})
}
