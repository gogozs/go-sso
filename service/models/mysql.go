package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go-weixin/config"
	"go-weixin/pkg/log"
	"time"
)

var (
	DB                                                *gorm.DB
	err                                               error
	dbType, dbName, user, password, host, tablePrefix string
	c                                                 *config.Config
)

func Migrate() {
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&UserProfile{})
}

type BaseModel struct {
	ID        uint      `json:"id" gorm:"primary_key" gorm:"AUTO_INCREMENT"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func init() {
	c = config.GetConfig()
	Init(c.MySQL)
}

// 数据库初始化
func Init(mysql config.MySQLConfig) {
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
		log.Info(err)
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
