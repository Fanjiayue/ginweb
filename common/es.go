package common

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"time"
)

type logDataEs struct {
	Topic string `json:"topic"`
	Data string	`json:"data"`
}

var EsClient *elastic.Client
var ch chan *logDataEs

func InitEs()(err error){
	url :=  viper.GetString("es.url")
	chanMaxSize := viper.GetInt("es.chan_max_size")
	EsClient,err = elastic.NewClient(elastic.SetSniff(false),elastic.SetURL(url))
	if err!=nil{
		return
	}
	fmt.Println("content to es success")
	ch = make(chan *logDataEs,chanMaxSize)
	go SendToEs()
	return
}

func SendToEsChan(msg *logDataEs){
	ch<-msg
}

//发送数据到es
func SendToEs() {
	for{
		select {
		case msg :=<-ch:
			put1,err :=EsClient.Index().
				Index(msg.Topic).
				BodyJson(msg).
				Do(context.Background())
			if err!=nil{
				fmt.Println("data send to falsed,err: ",err)
				continue
			}
			fmt.Printf("Index  %s to index %s,type %s\n",put1.Id,put1.Index,put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}

}
