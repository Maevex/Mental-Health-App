package main

import (
	"fmt"
	"log"
	"mental-health/config"
	"mental-health/routes"
	"net/http"
)

func main() {
	// Koneksi database
	config.ConnectDB()

	// Setup routes
	r := routes.SetupRoutes()

	// Jalankan server di port 8080
	fmt.Println("Server running on port http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
