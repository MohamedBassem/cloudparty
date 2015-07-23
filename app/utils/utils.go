package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func RandomString(length int) string {
	stringSeed := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	var ret string
	for i := 0; i < length; i++ {
		ret += string(stringSeed[rand.Intn(len(stringSeed))])
	}
	return ret
}
