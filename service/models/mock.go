package models

import (
	"go-weixin/config"
)

// 构建 mock database
func SetupTests()  {
	c := config.GetConfig()
	mysql := c.TestMysql
	Init(config.MySQLConfig(mysql))  // 指向测试数据库
	teardownTests() // 清空数据库
	Migrate()
}

func teardownTests() {
	DB.DropTableIfExists(&User{})
	DB.DropTableIfExists(&UserProfile{})
}
