package mysql_query

import (
	"go-sso/conf"
	"go-sso/db/model"
)

// 构建 mock database
func SetupTests() {
	model.Migrate()
}

// mysql 测试数据库
func SetupTestMysql() {
	c := conf.GetConfig()
	model.InitMysql(c) // 指向测试数据库
	teardownTests()    // 清空数据库
	model.Migrate()
}

func teardownTests() {
	model.DB.DropTableIfExists(&model.User{})
	model.DB.DropTableIfExists(&model.UserProfile{})
}
