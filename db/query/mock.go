package query

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
	mysql := c.TestMysql
	model.InitMysql(true, model.MySQLConfig{
		Host:     mysql.Host,
		Username: mysql.Username,
		Password: mysql.Password,
		Port:     mysql.Port,
		Dbname:   mysql.Dbname,
		Dbtype:   mysql.Dbtype,
		Prefix:   mysql.Prefix,
	}) // 指向测试数据库
	teardownTests() // 清空数据库
	model.Migrate()
}

func teardownTests() {
	model.DB.DropTableIfExists(&model.User{})
	model.DB.DropTableIfExists(&model.UserProfile{})
}
