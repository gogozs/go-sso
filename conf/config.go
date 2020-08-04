package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
)

var config Config

func GetConfig() *Config {
	return &config
}

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

func InitConfig() (*Config, error) {
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
		return nil, err
	}

	return &config, nil
}
