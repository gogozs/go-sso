package v1_test

import (
	"go-sso/conf"
	"go-sso/registry"
	"go-sso/service/api/routes"
	"go-sso/storage/mysql/model"
	"go-sso/storage/mysql/mysql_query"
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
	router = routes.InitRouter(c)
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
