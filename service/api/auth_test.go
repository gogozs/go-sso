package api

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-weixin/pkg/json"
	"go-weixin/pkg/log"
	"go-weixin/service/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	user models.User
)

func init() {
	InitRouter()
	models.SetupTests() // 初始化mock database
	user = models.User{Username: "test", Password: "testpassword", Role: "superuser"}
	err := models.CreateUser(user)
	log.Error(err)
}

func TestViewLogin(t *testing.T) {
	// 获取一个请求实例
	w := httptest.NewRecorder()
	// 构造请求
	// 参数依次是 请求方法、路由、参数
	u, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login/", bytes.NewReader(u))
	// 执行
	router.ServeHTTP(w, req)
	fmt.Println(w.Body)
	assert.Equal(t, 200, w.Code)
}

func TestViewRegister(t *testing.T) {

}
