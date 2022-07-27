package controller

import (
	middleware "ToDoApp/Middleware"
	service "ToDoApp/Service"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()

	//User login system
	router.HandleFunc("/user", service.SaveUser).Methods("POST")
	router.Handle("/login", middleware.CheckSessionDetails(http.HandlerFunc(service.LoginUser))).Methods("PATCH")
	// router.HandleFunc("/login", service.LoginUser).Methods("PATCH")

	// router.HandleFunc("/project", service.SaveProject).Methods("POST")
	// router.HandleFunc("/task", service.SaveTask).Methods("POST")

	// router.HandleFunc("/updateTaskStatus", service.UpdateTaskStatus).Methods("PATCH")
	// router.HandleFunc("/addTaskComment", service.AddClosureComment).Methods("PATCH")

	// router.HandleFunc("/getAllProjects", service.GetAllProjects).Methods("GET")
	// router.HandleFunc("/getAllTask", service.GetAllTask).Methods("GET")
	// router.HandleFunc("/getAllUser", service.GetAllUser).Methods("GET")
	// router.HandleFunc("/getDevTask", service.GetDevTask).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../Resources/template/")))

	router.Use(mux.CORSMethodMiddleware(router))
	return router
}
