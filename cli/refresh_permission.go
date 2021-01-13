package cli

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"go-sso/conf"
	"os"
	"path"
)

var permissionCmd = &cobra.Command{
	Use:   "refresh_permission",
	Short: "refresh  permission",
	Long:  `refresh  permission`,
	Run: func(cmd *cobra.Command, args []string) {
		refreshPermission()
	},
}

func refreshPermission() {
	confPath := conf.GetConfigPath()
	csvFile, err := os.Open(path.Join(confPath, "auth_policy.csv"))
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(csvFile)
	rows, err := csvReader.ReadAll() // `rows` is of type [][]string
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(rows)
	var (
		user, password, host, port, dbname string
	)
	mysqlArr := []conf.MySQLConfig{
		conf.GetConfig().MySQL,
		conf.GetConfig().TestMysql,
	}
	for _, mysql := range mysqlArr {
		user = mysql.Username
		password = mysql.Password
		host = mysql.Host
		port = mysql.Port
		dbname = mysql.Dbname
		table := mysql.Prefix + "casbin_rule"

		db, err := sql.Open("mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
				user, password, host, port, dbname,
			))
		if err != nil {
			fmt.Println("failed to open database:", err.Error())
			return
		}
		defer db.Close()
		for _, row := range rows {
			s1 := fmt.Sprintf(
				"INSERT INTO %s (p_type, v0, v1, v2) values ('%s', '%s', '%s', '%s');",
				table, row[0], row[1], row[2], row[3])
			fmt.Println(s1)
			_, err := db.Exec(s1)
			fmt.Println(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(permissionCmd)
}
