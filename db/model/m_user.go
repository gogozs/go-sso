package model

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"go-sso/pkg/log"
	"gopkg.in/go-playground/validator.v9"
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

// 回调自动创建用户资料
func (user *User) AfterCreate(db *gorm.DB) (err error) {
	if err = db.Create(&UserProfile{UserID: user.ID}).Error; err != nil {
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
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Telephone string `json:"telephone" validate:"required"`
	Email     string `json:"email"`
}

func (this *RegisterParams) Validate() error {
	validate := validator.New()
	return validate.Struct(this)
}

type ChangePasswordParams struct {
	RawPassword string `json:"raw_password" validate:"required,gte=6"`
	NewPassword string `json:"new_password" validate:"required,gte=6"`
}

func (this *ChangePasswordParams) Validate() error {
	validate := validator.New()
	return validate.Struct(this)
}

type ResetPasswordParams struct {
	Account     string `json:"account" validate:"required"`
	VerifyType  string `json:"verify_type" validate:"required"`
	Code        string `json:"code" validate:"required,eq=6"`
	NewPassword string `json:"new_password" validate:"required,gte=6"`
}

func (this *ResetPasswordParams) Validate() error {
	validate := validator.New()
	return validate.Struct(this)
}

type TelephoneLoginParams struct {
	Telephone string `json:"telephone" validate:"required,eq=11"`
	Code      string `json:"code" validate:"required,eq=6"`
}

func (this *TelephoneLoginParams) Validate() error {
	validate := validator.New()
	return validate.Struct(this)
}
