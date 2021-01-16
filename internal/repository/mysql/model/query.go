package model

import (
	"github.com/jinzhu/gorm"
	"go-sso/conf"
	"strconv"
)

func NewQuery(db *gorm.DB) *query {
	return &query{db: db}
}

type query struct {
	db *gorm.DB
}

func (q query) GetID(id string) uint {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0
	}
	return uint(idInt)
}

func (q query) InitPageParams(pl ...Pagination) (p Pagination) {
	if len(pl) == 0 {
		p = Pagination{
			PageNum:  1,
			PageSize: conf.GetConfig().Common.PageSize,
		}
	} else {
		p = pl[0]
	}
	return
}
