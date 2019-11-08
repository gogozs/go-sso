package model

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-qiuplus/pkg/log"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

var (
	AnonymousUser = User{}
)

type User struct {
	BaseModel
	Username    string         `json:"username" gorm:"type:varchar(64);unique_index;not null"`
	Password    string         `json:"password" gorm:"type:varchar(128)"`
	Telephone   string         `json:"telephone" gorm:"varchar(16)";unique_index`
	WxId        sql.NullString `json:"wx_id gorm:"type:varchar(64);unique_index`
	Email       string         `json:"email" gorm:"type:varchar(128);unique_index;sql:"default: ''"`
	Role        string         `json:"role" gorm:"default:'viewer'"`
	UserProfile UserProfile    // OneToOne
}

type UserProfile struct {
	UserID    uint           `gorm:"primary_key"` // OneToOne
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	LastLogin mysql.NullTime `json:"last_login"`
}

//回调自动创建用户资料
func (user *User) AfterCreate(db *gorm.DB) (err error) {
	if err = DB.Create(&UserProfile{UserID: user.ID}).Error; err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

type UserParams struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type RegisterParams struct {
	Username  string `json:"account" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Telephone string `json:"telephone" validate:"required"`
	Email     string `json:"email"`
	ValidCode string `json:"valid_code" validate:"required"`
}

func usernameFunc(f validator.FieldLevel) bool {
	if v := f.Field().String(); v == "invalid" {
		return false
	} else if b, err := regexp.MatchString(`^1([38][0-9]|14[57]|5[^4])\d{8}$`, v); !b || err != nil {
		return false
	}
	return true
}

func emailFunc(f validator.FieldLevel) bool {
	if v := f.Field().String(); v == "invalid" {
		return false
	} else if v == "" {
		return true
	} else if b, err := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, v); !b || err != nil {
		return false
	}
	return true
}

func telephoneFunc(f validator.FieldLevel) bool {
	if v := f.Field().String(); v == "invalid" {
		return false
	} else if b, err := regexp.MatchString(`^1([38][0-9]|14[57]|5[^4])\d{8}$`, v); !b || err != nil {
		return false
	}
	return true
}

func (this *RegisterParams) Validate() error {
	validate := validator.New()
	//validate.RegisterValidation("username", usernameFunc)
	//validate.RegisterValidation("telephone", telephoneFunc)
	//validate.RegisterValidation("email", emailFunc)
	return validate.Struct(this)
}
