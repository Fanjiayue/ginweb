package main

import (
	"github.com/gin-gonic/gin"
	"ginweb/common"
	"os"
	"github.com/spf13/viper"
)


func main(){
	kafka := common.InitKafka()
	defer kafka.Close()
	common.InitTail()
	InitConfig()
	redis :=common.InitRedis()
	defer redis.Close()
	db := common.InitDB()
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

