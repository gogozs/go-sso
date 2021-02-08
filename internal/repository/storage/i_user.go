package storage

import (
	"go-sso/internal/repository/storage/mysql"
)

type UserStorage interface {
	GetUserByAccount(account string) (user *mysql.User, err error)
	Get(id string) (obj *mysql.User, err error)
	Create(item *mysql.User) (obj *mysql.User, err error)
	GetUser(account, password string) (*mysql.User, bool)
	Exists(account, accountType string) bool
	IsValid(account, accountType string) bool
	ChangePassword(*mysql.User, string) error
}
