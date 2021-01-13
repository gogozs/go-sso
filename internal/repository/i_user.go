package repository

import (
	"go-sso/internal/repository/mysql/model"
)

type UserStorage interface {
	GetUserByAccount(account string) (user *model.User, err error)
	Get(id string) (obj *model.User, err error)
	Create(item *model.User) (obj *model.User, err error)
	CheckUser(account, password string) (*model.User, bool)
	Exists(account, accountType string) bool
	IsValid(account, accountType string) bool
	ChangePassword(*model.User, string) error
}
