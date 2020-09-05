package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jinzhu/gorm"
	"ginweb/model"
	"ginweb/common"
	"ginweb/util"
	"golang.org/x/crypto/bcrypt"
	"log"
	"ginweb/response"
	"github.com/gomodule/redigo/redis"
	_"github.com/go-sql-driver/mysql"
)
//创建用户
func Register(ctx *gin.Context) {
	DB := common.GetDB()
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
		name = util.RandomString(10)
	}

	//判断手机是否存在
	if isTelephoneExist(DB, telephone){
		response.Error(ctx,10002,"	手机号已存在")
		return
	}

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err!=nil{
		response.Response(ctx,http.StatusInternalServerError,500,nil,"加密错误")
		return
	}
	//创建用户
	newuser := model.User{
		Name : name,
		Telephone:telephone,
		Password:string(hasedPassword),
	}
	DB.Create(&newuser)

	//返回结果
	response.Success(ctx,nil,"注册成功")
	return
}

//查询手机号是否存在
func isTelephoneExist(db *gorm.DB,telephone string) bool{
	var user model.User
	db.Where("telephone = ?",telephone).First(&user)
	if user.ID !=0{
		return true
	}
	return false
}

//登陆
func Login(ctx *gin.Context){
	DB := common.GetDB()
	telephone := ctx.PostForm("telephone")
	password :=ctx.PostForm("password")

	var user model.User
	DB.Where("telephone = ?",telephone).First(&user)

	if user.ID ==0{
		response.Error(ctx,40002,"	手机号不存在")
		return
	}


	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password)); err !=nil{
		response.Error(ctx,40001,"	密码错误")
		return
	}

	token,err :=common.ReleaseToken(user)
	if err != nil{
		response.Response(ctx,http.StatusInternalServerError,500,nil,"系统错误")
		log.Printf("token genrate error: %v",err)
		return
	}


	response.Success(ctx,gin.H{"token":token,},"查询成功")
	return

}

func Userinfo(ctx *gin.Context){
	user,_ := ctx.Get("user")
	response.Success(ctx,gin.H{"user":response.ToUserResponse(user.(model.User))},"查询成功")
	return
}

func Test(ctx *gin.Context){
	Redis := common.GetRedis()
	_, err :=Redis.Do("SET","ID1","no1")
	if err!=nil{
		response.Response(ctx,http.StatusInternalServerError,500,nil,"系统错误")
		log.Printf("redis genrate error: %v",err)
		return
	}
	v,err:=redis.String(Redis.Do("GET","key"))

	response.Success(ctx,gin.H{"token":v,},"查询成功")
	return
}
