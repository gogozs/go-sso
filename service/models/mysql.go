package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-weixin/config"
	"go-weixin/pkg/log"
	"time"
)

var DB *gorm.DB

type BaseModel struct {
	ID        uint      `json:"id" gorm:"primary_key" gorm:"AUTO_INCREMENT"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 数据库初始化
func init() {
	c := config.GetConfig()
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)
	mysql := c.MySQL

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

	DB.SingularTable(true) // 创建table名单数
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer DB.Close()
}
