package routes

import (
	"mental-health/config"
	"mental-health/controllers"
	middlewares "mental-health/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// User routes
	r.Handle("/users", middlewares.AdminMiddleware(http.HandlerFunc(controllers.GetUsers))).Methods("GET")
	r.Handle("/myuser", middlewares.JWTMiddleware(http.HandlerFunc(controllers.GetUserByID))).Methods("GET")
	r.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	r.Handle("/userupdate", middlewares.JWTMiddleware(http.HandlerFunc(controllers.UpdateUser))).Methods("PUT")
	r.Handle("/userdel", middlewares.JWTMiddleware(http.HandlerFunc(controllers.DeleteUser))).Methods("DELETE")
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	// Chat routes
	r.Handle("/chat", middlewares.JWTMiddleware(http.HandlerFunc(controllers.ChatHandler(config.DB)))).Methods("POST")
	r.Handle("/sesi", middlewares.JWTMiddleware(http.HandlerFunc(controllers.GetSessions))).Methods("GET")
	r.Handle("/sesi/{id}", middlewares.JWTMiddleware(http.HandlerFunc(controllers.GetDetailSesiHandler))).Methods("GET")
	//consultant routes
	r.Handle("/consultantIns", middlewares.AdminMiddleware(http.HandlerFunc(controllers.CreateConsultant))).Methods("POST")
	r.Handle("/consultants", middlewares.AdminMiddleware(http.HandlerFunc(controllers.GetAllConsultants))).Methods("GET")
	r.Handle("/consultantDel/{id}", middlewares.AdminMiddleware(http.HandlerFunc(controllers.DeleteConsultant))).Methods("DELETE")
	r.Handle("/consultantUpdate/{id}", middlewares.AdminMiddleware(http.HandlerFunc(controllers.UpdateConsultant))).Methods("PUT")
	r.Handle("/consultant/{id}", middlewares.AdminMiddleware(http.HandlerFunc(controllers.GetConsultantByID))).Methods("GET")

	//checking public
	r.HandleFunc("/ping", controllers.PublicCheckHandler).Methods("GET")
	return r
}
