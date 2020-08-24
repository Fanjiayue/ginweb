package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jinzhu/gorm"
	"ginweb/model"
	"ginweb/common"
	"ginweb/util"
)

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
		ctx.JSON(http.StatusUnprocessableEntity,gin.H{
			"code":422,
			"msg":"手机号已存在",
		})
		return
	}


	//创建用户
	newuser := model.User{
		Name : name,
		Telephone:telephone,
		Password:password,
	}
	DB.Create(&newuser)

	//返回结果
	ctx.JSON(http.StatusCreated,gin.H{
		"msg":"注册成功",
	})
	return
}


func isTelephoneExist(db *gorm.DB,telephone string) bool{
	var user model.User
	db.Where("telephone = ?",telephone).First(&user)
	if user.ID !=0{
		return true
	}
	return false
}