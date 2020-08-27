package mysql_query

import (
	"go-sso/conf"
	"go-sso/db/model"
)

// 构建 mock database
func SetupTests(c *conf.Config) {
	SetupTestMysql(c)
}

// mysql 测试数据库
func SetupTestMysql(c *conf.Config) {
	c.MySQL = c.TestMysql
	model.InitMysql(c) // 指向测试数据库
	teardownTests()    // 清空数据库
	model.Migrate()
}

func teardownTests() {
	model.DB.DropTableIfExists(model.TableArr...)
}
