package query

import (
	"github.com/jinzhu/gorm"
	"go-qiuplus/conf"
	"go-qiuplus/db/model"
	"strconv"
)

var (
	UserQ      *UserQuery
)

func init() {
	UserQ = &UserQuery{BaseQuery{db: model.DB}}
}

type BaseQuery struct {
	db *gorm.DB
}

func (this *BaseQuery) GetID(id string) uint {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0
	}
	return uint(idInt)
}

func (this *BaseQuery) InitPageParams(pl ...model.Pagination) (p model.Pagination) {
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