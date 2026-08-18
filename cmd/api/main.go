package main

import (
	"log"
	"os"

	"backend-challenge-golang/internal/application/user"
	"backend-challenge-golang/internal/interfaces/http"
)

func main() {
	// Minimal bootstrap: Step 1 focuses on wiring + contracts.
	// Business implementations (MongoDB/JWT) will come in later steps.
	registerSvc := &user.NotImplementedRegisterUserService{}

	app := httpapi.NewApp(registerSvc)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}

	log.Printf("listening on :%s", addr)
	if err := app.Listen(":" + addr); err != nil {
		log.Fatal(err)
	}
}

