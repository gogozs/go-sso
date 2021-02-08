package util

import (
	"math/rand"
	"time"
)

// 验证随机验证码
func RandomCode(n int) string {
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, n)
	for i := 0; i < n; i++ {
		c := rand.Intn(10) + 48
		code[i] = uint8(c)
	}
	return string(code)
}
