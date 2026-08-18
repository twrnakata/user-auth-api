package main

import (
	"context"
	"log"

	httproute "backend-challenge-golang/internal/http/route"
	repositoryauth "backend-challenge-golang/internal/repository/auth"
	repositoryuser "backend-challenge-golang/internal/repository/user"
	serviceauth "backend-challenge-golang/internal/service/auth"
	serviceuser "backend-challenge-golang/internal/service/user"
	"backend-challenge-golang/pkg/configuration"
	"backend-challenge-golang/pkg/datetime"
	jwtpkg "backend-challenge-golang/pkg/jwt"
	"backend-challenge-golang/pkg/mongodb"
)

func main() {
	if err := datetime.SetDefaultTimeZone(datetime.TimeZoneAsiaBangkok); err != nil {
		log.Fatal(err)
	}

	if err := configuration.InitConfig(); err != nil {
		log.Fatal(err)
	}

	executionContext := context.Background()
	mongoClient, err := mongodb.Connect(executionContext, configuration.Env.MONGO_URI)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if disconnectErr := mongoClient.Disconnect(executionContext); disconnectErr != nil {
			log.Printf("mongo disconnect error: %v", disconnectErr)
		}
	}()

	userCollection := mongodb.UserCollection(mongoClient, configuration.Env.MONGO_DATABASE)

	registerRepository, err := repositoryauth.NewAuthRegisterRepository(executionContext, userCollection)
	if err != nil {
		log.Fatal(err)
	}

	loginRepository, err := repositoryauth.NewAuthLoginRepository(userCollection)
	if err != nil {
		log.Fatal(err)
	}

	getUserRepository, err := repositoryuser.NewGetUserRepository(userCollection)
	if err != nil {
		log.Fatal(err)
	}

	jwtService, err := jwtpkg.NewJWTService(configuration.Env.JWT_SECRET, jwtpkg.DefaultExpireDuration)
	if err != nil {
		log.Fatal(err)
	}

	registerService := &serviceauth.RegisterUserService{
		Repository: registerRepository,
	}
	loginService := &serviceauth.AuthLoginService{
		Repository:   loginRepository,
		TokenService: jwtService,
	}

	getUserService := &serviceuser.GetUserService{
		Repository: getUserRepository,
	}

	application := httproute.NewApp(registerService, loginService, getUserService, jwtService)

	log.Printf("listening on :%s", configuration.Env.PORT)
	if err := application.Listen(":" + configuration.Env.PORT); err != nil {
		log.Fatal(err)
	}
}
