package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or enviornment variables
type Config struct {
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	JWTSecretKey         string        `mapstructure:"JWT_SECRET_KEY"`
	TokenDuration        time.Duration `mapstructure:"TOKEN_DURATION"`
	DBHost               string        `mapstructure:"DB_HOST"`
	DBUser               string        `mapstructure:"DB_USER"`
	DBPassword           string        `mapstructure:"DB_PASSWORD"`
	DBName               string        `mapstructure:"DB_NAME"`
	DBPort               string        `mapstructure:"DB_PORT"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	LogArchiveDuration   time.Duration `mapstructure:"LOG_ARCHIVE_DURATION"`
	LogDeleteDuration    time.Duration `mapstructure:"LOG_DELETE_DURATION"`
	RedisUrl             string        `mapstructure:"REDIS_URL"`
}

// LoadConfig reads configuration from file or enviornment variable
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
