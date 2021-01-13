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
)

var (
	TableArr = []interface{}{
		&User{},
		&UserProfile{},
	}
)

func Migrate() {
	DB.AutoMigrate(TableArr...)
}

type MySQLConfig struct {
	Host     string
	Username string
	Password string
	Port     string
	Dbname   string
	Dbtype   string
	Prefix   string
}

// 数据库初始化
func InitMysql(config *conf.Config) {
	mysql := config.MySQL
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
		dbName))

	if err != nil {
		log.Error(err.Error())
	}
	initDBConfig(config.Common.Debug)
}

func initDBConfig(debug bool) {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if defaultTableName == tablePrefix+"casbin_rule" {
			return defaultTableName
		}
		return tablePrefix + defaultTableName
	}
	if debug {
		DB.LogMode(true) // 开启 sql 日志
	}
	DB.SingularTable(true) // 创建table名单数
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer DB.Close()
}
