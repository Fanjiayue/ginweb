package common

import (
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
)

var (
	tailObj *tail.Tail
	//LogChan chan string
	)


func InitTail() *tail.Tail{
	fileName := viper.GetString("tail.filename")
	config :=tail.Config{
		ReOpen:	true,									//重新打开
		Follow: true,									//是否跟随
		Location: &tail.SeekInfo{Offset: 0,Whence: 2},	//从文件的哪个位置开始读取
		MustExist: false,								//文件不存在不报错
		Poll: true,
	}
	var err error
	tailObj, err = tail.TailFile(fileName,config)
	if err != nil{
		panic("tail file falsed,err :"+err.Error())
	}

	return tailObj
}

func GetTailObj() *tail.Tail{
	return tailObj
}

func ReadChan() <-chan *tail.Line{
	return tailObj.Lines

}