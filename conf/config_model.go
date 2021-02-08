package conf

import "time"

type Config struct {
	MySQL     MySQLConfig
	TestMysql MySQLConfig
	Redis     RedisConfig
	Jwt       JwtConfig
	Cookie    CookieConfig
	Common    CommonConfig
	Email     EmailConfig
	Weixin    WeixinConfig
	Cors      CorsConfig
	AliConfig AliConfig
	AliSms    []AliSmsTemplate
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

type RedisConfig struct {
	Host        string
	Port        int
	Password    string
	Dbname      int
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type CorsConfig []string

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

type AliConfig struct {
	AccessKey    string
	AccessSecret string
}

type AliSmsTemplate struct {
	TemplateName   string
	SignName       string
	TemplateId     string
	TemplateParams string
}
