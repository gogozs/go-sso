package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

type Config struct {
	MySQL     MySQLConfig
	TestMysql TestMysqlConfig
	Cache     CacheConfig
	Jwt       JwtConfig
	Cookie    CookieConfig
	Common    CommonConfig
	Email     EmailConfig
	Weixin    WeixinConfig
	Cors      CorsConfig
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

type TestMysqlConfig struct {
	Host     string
	Username string
	Password string
	Port     string
	Dbname   string
	Dbtype   string
	Prefix   string
}

type CacheConfig struct {
	Host        string
	Password    string
	Dbname      int
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type CorsConfig []string

var config Config

func GetConfig() *Config {
	return &config
}

type CommonConfig struct {
	Debug        bool
	AppSecret    string
	TemplatePath string // 静态文件相对路径

	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PageSize     int

	LogFile string
	Level   string
}

// jwt
type JwtConfig struct {
	SECRET string
	EXP    time.Duration // 过期时间
	ALG    string        // 算法
}

// cookie
type CookieConfig struct {
	NAME string
}

type EmailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Admin    []string
}

type WeixinConfig struct {
	AppID          string
	AppSecret      string
	GrantType      string
	EncodingAESKey string // 加密密钥
	Token          string //官网中配置相同
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

func init() {
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
	err = viper.Unmarshal(&config) // 将配置信息绑定到结构体上
	if err != nil {
		panic(err)
	}
}
