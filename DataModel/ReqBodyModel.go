package datamodel

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
		Tid    int    `json:"tid"`
		TName  string `json:"tName"`
		TDesc  string `json:"tDesc"`
		UserId int    `json:"userId"`
	}
	SaveUserReqBody struct {
		UName    string `json:"uname"`
		Password string `json:"password"`
		Type     int    `json:"type"`
	}
)
