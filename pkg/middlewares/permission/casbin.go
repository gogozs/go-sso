package permission

import (
	"encoding/csv"
	"fmt"
	"github.com/casbin/casbin"
	_ "github.com/go-sql-driver/mysql"
	"go-weixin/config"
	"go-weixin/service/models"
	"log"
	"os"
)

// 加载权限设置
func init() {
	csvFile, err := os.Open("config/auth_policy.csv")
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(csvFile)
	rows, err := csvReader.ReadAll() // `rows` is of type [][]string
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for _, row := range rows {
		cm := CasbinModel{
			Ptype: row[0],
			RoleName: row[1],
			Path: row[2],
			Method: row[3],
		}
		AddCasbin(cm)
	}
}

//权限结构
type CasbinModel struct {
	models.BaseModel
	Ptype    string `json:"p_type"`
	RoleName string `json:"role_name"`
	Path     string `json:"path"`
	Method   string `json:"method"`
}


//添加权限
func AddCasbin(cm CasbinModel) bool {
	e := Casbin()
	return e.AddPolicy(cm.RoleName, cm.Path, cm.Method)
}

//持久化到数据库
func Casbin() *casbin.Enforcer {
	var (
		err                        error
		user, password, host, port string
	)

	mysql := config.GetConfig().MySQL

	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	user = mysql.Username
	password = mysql.Password
	host = mysql.Host
	port = mysql.Port

	// Initialize a Gorm adapter and use it in a Casbin enforcer:
	// The adapter will use the MySQL database named "casbin".
	// If it doesn't exist, the adapter will create it automatically.
	// You can also use an already existing gorm instance with gormadapter.NewAdapterByDB(gormInstance)
	a := NewAdapter("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		user,
		password,
		host,
		port,
	)) // Your driver and data source.
	enforcer := casbin.NewEnforcer("config/auth_model.conf", a)
	enforcer.LoadPolicy()
	return enforcer
}

//func AddCasbin(c *gin.Context) {
//	rolename := c.PostForm("rolename")
//	path := c.PostForm("path")
//	method := c.PostForm("method")
//	ptype := "p"
//	casbin := CasbinModel{
//		Ptype:    ptype,
//		RoleName: rolename,
//		Path:     path,
//		Method:   method,
//	}
//	isok := casbins.AddCasbin(casbin)
//	if isok {
//		c.JSON(http.StatusOK, gin.H{
//			"success": true,
//			"msg":     "保存成功",
//		})
//	} else {
//		c.JSON(http.StatusOK, gin.H{
//			"success": false,
//			"msg":     "保存失败",
//		})
//	}
//}
