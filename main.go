package main

import (
	"github.com/gin-gonic/gin"
	"ginweb/common"
	"os"
	"github.com/spf13/viper"
)


func main(){
	InitConfig()
	kafka := common.InitKafka()  //初始化kafka
	defer kafka.Close()
	common.InitTail()            //初始化tail
	redis :=common.InitRedis()   //初始化redis
	defer redis.Close()
	db := common.InitDB()		//初始化mysql
	defer db.Close()
	r := gin.Default()
	r = ControllerRouter(r)
	port := viper.GetString("server.port")
	if port !="" {
		panic(r.Run(":"+port))
	}
	panic(r.Run())
}


func InitConfig(){
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir+"/config")
	err := viper.ReadInConfig()
	if err!=nil{
		panic(err)
	}
}

