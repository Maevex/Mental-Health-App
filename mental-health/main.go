package main

import (
	"fmt"
	"log"
	"mental-health/config"
	"mental-health/routes"
	"net/http"

	"github.com/gorilla/handlers" // Import CORS handler
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	// Setup routes
	r := routes.SetupRoutes()

	// Menambahkan middleware CORS
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{
			"http://localhost:8081",
			"http://157.66.34.218", // IP VPS kamu
		}), // Origin frontend kamu
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}), // Metode yang diizinkan
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),           // Headers yang diizinkan
	)

	// Menjalankan server dengan middleware CORS
	fmt.Println("Server running on port http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", cors(r))) // Memasukkan middleware CORS
}
