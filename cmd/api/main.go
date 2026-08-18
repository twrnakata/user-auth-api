package main

import (
	"log"
	"os"

	domainauth "backend-challenge-golang/internal/domain/auth"
	httproute "backend-challenge-golang/internal/http/route"
	"backend-challenge-golang/pkg/datetime"
)

func main() {
	if err := datetime.SetDefaultTimeZone(datetime.TimeZoneAsiaBangkok); err != nil {
		log.Fatal(err)
	}

	// Minimal bootstrap: Step 1 focuses on wiring + contracts.
	// Business implementations (MongoDB/JWT) will come in later steps.
	registerService := &domainauth.NotImplementedRegisterUserService{}
	loginService := &domainauth.NotImplementedLoginUserService{}

	application := httproute.NewApp(registerService, loginService)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}

	log.Printf("listening on :%s", addr)
	if err := application.Listen(":" + addr); err != nil {
		log.Fatal(err)
	}
}
