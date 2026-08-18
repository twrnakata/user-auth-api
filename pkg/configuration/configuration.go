package configuration

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/twrnakata/user-auth-api/pkg/apperror"
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
		return apperror.ErrMongoURIRequired
	}
	if Env.MONGO_DATABASE == "" {
		return apperror.ErrMongoDatabaseRequired
	}
	if Env.JWT_SECRET == "" {
		return apperror.ErrJWTSecretRequired
	}

	return nil
}
