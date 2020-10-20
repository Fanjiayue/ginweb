package common

import (
	"fmt"
	"ginweb/model"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)
var DB *gorm.DB

func InitDB() *gorm.DB{
	driverName := viper.GetString("datasource.driverName")
	host :=  viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True", username, password, host, port,database,charset)
	var err error
	DB, err = gorm.Open(driverName,args)
	if err!=nil{
		panic("failed to connect database, err: " + err.Error())
	}

	maxIdleConns := viper.GetInt("datasource.maxIdleConns")
	maxOpenConns := viper.GetInt("datasource.maxOpenConns")
	DB.DB().SetMaxIdleConns(maxIdleConns)//SetMaxIdleConns用于设置闲置的连接数
	DB.DB().SetMaxOpenConns(maxOpenConns)//SetMaxOpenConns用于设置最大打开的连接数

	// 启用Logger，显示详细日志
	//db.LogMode(true)

	// 自动迁移模式
	//db.AutoMigrate(&Model.UserModel{},
	//	&Model.UserDetailModel{},
	//	&Model.UserAuthsModel{},
	//)
	DB.AutoMigrate(&model.User{})

	fmt.Println("database connect success")

	return DB

}

func GetDB() *gorm.DB{
	return DB
}
