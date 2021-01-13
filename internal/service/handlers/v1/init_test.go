package v1_test

import (
	"go-sso/conf"
	"go-sso/internal/registry"
	"go-sso/internal/repository/mysql/model"
	"go-sso/internal/repository/mysql/mysql_query"
	"go-sso/internal/service/routers"
	"go-sso/pkg/permission"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetUp()
	m.Run()
	TearDown()
	os.Exit(0)
}

func SetUp() {
	c := InitConfig()
	mysql_query.SetupTests(c) // 初始化mock database
	registry.SetStorage(model.DB)
	InitRouter(c)
}

func TearDown() {
}

func InitRouter(c *conf.Config) {
	permission.InitPermission(c)
	router = routers.InitRouter(c)
	user = model.User{Username: username, Password: password, Role: "superuser", Telephone: telephone, Email: email}
	_, err := registry.GetStorage().Create(&user)
	if err != nil {
		panic(err)
	}
}

func InitConfig() *conf.Config {
	c, err := conf.InitConfig()
	if err != nil {
		panic(err)
	}
	return c
}
