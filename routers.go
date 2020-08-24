package main

import (
	"github.com/gin-gonic/gin"
	"ginweb/controller"
)

func ControllerRouter(r *gin.Engine) *gin.Engine{
	r.POST("/api/auth/register", controller.Register)
	return r
}
