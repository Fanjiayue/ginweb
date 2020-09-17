package common

import (
	"fmt"
	"time"
)

var tskMgr *tailLogMgr

type tailLogMgr struct {
	logEntry []*LogEntry
	taskMap map[string]*TailTask
	newConfChan chan []*LogEntry
}

func InitTskMgr(logEntryConf []*LogEntry){
	tskMgr = &tailLogMgr{
		logEntry:logEntryConf, //把当前日志收集的配置信息保存下来
		taskMap:make(map[string]*TailTask,16),
		newConfChan:make(chan []*LogEntry),//无缓冲区通道
	}
	for _,value := range logEntryConf{
		tailClient :=NewTailTask(value.Path, value.Topic)
		mk := fmt.Sprintf("%s_%s",value.Path,value.Topic)
		tskMgr.taskMap[mk] = tailClient
	}
	go tskMgr.run()
	go tskMgr.WatchConf()
}

func (t *tailLogMgr)run(){
	for{
		select {
		case newConf := <-t.newConfChan:
			for _,conf := range newConf{
				mk := fmt.Sprintf("%s_%s",conf.Path,conf.Topic)
				_,ok := t.taskMap[mk]
				if ok {
					fmt.Println("原来就有，无需操作")
					//原来就有，无需操作
					continue
				}else{
					//新增或修改的
					fmt.Println("新增或修改的")
					tailClient :=NewTailTask(conf.Path, conf.Topic)
					t.taskMap[mk] = tailClient
				}
			}
			//找出原来t.logEntry中有的，newConf中没有的，去删除

		default:
		 	time.Sleep(time.Second)
		}
	}
}


//向外暴露tskMgr内的newConfChan
func NewConfChan() chan<- []*LogEntry{
	return tskMgr.newConfChan
}


