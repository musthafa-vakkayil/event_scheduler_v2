package handlers

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	"github.com/musthafa-vakkayil/event_scheduler_v2/mocks"
	"github.com/musthafa-vakkayil/event_scheduler_v2/utils"
	"github.com/stretchr/testify/assert"
)

func SetupTestServer(t *testing.T, mockRepo *mocks.Repository, mockCache *mocks.Cache) *Server {
	config := config.Config{
		TokenDuration: time.Minute,
		JWTSecretKey:  utils.RandomString(32),
	}

	server, err := NewServer(config, mockRepo, mockCache)
	assert.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
