package registry

import (
	"github.com/jinzhu/gorm"
	"go-sso/storage"
	"go-sso/storage/mysql/mysql_query"
)

var (
	s storage.Storage
)

func SetStorage(db *gorm.DB) storage.Storage {
	s = mysql_query.NewQuery(db)
	return s
}

func GetStorage() storage.Storage {
	return s
}
