package utility

import (
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
