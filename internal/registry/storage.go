package registry

import (
	"github.com/jinzhu/gorm"
	"go-sso/internal/repository"
	"go-sso/internal/repository/mysql/mysql_query"
)

var (
	s repository.Storage
)

func SetStorage(db *gorm.DB) repository.Storage {
	s = mysql_query.NewQuery(db)
	return s
}

func GetStorage() repository.Storage {
	return s
}
