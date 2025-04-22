package main

import (
	"fmt"
	"log"
	"mental-health/config"
	"mental-health/routes"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	r := routes.SetupRoutes()

	fmt.Println("Server running on port http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
