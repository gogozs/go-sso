package permission

import (
	"fmt"
	"github.com/casbin/casbin"
	_ "github.com/go-sql-driver/mysql"
	"go-sso/conf"
	"go-sso/db/model"
	"go-sso/pkg/log"
	"path/filepath"
)

var (
	enforcer *casbin.Enforcer
)

func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

// 加载权限设置
func init() {
	enforcer = CreateCasbin()
}

//权限结构
type CasbinModel struct {
	model.BaseModel
	Ptype    string `json:"p_type"`
	RoleName string `json:"role_name"`
	Path     string `json:"path"`
	Method   string `json:"method"`
}


//添加权限
func AddCasbin(cm CasbinModel) bool {
	return enforcer.AddPolicy(cm.RoleName, cm.Path, cm.Method)
}

//持久化到数据库
func CreateCasbin() *casbin.Enforcer {
	var (
		user, password, host, port, dbname string
	)
	mysql := conf.GetConfig().MySQL
	user = mysql.Username
	password = mysql.Password
	host = mysql.Host
	port = mysql.Port
	dbname = mysql.Dbname

	// Initialize a Gorm adapter and use it in a Casbin enforcer:
	// The adapter will use the MySQL database named "casbin".
	// If it doesn't exist, the adapter will create it automatically.
	// You can also use an already existing gorm instance with gormadapter.NewAdapterByDB(gormInstance)
	a := NewAdapter("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		user,
		password,
		host,
		port,
		dbname,
	), true) // Your driver and data source.
	p := filepath.Join(conf.GetConfigPath(), "auth_model.conf")
	e := casbin.NewEnforcer(p, a)

	err := e.LoadPolicy()
	if err != nil {
		log.Error(err.Error())
	}
	return e
}

