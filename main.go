package main

import (
	"fmt"
	"ginweb/common"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)


func main(){
	InitConfig()
	kafka := common.InitKafka()  //初始化kafka
	defer  kafka.Close()
	etcdClient := common.InitEtcd() //初始化etcd
	defer  etcdClient.Close()
	logEntryConf, err :=common.GetConf()
	if err != nil {
		fmt.Printf("common.GetConf failed,err:%v\n",err)
		return
	}
	fmt.Printf("get conf from etcd success, %v\n",logEntryConf)
	common.InitTskMgr(logEntryConf)


	//redis :=common.InitRedis()   //初始化redis
	//defer redis.Close()
	//db := common.InitDB()		//初始化mysql
	//defer db.Close()
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

