package service

import (
	datamodel "ToDoApp/DataModel"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dbUrl = "?charset=utf8mb4&parseTime=True&loc=Local"

type (
	LoginReqBody struct {
		Uname    string `json:"uname"`
		Password string `json:"password"`
	}
	UpdateTaskStatusReqBody struct {
		Id     int    `json:"id"`
		Status string `json:"status"`
	}
	AddTaskCommentReqBody struct {
		Id     int
		Status string
	}
	GetDevTaskReqBody struct {
		Username string
	}
	SaveTaskReqBody struct {
		Tid   int    `json:"tid"`
		TName string `json:"tName"`
		TDesc string `json:"tDesc"`
	}
	SaveUserReqBody struct {
		UName    string `json:"uname"`
		Password string `json:"password"`
		Type     string `json:"type"`
	}
)

func InitDbConnection(dbName, uname, password, url string) {
	dbUrl = uname + ":" + password + "@" + url + "/" + dbName + dbUrl
	fmt.Println("------->", dbUrl)
	var err error
	DB, err = gorm.Open(mysql.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&datamodel.Project{})
	DB.AutoMigrate(&datamodel.Task{})
	DB.AutoMigrate(&datamodel.User{})
}

func SaveProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var project datamodel.Project

	json.NewDecoder(r.Body).Decode(&project)
	var checkProject datamodel.Project
	DB.Table("projects").Where("pid", project.Pid).Select("*").Scan(&checkProject)
	if len(checkProject.Name) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Project already exist please select some new username")
		return
	}

	if err := DB.Create(project).Error; err != nil {
		fmt.Println("***********Error********\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There was some error")
		return
	}
	json.NewEncoder(w).Encode(project)
}

func SaveTask(w http.ResponseWriter, r *http.Request) {
	var task SaveTaskReqBody
	var checkTask datamodel.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		panic(err)
	}
	w.Header().Set("content-type", "application/json")

	DB.Table("tasks").Where("tid", task.Tid).Select("*").Scan(&checkTask)
	if len(checkTask.TaskName) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Task already exist please select some new username")
		return
	}
	checkTask.Tid = task.Tid
	checkTask.TaskName = task.TName
	checkTask.TaskDesc = task.TDesc
	if err := DB.Select("Tid", "TaskName", "TaskDesc").Create(checkTask).Error; err != nil {
		fmt.Println("***********Error********\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There was some error")
		return
	}
	json.NewEncoder(w).Encode(task)
}

func SaveUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var userReq SaveUserReqBody
	var user datamodel.User

	json.NewDecoder(r.Body).Decode(&userReq)

	DB.Table("users").Where("username", userReq.UName).Select("*").Scan(&user)

	if len(user.Username) > 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode("User already exist please select some new username")
		return
	}
	user.Password = userReq.Password
	user.Username = userReq.UName
	user.Type = userReq.Type
	if err := DB.Select("Uid", "Username", "Password").Create(user).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There was some error")
		return
	}
	json.NewEncoder(w).Encode(user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var loginReqBody LoginReqBody
	var user datamodel.User

	json.NewDecoder(r.Body).Decode(&loginReqBody)
	DB.Table("users").Where("username", loginReqBody.Uname).Select("*").Scan(&user)
	if len(user.Username) <= 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("No Such user found")
	} else if user.Password != loginReqBody.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Password did not match")
	} else {
		user.LoginStatus = true
		DB.Save(&user)
		json.NewEncoder(w).Encode(user)
	}
}

func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var project []datamodel.Project
	DB.Find(&project)
	json.NewEncoder(w).Encode(project)
}

func GetAllTask(w http.ResponseWriter, r *http.Request) {
	var task []datamodel.Task
	w.Header().Set("content-type", "application/json")
	DB.Find(&task)
	json.NewEncoder(w).Encode(task)
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var user []datamodel.User
	DB.Find(&user)
	json.NewEncoder(w).Encode(user)
}

func UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var updateTaskStatus UpdateTaskStatusReqBody
	json.NewDecoder(r.Body).Decode(&updateTaskStatus)
	if err := DB.Model(&datamodel.Task{}).Where("tid", updateTaskStatus.Id).Update("status", updateTaskStatus.Status).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There was some error")
		return
	}
	json.NewEncoder(w).Encode("Status updated")
}

func AddClosureComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var addTaskCommentReqBody AddTaskCommentReqBody
	json.NewDecoder(r.Body).Decode(&addTaskCommentReqBody)
	if err := DB.Model(&datamodel.Task{}).Where("tid", addTaskCommentReqBody.Id).Update("closure_comment", addTaskCommentReqBody.Status); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There was some error")
		return
	}
	json.NewEncoder(w).Encode("Comment added")
}

func GetDevTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var getDevTaskReqBody GetDevTaskReqBody
	var user datamodel.User

	json.NewDecoder(r.Body).Decode(&getDevTaskReqBody)
	DB.Table("users").Where("username", getDevTaskReqBody.Username).Select("*").Scan(&user)
	DB.First(&user)
	if len(user.Username) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("No Such user found")
		return
	}

	var tasks []datamodel.Task
	DB.Where("user_id", user.Uid).Find(&tasks)
	json.NewEncoder(w).Encode(tasks)
}
