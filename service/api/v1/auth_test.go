package v1

import (
	"go-weixin/service/api"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestViewLogin(t *testing.T) {
	router := api.InitRouter()
	// 获取一个请求实例
	w := httptest.NewRecorder()
	// 构造请求
	// 参数依次是 请求方法、路由、参数
	req, _ := http.NewRequest("POST", "/api/v1/auth/login/", nil)
	// 执行
	router.ServeHTTP(w, req)
	// 断言

	assert.Equal(t, 200, w.Code)
}


func TestViewRegister(t *testing.T) {

}