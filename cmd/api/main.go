package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	httproute "github.com/twrnakata/user-auth-api/internal/http/route"
	"github.com/twrnakata/user-auth-api/internal/job"
	repositoryauth "github.com/twrnakata/user-auth-api/internal/repository/auth"
	repositoryuser "github.com/twrnakata/user-auth-api/internal/repository/user"
	serviceauth "github.com/twrnakata/user-auth-api/internal/service/auth"
	serviceuser "github.com/twrnakata/user-auth-api/internal/service/user"
	"github.com/twrnakata/user-auth-api/pkg/configuration"
	"github.com/twrnakata/user-auth-api/pkg/datetime"
	jwtpkg "github.com/twrnakata/user-auth-api/pkg/jwt"
	"github.com/twrnakata/user-auth-api/pkg/mongodb"
)

func main() {
	if err := datetime.SetDefaultTimeZone(datetime.TimeZoneAsiaBangkok); err != nil {
		log.Fatal(err)
	}

	if err := configuration.InitConfig(); err != nil {
		log.Fatal(err)
	}

	executionContext := context.Background()
	database, err := mongodb.Connect(executionContext, configuration.Env.MONGO_URI)
	if err != nil {
		log.Fatal(err)
	}

	userCollection := mongodb.UserCollection(database, configuration.Env.MONGO_DATABASE)

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

	listUserRepository, err := repositoryuser.NewListUserRepository(userCollection)
	if err != nil {
		log.Fatal(err)
	}

	updateUserRepository, err := repositoryuser.NewUpdateUserRepository(userCollection)
	if err != nil {
		log.Fatal(err)
	}

	deleteUserRepository, err := repositoryuser.NewDeleteUserRepository(userCollection)
	if err != nil {
		log.Fatal(err)
	}

	countUserRepository, err := repositoryuser.NewCountUserRepository(userCollection)
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
	listUserService := &serviceuser.ListUserService{
		Repository: listUserRepository,
	}
	updateUserService := &serviceuser.UpdateUserService{
		Repository: updateUserRepository,
	}
	deleteUserService := &serviceuser.DeleteUserService{
		Repository: deleteUserRepository,
	}
	countUserService := &serviceuser.CountUserService{
		Repository: countUserRepository,
	}

	userCountJob, err := job.NewUserCountJob(countUserService, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	jobContext, cancelJob := context.WithCancel(context.Background())
	jobDone := &sync.WaitGroup{}
	jobDone.Add(1)
	go func() {
		defer jobDone.Done()
		userCountJob.Run(jobContext)
	}()

	application := httproute.NewApp(registerService, loginService, listUserService, getUserService, updateUserService, deleteUserService, jwtService)

	listenErr := make(chan error, 1)
	go func() {
		listenErr <- application.Listen(":" + configuration.Env.PORT)
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	log.Printf("listening on :%s", configuration.Env.PORT)

	select {
	case err := <-listenErr:
		cancelJob()
		jobDone.Wait()
		if disconnectErr := database.Disconnect(executionContext); disconnectErr != nil {
			log.Printf("mongo disconnect error: %v", disconnectErr)
		}
		if err != nil {
			log.Fatal(err)
		}
	case <-signalChannel:
		log.Printf("shutdown signal received")
		if err := gracefulShutdown(executionContext, application, cancelJob, jobDone, database, nil); err != nil {
			log.Printf("graceful shutdown error: %v", err)
		}
		if err := <-listenErr; err != nil {
			log.Printf("listen error: %v", err)
		}
	}
}
