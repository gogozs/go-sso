package model

import (
	"go-sso/conf"
)

// 构建 mock database
func SetupTests(c *conf.Config) {
	SetupTestMysql(c)
}

// mysql 测试数据库
func SetupTestMysql(c *conf.Config) {
	c.MySQL = c.TestMysql
	InitMysql(c)    // 指向测试数据库
	teardownTests() // 清空数据库
	Migrate()
}

func teardownTests() {
	DB.DropTableIfExists(TableArr...)
}
