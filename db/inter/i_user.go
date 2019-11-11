package inter

import (
	"go-sso/db/model"
)

type IUser interface {
	GetUserByAccount(username string) (user *model.User, err error)
	Get(id string) (obj *model.User, err error)
	Create(item *model.User) (obj *model.User, err error)
	CheckUser(account, password string) bool
	Exists(account, accountType string) bool
	IsValid(account, accountType string) bool
}
