package routes

import (
	"mental-health/controllers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/users", controllers.GetUsers).Methods("GET")

	return r
}
