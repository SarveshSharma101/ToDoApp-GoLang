package controller

import (
	middleware "ToDoApp/Middleware"
	service "ToDoApp/Service"
	"net/http"

	"github.com/gorilla/mux"
)

//Initialize the router
func GetRouter() *mux.Router {
	router := mux.NewRouter()

	//User login system
	router.HandleFunc("/user", service.SaveUser).Methods("POST")
	router.Handle("/login", middleware.CheckSessionDetails(http.HandlerFunc(service.LoginUser))).Methods("PATCH")

	//create project and task
	router.Handle("/project", middleware.CheckUserSession(http.HandlerFunc(service.SaveProject))).Methods("POST")
	router.Handle("/task", middleware.CheckUserSession(http.HandlerFunc(service.SaveTask))).Methods("POST")

	router.Handle("/getAllProjects", middleware.CheckUserSession(http.HandlerFunc(service.GetAllProjects))).Methods("GET")
	router.Handle("/getAllTask", middleware.CheckUserSession(http.HandlerFunc(service.GetAllTask))).Methods("GET")
	router.Handle("/getAllUser", middleware.CheckUserSession(http.HandlerFunc(service.GetAllUser))).Methods("GET")
	router.Handle("/getDevTask/{userId}", middleware.CheckUserSession(http.HandlerFunc(service.GetDevTask))).Methods("GET")
	router.Handle("/getAllDev", middleware.CheckUserSession(http.HandlerFunc(service.GetAllDev))).Methods("GET")
	router.Handle("/addTaskComment", middleware.CheckUserSession(http.HandlerFunc(service.AddClosureComment))).Methods("PATCH")
	router.Handle("/updateTaskStatus", middleware.CheckUserSession(http.HandlerFunc(service.UpdateTaskStatus))).Methods("PATCH")

	router.Handle("/logout", middleware.CheckUserSession(http.HandlerFunc(service.Logout))).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../Resources/template/")))

	router.Use(mux.CORSMethodMiddleware(router))
	return router
}
