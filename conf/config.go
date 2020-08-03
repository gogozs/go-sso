package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"go-sso/db/model"
	"go-sso/pkg/log"
	"os"
	"path"
	"path/filepath"
)

func ExeDir() string {
	dir, exists := os.LookupEnv("GO_SSO_WORKDIR")
	if exists {
		return dir
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := path.Dir(ex)
		return exPath
	}
}

func GetConfigPath() string {
	basePath := ExeDir()
	confPath := path.Join(basePath, "conf/")
	return confPath
}

func InitConfig() error {
	// 需要配置项目根目录的环境变量，方便执行test
	confPath := GetConfigPath()
	fmt.Println("配置目录: ", confPath)
	serviceEnv := os.Getenv("service_env")
	if serviceEnv == "" {
		serviceEnv = "local"
	}
	viper.SetConfigName(serviceEnv) // 设置配置文件名 (不带后缀)
	viper.AddConfigPath(confPath)   // 第一个搜索路径
	viper.WatchConfig()             // 监控配置文件热重载
	err := viper.ReadInConfig()     // 读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	initMysql()
	initLogger()

	return nil
}

func initMysql() {
	model.InitMysql(config.Common.Debug, model.MySQLConfig{
		Host:     config.MySQL.Host,
		Username: config.MySQL.Username,
		Password: config.MySQL.Password,
		Port:     config.MySQL.Port,
		Dbname:   config.MySQL.Dbname,
		Dbtype:   config.MySQL.Dbtype,
		Prefix:   config.MySQL.Prefix,
	})
}

func initLogger() {
	logPath := filepath.Join(ExeDir(), "log-files/", config.Common.LogFile)
	log.InitLogger(logPath, config.Common.Level)
}
