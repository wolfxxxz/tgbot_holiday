package config

import (
	"fmt"
	"holiday_bot/internal/apperrors"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Token             string `env:"TOKEN"`
	LogLevel          string `env:"LOGGER_LEVEL"`
	MongoHost         string `env:"MONGO_URL"`
	MongoPort         string `env:"MONGO_PORT"`
	UserName          string `env:"USER_NAME"`
	DBName            string `env:"DB_NAME"`
	Password          string `env:"PASSWORD"`
	TimeoutMongoQuery string `env:"TIMEOUT_MONGO_QUERY"`
}

func NewConfig() *Config {
	return &Config{}
}

func (v *Config) ParseConfig(path string, log *logrus.Logger) error {
	err := godotenv.Load(path)
	if err != nil {
		errMsg := fmt.Sprintf(" %s", err.Error())
		appError := apperrors.EnvConfigParseError.AppendMessage(errMsg)
		log.Error(appError)
		return appError
	}

	if err := env.Parse(v); err != nil {
		errMsg := fmt.Sprintf("%+v\n", err)
		appErr := apperrors.EnvConfigParseError.AppendMessage(errMsg)
		log.Error(appErr)
		return appErr
	}

	log.Info("Config has been parsed, succesfully!!!")
	return nil
}
