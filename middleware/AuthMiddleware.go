package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"ginweb/common"
	"ginweb/model"
	"ginweb/response"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		tokenString :=ctx.GetHeader("Authorization")
		//验证token
		if tokenString == "" || !strings.HasPrefix(tokenString,"Bearer "){
			response.Error(ctx,40000,"未登陆")
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]
		token,claims,err := common.ParseToken(tokenString)
		if err !=nil || !token.Valid{
			response.Error(ctx,40000,"未登陆")
			ctx.Abort()
			return
		}

		//获取token中信息
		userId := claims.UserId
		DB :=common.GetDB()
		var user model.User
		DB.First(&user,userId)

		//用户
		if user.ID == 0{
			response.Error(ctx,40002,"用户已删除")
			ctx.Abort()
			return
		}

		//存用户信息
		ctx.Set("user",user)

		ctx.Next()
	}
}
