package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-weixin/service/models"
)

func main()  {
	db, _ := gorm.Open("sqlite3", "/tmp/gorm.db")
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.UserProfile{})
	db.LogMode(true)
	user := models.User{Username: "test", Password: "test", Role: "superuser"}
	err := db.Create(&user)
	fmt.Println(err)
}
