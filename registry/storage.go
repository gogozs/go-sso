package registry

import (
	"github.com/jinzhu/gorm"
	"go-sso/internal/repository/storage"
	"go-sso/internal/repository/storage/mysql"
)

var (
	s storage.Storage
)

func InitStorage(db *gorm.DB) storage.Storage {
	s = mysql.NewQuery(db)
	return s
}

func GetStorage() storage.Storage {
	return s
}
