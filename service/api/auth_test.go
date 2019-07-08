package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestViewLogin(t *testing.T) {
	router := InitRouter()
	// 获取一个请求实例
	w := httptest.NewRecorder()
	// 构造请求
	// 参数依次是 请求方法、路由、参数
	req, _ := http.NewRequest("POST", "/api/v1/auth/login/", nil)
	// 执行
	router.ServeHTTP(w, req)
	fmt.Println(w.Body)
	assert.Equal(t, 400, w.Code)
}


func TestViewRegister(t *testing.T) {

}