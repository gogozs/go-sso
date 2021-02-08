package v1_test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-sso/internal/repository/storage/mysql"
	"go-sso/internal/service/viewset"
	"go-sso/pkg/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	user      mysql.User
	router    *gin.Engine
	username  = "test"
	password  = "testpassword"
	telephone = "12345678901"
	email     = "test@test.com"

	username1  = "zs"
	username2  = "zs123"
	telephone1 = "12345678901"
	telephone2 = "18817550909"

	upTestCases = []mysql.UserParams{
		{Account: username, Password: password},
		{Account: telephone, Password: password},
		{Account: email, Password: password},
	}
	registerTestCases = []struct {
		m        mysql.RegisterParams
		expected int
	}{
		{mysql.RegisterParams{Username: username, Password: password, Telephone: telephone}, 400},
		{mysql.RegisterParams{Username: username, Password: password, Telephone: telephone1}, 400},
		{mysql.RegisterParams{Username: username, Password: password, Telephone: telephone2}, 400},
		{mysql.RegisterParams{Username: username1, Password: password, Telephone: telephone2}, 400},
		{mysql.RegisterParams{Username: username2, Password: password, Telephone: telephone2}, 200},
	}
)

func TestViewLogin(t *testing.T) {
	for _, testCase := range upTestCases {
		// 获取一个请求实例
		w := httptest.NewRecorder()
		u, _ := json.Marshal(testCase)
		// 构造请求
		req, _ := http.NewRequest("POST", "/api/public/v1/auth/login/", bytes.NewReader(u))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}
}

func TestViewRegister(t *testing.T) {
	for _, testCase := range registerTestCases {
		w := httptest.NewRecorder()
		u, _ := json.Marshal(testCase.m)
		req, _ := http.NewRequest("POST", "/api/public/v1/auth/register/", bytes.NewReader(u))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, testCase.expected, w.Code)
	}

}

func TestChangePassword(t *testing.T) {
	// get token
	w1 := httptest.NewRecorder()
	u1, _ := json.Marshal(mysql.UserParams{Account: username, Password: password})
	// 构造请求
	req1, _ := http.NewRequest("POST", "/api/public/v1/auth/login/", bytes.NewReader(u1))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)
	m := viewset.Response{}
	_ = json.Unmarshal(w1.Body.Bytes(), &m)
	token := m.Data.(map[string]interface{})["token"].(string)

	// change password
	w2 := httptest.NewRecorder()
	newPassword := password + "111"
	cp := mysql.ChangePasswordParams{RawPassword: password, NewPassword: newPassword}
	u2, _ := json.Marshal(&cp)
	req2, _ := http.NewRequest("POST", "/api/v1/auth/change-password/", bytes.NewReader(u2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", fmt.Sprintf("Token %sms", token))
	router.ServeHTTP(w2, req2)
	require.Equal(t, 200, w2.Code)
	fmt.Println(w2.Body)

	// newPassword login
	w3 := httptest.NewRecorder()
	u3, _ := json.Marshal(mysql.UserParams{Account: username, Password: newPassword})
	req3, _ := http.NewRequest("POST", "/api/public/v1/auth/login/", bytes.NewReader(u3))
	req3.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w3, req3)
	assert.Equal(t, 200, w3.Code)
	fmt.Println(w3.Body)
}
