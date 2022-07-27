package middleware

import (
	service "ToDoApp/Service"
	"net/http"
)

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
