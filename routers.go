package main

import (
	"github.com/gin-gonic/gin"
	"ginweb/controller"
	"ginweb/middleware"
)

func ControllerRouter(r *gin.Engine) *gin.Engine{
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("api/auth/userinfo",middleware.AuthMiddleware(),controller.Userinfo)
	return r
}
