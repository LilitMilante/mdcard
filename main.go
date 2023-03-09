package main

import (
	"database/sql"
	"fmt"

	"medical-card/internal/api"
	"medical-card/internal/dal"
	"medical-card/internal/service"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 8080
	user     = "dev"
	password = "dev"
	dbname   = "mdcard"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	patientRepository := dal.NewPatientRepository(db)
	patientService := service.NewPatientService(patientRepository)
	patientHandler := api.NewPatientHandler(patientService)
	server := api.NewServer("8081", patientHandler)

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
