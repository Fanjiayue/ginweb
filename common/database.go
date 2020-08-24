package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"ginweb/model"
)
var DB *gorm.DB

func InitDB() *gorm.DB{
	driverName := "mysql"
	host := "127.0.0.1"
	port := "3306"
	database := "ginweb"
	username := "root"
	password := "Fjy13819372207"
	charset := "utf8mb4"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True", username, password, host, port,database,charset)
	db, err := gorm.Open(driverName,args)
	if err!=nil{
		panic("failed to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&model.User{})

	DB = db
	return db

}

func GetDB() *gorm.DB{
	return DB
}
