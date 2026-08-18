package configuration

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type Environment struct {
	PORT           string
	MONGO_URI      string
	MONGO_DATABASE string
	MONGO_USERNAME string
	MONGO_PASSWORD string
	JWT_SECRET     string
}

var Env Environment

func InitConfig() error {
	viper.Reset()
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	Env = Environment{
		PORT:           strings.TrimSpace(viper.GetString("PORT")),
		MONGO_URI:      strings.TrimSpace(viper.GetString("MONGO_URI")),
		MONGO_DATABASE: strings.TrimSpace(viper.GetString("MONGO_DATABASE")),
		MONGO_USERNAME: strings.TrimSpace(viper.GetString("MONGO_USERNAME")),
		MONGO_PASSWORD: strings.TrimSpace(viper.GetString("MONGO_PASSWORD")),
		JWT_SECRET:     strings.TrimSpace(viper.GetString("JWT_SECRET")),
	}

	if Env.PORT == "" {
		Env.PORT = "8080"
	}
	if Env.MONGO_URI == "" {
		return errors.New("MONGO_URI is required")
	}
	if Env.MONGO_DATABASE == "" {
		return errors.New("MONGO_DATABASE is required")
	}
	if Env.JWT_SECRET == "" {
		return errors.New("JWT_SECRET is required")
	}

	return nil
}
