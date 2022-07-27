package middleware

import (
	service "ToDoApp/Service"
	"encoding/json"
	"net/http"
)

//Middleware to check if the session already exist for the user
func CheckSessionDetails(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("sessionId")
		if err != nil {
			panic(err)
		}
		if service.CheckInRedis(c.Value) {
			w.Write([]byte("Welcome again"))
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func CheckUserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("sessionId")
		if err != nil {
			panic(err)
		}
		if service.CheckIfUserHasValidSession(c.Value) {
			next.ServeHTTP(w, r)
		} else {
			json.NewEncoder(w).Encode("Please login and generate a valid sessionId")
		}
	})
}
