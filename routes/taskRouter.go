package routes

import (
	"github.com/DebanjanBarman/todo/controller"
	"github.com/gorilla/mux"
	"net/http"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/tasks", controller.CreateNewTask).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/tasks", controller.GetTasks).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/tasks/{id}", controller.GetTask).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/tasks/{id}", controller.UpdateTask).Methods(http.MethodPatch)
	router.HandleFunc("/api/v1/tasks/{id}", controller.DeleteTask).Methods(http.MethodDelete)

	return router
}
