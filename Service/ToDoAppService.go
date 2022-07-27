package service

import (
	datamodel "ToDoApp/DataModel"
	utility "ToDoApp/Utility"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dbUrl = "?charset=utf8mb4&parseTime=True&loc=Local"

//Initialize DB and redis connection
func InitDbConnection(dbName, uname, password, url string) {
	dbUrl = uname + ":" + password + "@" + url + "/" + dbName + dbUrl
	var err error
	DB, err = gorm.Open(mysql.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&datamodel.Project{})
	DB.AutoMigrate(&datamodel.Task{})
	DB.AutoMigrate(&datamodel.User{})
	InitRedisClient()
}

//Register a user
func SaveUser(w http.ResponseWriter, r *http.Request) {
	//set content type to json
	w.Header().Set("content-type", "application/json")

	var userReq datamodel.SaveUserReqBody
	var user datamodel.User

	//parse the req body to json datamodel
	json.NewDecoder(r.Body).Decode(&userReq)
	//validate request
	isValid, msg := ValidateUser(userReq)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}

	//check if such a user already exist
	DB.Table("users").Where("username", userReq.UName).Select("*").Scan(&user)
	if len(user.Username) > 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode("User already exist please select some new username")
		return
	}

	//store the user in DB
	user.Password = userReq.Password
	user.Username = userReq.UName
	user.Type = userReq.Type
	//throw error if error occurs while saving the user
	if err := DB.Select("Uid", "Username", "Password", "Type").Create(&user).Error; err != nil {
		fmt.Println("Error! >>>>>>>>", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There was some error")
		return
	}
	//return success response if no error
	json.NewEncoder(w).Encode(user)
}

//Login the user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	//set content type to json
	w.Header().Set("content-type", "application/json")

	var loginReqBody datamodel.LoginReqBody
	var user datamodel.User

	//parse the req body to datamodel
	json.NewDecoder(r.Body).Decode(&loginReqBody)

	//check if user exist and validate the password
	DB.Table("users").Where("username", loginReqBody.Uname).Select("*").Scan(&user)
	if len(user.Username) <= 0 {

		//failure response if user not found
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("No Such user found")
	} else if user.Password != loginReqBody.Password {

		//failure response if password don't match
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Password did not match")
	} else {

		//if password match create a sessionId and store in redis
		sessionId := utility.GetRandomAlphaNumbericString(10)
		for CheckInRedis(sessionId) {
			sessionId = utility.GetRandomAlphaNumbericString(10)
		}
		storeSession(sessionId, datamodel.RedisUser{
			Username: user.Username,
			Role:     user.Type,
		})

		//send back response with sessionId in cookie details of response
		cookie := &http.Cookie{
			Name:   "sessionId",
			Value:  sessionId,
			MaxAge: 24 * 60 * 60,
		}
		w.Header().Add("session", cookie.String())
		json.NewEncoder(w).Encode(user)
	}
}

//Create a new project
func SaveProject(w http.ResponseWriter, r *http.Request) {
	//set content type to json
	w.Header().Set("content-type", "application/json")

	//parse req body to datamodel
	var project datamodel.Project

	json.NewDecoder(r.Body).Decode(&project)
	var checkProject datamodel.Project

	//check if the project already exist
	DB.Table("projects").Where("pid", project.Pid).Select("*").Scan(&checkProject)
	if len(checkProject.Name) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Project already exist please select some new username")
		return
	}

	//save the project
	if err := DB.Create(&project).Error; err != nil {
		fmt.Println("***********Error********\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There was some error")
		return
	}
	//send success response
	json.NewEncoder(w).Encode(project)
}

//Create a new Task
func SaveTask(w http.ResponseWriter, r *http.Request) {
	//set content type as json
	w.Header().Set("content-type", "application/json")

	//parse req to datamodel
	var task datamodel.SaveTaskReqBody
	var checkTask datamodel.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		panic(err)
	}

	//check if task already exist
	DB.Table("tasks").Where("tid", task.Tid).Select("*").Scan(&checkTask)
	if len(checkTask.TaskName) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Task already exist please select some new username")
		return
	}
	//save the task
	checkTask.Tid = task.Tid
	checkTask.TaskName = task.TName
	checkTask.TaskDesc = task.TDesc
	if err := DB.Select("Tid", "TaskName", "TaskDesc").Create(&checkTask).Error; err != nil {
		fmt.Println("***********Error********\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There was some error")
		return
	}
	json.NewEncoder(w).Encode(task)
}

//get all project
func GetAllProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var project []datamodel.Project
	DB.Find(&project)
	json.NewEncoder(w).Encode(project)
}

//get all task
func GetAllTask(w http.ResponseWriter, r *http.Request) {
	var task []datamodel.Task
	w.Header().Set("content-type", "application/json")
	DB.Find(&task)
	json.NewEncoder(w).Encode(task)
}

//get all users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var user []datamodel.User
	DB.Find(&user)
	json.NewEncoder(w).Encode(user)
}

//get all users
func GetAllDev(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var user []datamodel.User
	DB.Where("type = ?", "2").Find(&user)
	json.NewEncoder(w).Encode(user)
}

//get all task for dev
func GetDevTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	userId := mux.Vars(r)["userId"]

	var tasks []datamodel.Task
	DB.Where("user_id", userId).Find(&tasks)
	json.NewEncoder(w).Encode(tasks)
}

//update task status
func UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var updateTaskStatus datamodel.UpdateTaskStatusReqBody
	json.NewDecoder(r.Body).Decode(&updateTaskStatus)
	if err := DB.Model(&datamodel.Task{}).Where("tid", updateTaskStatus.Id).Update("status", updateTaskStatus.Status).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There was some error")
		return
	}
	json.NewEncoder(w).Encode("Status updated")
}

//Add closure comments
func AddClosureComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var addTaskCommentReqBody datamodel.AddTaskCommentReqBody
	json.NewDecoder(r.Body).Decode(&addTaskCommentReqBody)
	if err := DB.Model(&datamodel.Task{}).Where("tid", addTaskCommentReqBody.Id).Update("closure_comment", addTaskCommentReqBody.Status); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There was some error")
		return
	}
	json.NewEncoder(w).Encode("Comment added")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	c, err := r.Cookie("sessionId")
	if err != nil {
		panic(err)
	}
	DeleteSession(c.Value)
}

// =====================================================================================================================

func ValidateUser(userReq datamodel.SaveUserReqBody) (bool, string) {
	if len(userReq.UName) == 0 || len(userReq.Password) == 0 {
		return false, "username/password cannot be empty"
	} else if len(userReq.Password) < 8 {
		return false, "password length should be >=8"
	} else if !checkType(userReq.Type) {
		return false, "Type value must be {0,1,2}"
	}
	return true, "request is okay"
}

func checkType(k int) bool {
	_type := []int{0, 1, 2}
	for _, v := range _type {
		if v == k {
			return true
		}
	}
	return false
}
