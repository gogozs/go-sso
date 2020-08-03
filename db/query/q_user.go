package query

import (
	"go-sso/db/model"
	"go-sso/util"
	"regexp"
)

// 校验用户名是否合法 1.不能有 @  2.不能纯数字
func (this *query) IsValid(account, accountType string) bool {
	switch accountType {
	case "telephone":
		if b, err := regexp.MatchString(`^1([38][0-9]|14[57]|5[^4])\d{8}$`, account); !b || err != nil {
			return false
		}
	case "username":
		// 非纯数字
		if b, err := regexp.MatchString(`^[a-zA-Z]+[a-zA-Z0-9_]{2,15}$`, account); !b || err != nil {
			return false
		}
	case "email":
		if b, err := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, account); !b || err != nil {
			return false
		}
	}
	return true
}

// 检查是否被占用
func (this *query) Exists(account, accountType string) bool {
	var user model.User
	switch accountType {
	case "telephone":
		if err := model.DB.Where("telephone = ?", account).Find(&user).Error; err == nil {
			return true
		}
	case "username":
		if err := model.DB.Where("username = ?", account).Find(&user).Error; err == nil {
			return true
		}
	case "email":
		if err := model.DB.Where("email = ?", account).Find(&user).Error; err == nil {
			return true
		}
	}
	return false
}

func (this *query) GetUserByAccount(account string) (obj *model.User, err error) {
	obj = &model.User{}
	if err = model.DB.Where("username = ? OR telephone = ? OR email = ?",
		account, account, account).First(obj).Error; err != nil {
		return
	} else {
		return
	}
}

func (this *query) Get(id string) (obj *model.User, err error) {
	obj = &model.User{}
	err = this.db.Where("ID = ?", this.GetID(id)).First(obj).Error
	return
}

func (this *query) Create(item *model.User) (obj *model.User, err error) {
	if item.Password != "" {
		var err error
		item.Password, err = util.GeneratePassword(item.Password)
		if err != nil {
			return nil, err
		}
	}
	err = model.DB.Create(&item).Error // gorm会自动将id插入struct中
	return item, err
}

func (this *query) CheckUser(account, password string) (*model.User, bool) {
	user, err := this.GetUserByAccount(account)
	if err != nil {
		return nil, false
	} else {
		if err = util.ComparePassword(password, user.Password); err != nil {
			return nil, false
		}
	}
	return user, true
}

func (this *query) ChangePassword(u *model.User, newPassword string) (err error) {
	u.Password, err = util.GeneratePassword(newPassword)
	if err != nil {
		return
	}
	err = model.DB.Model(u).Update(*u).Error
	return
}
