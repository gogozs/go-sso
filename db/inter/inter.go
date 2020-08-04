package inter

import (
	"github.com/jinzhu/gorm"
	"go-sso/db/mysql_query"
)

type IQuery interface {
	IUser
}

var (
	q IQuery
)

func InitQuery(db *gorm.DB) IQuery {
	q = mysql_query.NewQuery(db)
	return q
}

func GetQuery() IQuery {
	return q
}
