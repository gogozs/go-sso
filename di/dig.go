package di

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-sso/conf"
	"go-sso/db/inter"
	"go-sso/db/model"
	"go-sso/pkg/email_tool"
	"go-sso/pkg/log"
	"go-sso/pkg/sms"
	"go-sso/pkg/wx/wx_client"
	"go-sso/service/api/routes"
	"go.uber.org/dig"
	"path/filepath"
)

type DigConfig struct {
	dig.Out

	Config *conf.Config
	DB     *gorm.DB
}

func PrintConfig(d *DigConfig) {
	fmt.Println(d)
}

func BuildContainer() *dig.Container {
	Container := dig.New()
	Container.Provide(conf.InitConfig)
	Container.Provide(initMysql)
	Container.Provide(initQuery)
	return Container
}

func InitConfig(config *conf.Config) {
	initLogger(config)
	email_tool.InitEmail(config)
	wx_client.InitWeixin(config)
	sms.InitAliConfig(config)
	routes.InitRouter(config)
}

func initQuery(db *gorm.DB) {
	inter.InitQuery(db)
}

func initMysql(config *conf.Config) *gorm.DB {
	model.InitMysql(config)
	return model.DB
}

func initLogger(config *conf.Config) {
	logPath := filepath.Join(conf.ExeDir(), "log-files/", config.Common.LogFile)
	log.InitLogger(logPath, config.Common.Level)
}
