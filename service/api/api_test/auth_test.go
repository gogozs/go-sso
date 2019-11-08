package api_test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	model2 "go-qiuplus/db/model"
	"go-qiuplus/db/query"
	"go-qiuplus/pkg/json"
	"go-qiuplus/pkg/log"
	"go-qiuplus/service/api/routes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	user   model2.User
	router *gin.Engine
)

func init() {
	router = routes.GetRouter()
	query.SetupTests() // 初始化mock database
	user = model2.User{Username: "test", Password: "testpassword", Role: "superuser"}
	err := model2.CreateUser(user)
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
