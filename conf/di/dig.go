package di

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go-sso/conf"
	"go-sso/internal/repository"
	"go-sso/internal/repository/mysql/model"
	"go-sso/internal/routers"
	"go-sso/pkg/email_tool"
	"go-sso/pkg/log"
	"go-sso/pkg/permission"
	"go-sso/pkg/sms"
	"go-sso/pkg/wx/wx_client"
	"go-sso/registry"
	"go.uber.org/dig"
	"path/filepath"
)

type DigConfig struct {
	dig.In

	Config  *conf.Config
	DB      *gorm.DB
	Storage repository.Storage
	Engine  *gin.Engine `optional:"true"`
}

func RunServer(config DigConfig) {
	initLogger(config.Config)
	permission.InitPermission(config.Config)
	email_tool.InitEmail(config.Config)
	wx_client.InitWeixin(config.Config)
	sms.InitAliConfig(config.Config)
}

func BuildContainer() *dig.Container {
	Container := dig.New()
	_ = Container.Provide(conf.InitConfig)
	_ = Container.Provide(model.InitMysql)
	_ = Container.Provide(registry.SetStorage)
	_ = Container.Provide(routers.InitRouter)
	return Container
}

func initLogger(config *conf.Config) {
	logPath := filepath.Join(conf.ExeDir(), "log-files/", config.Common.LogFile)
	log.InitLogger(logPath, config.Common.Level)
}
