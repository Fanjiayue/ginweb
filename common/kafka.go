package common

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

var KfkClinet sarama.SyncProducer  //声明一个全局变量连接kafka的生产者client

func InitKafka() sarama.SyncProducer{
	host :=  viper.GetString("kafka.host")
	port := viper.GetString("kafka.port")
	address := host+":"+port
	fmt.Println(address)
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	var err error
	KfkClinet, err = sarama.NewSyncProducer([]string{address}, config)
	if err != nil {
		panic("failed to connect kafka, err: " + err.Error())
	}
	fmt.Println("kafka connect success")
	return KfkClinet
}


func SendToKafka(topic, value string) (err error){
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	//msg.Key = sarama.StringEncoder("log")
	msg.Value = sarama.StringEncoder(value)

	pid, offset, err := KfkClinet.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
	return nil
}