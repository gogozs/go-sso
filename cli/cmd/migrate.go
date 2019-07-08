package cmd

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"go-weixin/config"
	"go-weixin/service/models"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Short: "migrate",
	Long: `database migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		Migrate()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}


func Migrate()  {
	c := config.GetConfig()
	var (
		dbType, dbName, user, password, host, tablePrefix string
	)
	mysql := c.MySQL
	dbType = mysql.Dbtype
	dbName = mysql.Dbname
	user = mysql.Username
	password = mysql.Password
	host = mysql.Host
	tablePrefix = mysql.Prefix

	var db *gorm.DB
	db, _ = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName, ))
	db.SingularTable(true) // 表名单数

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if defaultTableName == tablePrefix+"casbin_rule" {
			return defaultTableName
		}
		return tablePrefix + defaultTableName
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.UserProfile{})
}
