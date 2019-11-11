package util

import (
	"math/rand"
	"time"
)

// 验证随机验证码
func RandomCode() string {
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i:=0;i<6;i++ {
		c := rand.Intn(10) + 48
		code[i] = uint8(c)
	}
	return string(code)
}
