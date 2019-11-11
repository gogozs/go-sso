package cli

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"go-sso/db/model"
	"go-sso/db/query"
)

var createUserCmd = &cobra.Command{
	Use:   "createsuperuser",
	Short: "创建管理员用户",
	Long:  "创建管理员用户",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("请输入正确的用户名和密码")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		CreateUser(args)
	},
}

func CreateUser(args []string) {
	username := args[0]
	password := args[1]
	user := &model.User{Username: username, Password:password, Role:"superuser"}
	qu := query.UserQuery{}
	_, err := qu.Create(user)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("create user success!")
	}
}

func init() {
	rootCmd.AddCommand(createUserCmd)
}