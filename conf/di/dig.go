package di

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-sso/conf"
	"go-sso/internal/registry"
	"go-sso/internal/repository"
	"go-sso/internal/repository/mysql/model"
	"go-sso/internal/service/routers"
	"go-sso/pkg/email_tool"
	"go-sso/pkg/log"
	"go-sso/pkg/permission"
	"go-sso/pkg/sms"
	"go-sso/pkg/wx/wx_client"
	"go.uber.org/dig"
	"path/filepath"
)

type DigConfig struct {
	dig.In

	Config *conf.Config
	DB     *gorm.DB
	Query  repository.Storage
}

func PrintConfig(config DigConfig) {
	fmt.Printf("%+v", config)
}

func BuildContainer() *dig.Container {
	Container := dig.New()
	_ = Container.Provide(conf.InitConfig)
	_ = Container.Provide(InitConfig)
	_ = Container.Provide(initQuery)
	return Container
}

func InitConfig(config *conf.Config) *gorm.DB {
	initLogger(config)
	permission.InitPermission(config)
	email_tool.InitEmail(config)
	wx_client.InitWeixin(config)
	sms.InitAliConfig(config)
	routers.InitRouter(config)
	model.InitMysql(config)

	return model.DB
}

func initQuery(db *gorm.DB) repository.Storage {
	return registry.SetStorage(db)
}

func initLogger(config *conf.Config) {
	logPath := filepath.Join(conf.ExeDir(), "log-files/", config.Common.LogFile)
	log.InitLogger(logPath, config.Common.Level)
}
