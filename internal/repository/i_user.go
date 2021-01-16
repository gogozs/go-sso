package repository

import (
	m "go-sso/internal/repository/mysql/model"
)

type UserStorage interface {
	GetUserByAccount(account string) (user *m.User, err error)
	Get(id string) (obj *m.User, err error)
	Create(item *m.User) (obj *m.User, err error)
	GetUser(account, password string) (*m.User, bool)
	Exists(account, accountType string) bool
	IsValid(account, accountType string) bool
	ChangePassword(*m.User, string) error
}
