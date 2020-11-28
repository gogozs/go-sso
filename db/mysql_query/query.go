package mysql_query

import (
	"github.com/jinzhu/gorm"
	"go-sso/conf"
	"go-sso/db/model"
	"strconv"
)

func NewQuery(db *gorm.DB) *query {
	return &query{db: db}
}

type query struct {
	db *gorm.DB
}

func (q *query) GetID(id string) uint {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0
	}
	return uint(idInt)
}

func (q *query) InitPageParams(pl ...model.Pagination) (p model.Pagination) {
	if len(pl) == 0 {
		p = model.Pagination{
			PageNum:  1,
			PageSize: conf.GetConfig().Common.PageSize,
		}
	} else {
		p = pl[0]
	}
	return
}
