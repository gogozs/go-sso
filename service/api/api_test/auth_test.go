package api_test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-sso/db/model"
	"go-sso/db/query"
	"go-sso/pkg/json"
	"go-sso/pkg/log"
	"go-sso/service/api/routes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	user   model.User
	router *gin.Engine
	username = "test"
	password = "testpassword"
	telephone = "12345678901"
	email = "test@test.com"

	upTestCases = []model.UserParams{
		{Account:username, Password:password},
		{Account:telephone, Password:password},
		{Account:email, Password:password},
	}
)

func init() {
	router = routes.GetRouter()
	query.SetupTests() // 初始化mock database
	user = model.User{Username: username, Password: password, Role: "superuser", Telephone: telephone, Email: email}
	uq := &query.UserQuery{}
	_, err := uq.Create(&user)
	log.Error(err)
}

func TestViewLogin(t *testing.T) {
	// 获取一个请求实例
	w := httptest.NewRecorder()
	// 构造请求
	// 参数依次是 请求方法、路由、参数
	for _, testCase := range upTestCases {
		u, _ := json.Marshal(testCase)
		req, _ := http.NewRequest("POST", "/api/public/v1/auth/login/", bytes.NewReader(u))
		// 执行
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}
}

func TestViewRegister(t *testing.T) {

}
