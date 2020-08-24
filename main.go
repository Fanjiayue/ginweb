package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"math/rand"
	"time"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null;unique"`
	Password string `gorm:"size:255;not null"`

}

func main(){
	db := InitDB()
	defer db.Close()
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		//获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password :=ctx.PostForm("password")

		//验证手机号
		if len(telephone)!= 11{
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{
				"code":422,
				"msg":"请正确填写手机号",
			})
			return
		}
		//验证密码
		if len(password)< 6{
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{
				"code":422,
				"msg":"密码不能少于6位",
			})
			return
		}
		//验证名称，如果没有穿就随机给一个10位的字符串
		if len(name)==0{
			name = RandomString(10)
		}

		//判断手机是否存在
		if isTelephoneExist(db, telephone){
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{
				"code":422,
				"msg":"手机号已存在",
			})
			return
		}


		//创建用户
		newuser := User{
			Name : name,
			Telephone:telephone,
			Password:password,
		}
		db.Create(&newuser)

		//返回结果
		ctx.JSON(http.StatusCreated,gin.H{
			"msg":"注册成功",
		})
		return
	})
	r.Run()
}

func RandomString(n int) string{
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte,n)
	rand.Seed(time.Now().Unix())
	for i := range result{
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

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
	db.AutoMigrate(&User{})
	return db

}

func isTelephoneExist(db *gorm.DB,telephone string) bool{
	var user User
	db.Where("telephone = ?",telephone).First(&user)
	if user.ID !=0{
		return true
	}
	return false
}