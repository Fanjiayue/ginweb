package common

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"time"
)

type logData struct {
	topic string
	data string
}



var (
	KfkClinet sarama.SyncProducer  //声明一个全局变量连接kafka的生产者client
	logDataChan chan *logData
)

func InitKafka() sarama.SyncProducer{
	host :=  viper.GetString("kafka.host")
	port := viper.GetString("kafka.port")
	chanMaxSize := viper.GetInt("kafka.chan_max_size")
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
	logDataChan = make(chan *logData,chanMaxSize)  //初始化日志通道
	go sendToKafka()//开启后台的goroutine，将channel中的数据取出来发往kafka
	return KfkClinet
}

//将数据发送到通道中
func SendToChan(topic, data string){
	msg := &logData{
		topic:topic,
		data:data,
	}
	logDataChan<-msg
}
//将数据发送到kafka
func sendToKafka() {
	for{
		select {
		case logmsg :=<-logDataChan:
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = logmsg.topic
			//msg.Key = sarama.StringEncoder("log")
			msg.Value = sarama.StringEncoder(logmsg.data)

			pid, offset, err := KfkClinet.SendMessage(msg)
			if err != nil {
				fmt.Printf("kafka sendmessage falsed, err: %v\n",err)
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
		default:
			time.Sleep(time.Millisecond*50)
		}

	}

}

////将数据发送到kafka
//func SendToKafka(topic, value string) {
//	// 构造一个消息
//	msg := &sarama.ProducerMessage{}
//	msg.Topic = topic
//	//msg.Key = sarama.StringEncoder("log")
//	msg.Value = sarama.StringEncoder(value)
//
//	pid, offset, err := KfkClinet.SendMessage(msg)
//	if err != nil {
//		fmt.Printf("kafka sendmessage falsed, err: %v\n",err)
//		return
//	}
//	fmt.Printf("pid:%v offset:%v\n", pid, offset)
//	return
//}

func CustomerKakfka(topic string) error{
	host :=  viper.GetString("kafka.host")
	port := viper.GetString("kafka.port")
	address := host+":"+port
	consumer, err := sarama.NewConsumer([]string{address}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return err
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return err
	}

	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return err
		}
		fmt.Println("consumer from kafka started success")
		//defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				//发送给es
				le := logDataEs{Topic: topic,Data: string(msg.Value)}
				SendToEsChan(&le)
			}

		}(pc)
	}
	return err
}
