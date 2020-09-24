package common

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

//var (
//	tailOb *tail.Tail
//	//LogChan chan string
//	)

// TailTask : 一个日志手机项任务
type TailTask struct {
	path string
	topic string
	instance *tail.Tail
	//为了能够推出t.run()
	ctx context.Context
	cannelFunc context.CancelFunc
}

func NewTailTask(path ,topic string) (tailClient *TailTask){
	ctx, cannel :=context.WithCancel(context.Background())
	tailClient = &TailTask{
		path:path,
		topic:topic,
		ctx:ctx,
		cannelFunc:cannel,
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
		case <-t.ctx.Done():
			fmt.Printf("tail task:%s_%s 退出",t.path,t.topic)
			return
		case line:=<-t.instance.Lines:
			 //SendToKafka(t.topic,line.Text)//函数调用函数
			 //优化：先把日志数据发送到通道中，再在kafka中起goroutine去取日志发送到kafka中
			SendToChan(t.topic,line.Text)
		default:
			time.Sleep(time.Millisecond*500)
		}
	}
}