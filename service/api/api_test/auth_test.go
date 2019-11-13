package api_test

import (
	"bytes"
	"fmt"
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

	username1 = "zs"
	username2 = "zs123"
	telephone1 = "12345678901"
	telephone2 = "18817550909"
	email1 = "xxx@189.com"
	email2 = "yyy@189.com"

	upTestCases = []model.UserParams{
		{Account:username, Password:password},
		{Account:telephone, Password:password},
		{Account:email, Password:password},
	}
	registerTestCases = []struct{
		m model.RegisterParams
		expected int
	}{
		{model.RegisterParams{Username: username, Password: password, Telephone:telephone}, 400},
		{model.RegisterParams{Username: username, Password: password, Telephone:telephone2}, 400},
		{model.RegisterParams{Username: username1, Password: password, Telephone:telephone2}, 400},
		{model.RegisterParams{Username: username2, Password: password, Telephone:telephone2}, 200},
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
	for _, testCase := range upTestCases {
		// 获取一个请求实例
		w := httptest.NewRecorder()
		u, _ := json.Marshal(testCase)
		// 构造请求
		req, _ := http.NewRequest("POST", "/api/public/v1/auth/login/", bytes.NewReader(u))
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}
}

func TestViewRegister(t *testing.T) {
	for _, testCase := range registerTestCases {
		w := httptest.NewRecorder()
		u, _ := json.Marshal(testCase.m)
		req, _ := http.NewRequest("POST", "/api/public/v1/auth/register/", bytes.NewReader(u))
		router.ServeHTTP(w, req)
		fmt.Println(testCase, w.Body)
		assert.Equal(t, testCase.expected, w.Code)
	}

}
