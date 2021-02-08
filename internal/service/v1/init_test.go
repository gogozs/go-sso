package v1_test

import (
	"go-sso/conf"
	"go-sso/internal/repository/storage/mysql"
	"go-sso/internal/routers"
	"go-sso/pkg/permission"
	"go-sso/registry"
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
	mysql.SetupTests(c) // 初始化mock database
	registry.InitStorage(mysql.DB)
	InitRouter(c)
}

func TearDown() {
}

func InitRouter(c *conf.Config) {
	permission.InitPermission(c)
	router = routers.InitRouter(c, registry.GetStorage())
	user = mysql.User{Username: username, Password: password, Role: "superuser", Telephone: telephone, Email: email}
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
