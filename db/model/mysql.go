package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go-sso/conf"
	"go-sso/pkg/log"
)

var (
	DB                                                *gorm.DB
	err                                               error
	dbType, dbName, user, password, host, tablePrefix string
	c                                                 *conf.Config
)

func Migrate() {
	DB.AutoMigrate(
		&User{},
		&UserProfile{},
	)
}

func init() {
	c = conf.GetConfig()
	Init(c.MySQL)
}

// 数据库初始化
func Init(mysql conf.MySQLConfig) {
	dbType = mysql.Dbtype
	dbName = mysql.Dbname
	user = mysql.Username
	password = mysql.Password
	host = mysql.Host
	tablePrefix = mysql.Prefix

	DB, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName, ))

	if err != nil {
		log.Error(err.Error())
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if defaultTableName == tablePrefix+"casbin_rule" {
			return defaultTableName
		}
		return tablePrefix + defaultTableName
	}
	if c.Common.Debug {
		DB.LogMode(true) // 开启 sql 日志
	}
	DB.SingularTable(true) // 创建table名单数
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer DB.Close()
}
