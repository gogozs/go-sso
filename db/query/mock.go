package query

import (
	"go-sso/conf"
	"go-sso/db/model"
)

// 构建 mock database
func SetupTests()  {
	c := conf.GetConfig()
	mysql := c.TestMysql
	model.Init(conf.MySQLConfig(mysql)) // 指向测试数据库
	teardownTests()                     // 清空数据库
	model.Migrate()
}

func teardownTests() {
	model.DB.DropTableIfExists(&model.User{})
	model.DB.DropTableIfExists(&model.UserProfile{})
}
