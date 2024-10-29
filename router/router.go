package router

import (
	"go-postgres/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", middleware.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/resetpassword", middleware.Login).Methods("POST", "OPTIONS")
	return router
}
