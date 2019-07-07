package models

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)


type User struct {
	BaseModel
	Username    string `json:"username" gorm:"type:varchar(64);unique_index;not null"`
	Password    string `json:"password" gorm:"type:varchar(128)"`
	Role sql.NullString  `json:"role"`
	UserProfile UserProfile  // OneToOne
}

type UserProfile struct {
	UserID    uint `gorm:"primary_key"`  // OneToOne
	Email     sql.NullString    `json:"email" gorm:"type:varchar(128);unique_index;sql:"default: null"`
	Telephone string    `json:"telephone" gorm:"varchar(16)"`
	LastLogin mysql.NullTime `json:"last_login"`
}


//回调自动创建用户资料
func (user *User) AfterCreate(db *gorm.DB) (err error) {
	if err = DB.Create(&UserProfile{UserID: user.ID}).Error;err != nil {
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
