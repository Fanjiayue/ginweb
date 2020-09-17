package common

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

//var (
//	tailObj *tail.Tail
//	//LogChan chan string
//	)

// TailTask : 一个日志手机项任务
type TailTask struct {
	path string
	topic string
	instance *tail.Tail
}

func NewTailTask(path ,topic string) (tailClient *TailTask){
	tailClient = &TailTask{
		path:path,
		topic:topic,
	}
	tailClient.initTail()

	return
}
func (t *TailTask)initTail() (){
	config :=tail.Config{
		ReOpen:	true,									//重新打开
		Follow: true,									//是否跟随
		Location: &tail.SeekInfo{Offset: 0,Whence: 2},	//从文件的哪个位置开始读取
		MustExist: false,								//文件不存在不报错
		Poll: true,
	}
	var err error
	t.instance, err = tail.TailFile(t.path,config)
	if err!=nil{
		fmt.Printf("tail.TailFile falsed, err: %v\n",err)
		return
	}
	fmt.Println("tail.TailFile success")
	go t.run()

}



func (t *TailTask)run(){
	for {
		select {
		case line:=<-t.instance.Lines:
			 //SendToKafka(t.topic,line.Text)//函数调用函数
			 //优化：先把日志数据发送到通道中，再在kafka中起goroutine去取日志发送到kafka中
			SendToChan(t.topic,line.Text)
		default:
			time.Sleep(time.Millisecond*500)
		}
	}
}