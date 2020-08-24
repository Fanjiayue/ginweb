package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"ginweb/common"
)


func main(){
	db := common.InitDB()
	defer db.Close()
	r := gin.Default()
	r = ControllerRouter(r)
	r.Run()
}




