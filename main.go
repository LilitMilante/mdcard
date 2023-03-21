package main

import (
	"log"

	"medical-card/internal/api"
	"medical-card/internal/app"
	"medical-card/internal/dal"
	"medical-card/internal/service"
)

func main() {
	c, err := app.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := app.NewPostgresClient(c.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	patientRepository := dal.NewPatientRepository(db)
	patientService := service.NewPatientService(patientRepository)
	patientHandler := api.NewPatientHandler(patientService)
	authMw := api.NewAuthMiddleware(patientService)
	server := api.NewServer(c.Port, patientHandler, authMw)

	log.Println("server started at:", c.Port)
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
