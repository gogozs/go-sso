package models

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseModel
	Username    string `json:"username" gorm:"type:varchar(64);unique_index;not null"`
	Password    string `json:"password" gorm:"type:varchar(128)"`
	Role        string `json:"role";"default:viewer"`
	UserProfile UserProfile // OneToOne
}

type UserProfile struct {
	UserID    uint           `gorm:"primary_key"` // OneToOne
	Email     sql.NullString `json:"email" gorm:"type:varchar(128);unique_index;sql:"default: null"`
	Telephone string         `json:"telephone" gorm:"varchar(16)"`
	LastLogin mysql.NullTime `json:"last_login"`
}

//回调自动创建用户资料
func (user *User) AfterCreate(db *gorm.DB) (err error) {
	if err = DB.Create(&UserProfile{UserID: user.ID}).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetUser(username string) (user User, err error) {
	if err = DB.Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Println(err)
		return user, err
	} else {
		return user, err
	}
}

func CreateUser(user  User) error {
	if user.Password != "" {
		var err error
		user.Password, err = GeneratePassword(user.Password)
		if err != nil {
			return err
		}
	}
	u := DB.Create(&user) // gorm会自动将id插入struct中
	if u.Error != nil {
		return u.Error
	}
	return nil
}

func GeneratePassword(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(userPwd, hashPwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(userPwd))
	return err
}

func CheckUser(username, password string) bool {
	user, err :=  GetUser(username)
	if err != nil {
		return false
	} else {
		if err = ComparePassword(password, user.Password); err != nil {
			return false
		}
	}
	return true
}
