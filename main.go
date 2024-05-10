package main

import (
	"net/http"

	"Shawty/controllers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	// List routes
	router.HandleFunc("/list", controllers.ListLists).Methods("GET")
	router.HandleFunc("/list/create", controllers.CreateList).Methods("POST")
	router.HandleFunc("/list/delete", controllers.DeleteList).Methods("POST")
	router.HandleFunc("/list/update", controllers.UpdateListTitle).Methods("POST")

	// Task routes
	router.HandleFunc("/task", controllers.GetItems).Methods("GET")
	router.HandleFunc("/task/create", controllers.CreateTask).Methods("POST")
	//router.HandleFunc("/task/delete", controllers.DeleteTask).Methods("POST")
	//router.HandleFunc("/task/complete", controllers.CompleteTask).Methods("POST")

	http.ListenAndServe(":8000", router)
}
