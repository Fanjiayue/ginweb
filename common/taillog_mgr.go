package common

var tskMgr *tailLogMgr

type tailLogMgr struct {
	logEntry []*LogEntry
}

func InitTskMgr(logEntryConf []*LogEntry){
	tskMgr = &tailLogMgr{
		logEntry:logEntryConf, //把当前日志收集的配置信息保存下来
	}
	for _,value := range logEntryConf{
		NewTailTask(value.Path, value.Topic)
	}
}