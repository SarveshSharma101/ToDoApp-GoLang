package utility

import (
	datamodel "ToDoApp/DataModel"
	"math/rand"
	"time"
)

func GetRandomAlphaNumbericString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const r = 62
	alphaNum := make([]byte, length)
	for i := 0; i < length; i++ {
		alphaNum[i] = charset[rand.Intn(r)]
	}
	return string(alphaNum)
}

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
