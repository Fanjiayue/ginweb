package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

type LogEntry struct {
	Path string `json:"path"`  //日志存放路径
	Topic string `json:"topic"`//日志发往Kafka中的哪个topic
}

var EtcdClient *clientv3.Client

func InitEtcd() *clientv3.Client{
	host := viper.GetString("etcd.host")
	port := viper.GetString("etcd.port")
	hosts := strings.Split(host,",")
	ports := strings.Split(port,",")
	n := len(hosts)
	if n>10{
		panic("集群超过规定长度")
	}
	m := len(ports)
	if m!=n{
		panic("配置文件host和port长度不一致")
	}
	var address = []string{"","","","","","","","","",""}
	for key, _ := range hosts {
		address[key] = hosts[key]+":"+ports[key]
	}
	var err error
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   address,		//连接地址，多个
		DialTimeout: 5 * time.Second,                    //超时时间
	})
	if err != nil {
		panic("failed to connect etcd, err: " + err.Error())
	}
	fmt.Println("etcd connect success")
	return EtcdClient
}

func GetConf() (logEntryConf []*LogEntry,err error) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := EtcdClient.Get(ctx, viper.GetString("etcd.key"))
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		//fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &logEntryConf)
		if err != nil {
			fmt.Printf("Unmarshal etcd value failed, err:%v\n", err)
			return
		}
	}
	return
}

func PutConf(key string,value string) (err error){
	host := viper.GetString("etcd.host")
	port := viper.GetString("etcd.port")
	hosts := strings.Split(host,",")
	ports := strings.Split(port,",")
	n := len(hosts)
	if n>10{
		panic("集群超过规定长度")
	}
	m := len(ports)
	if m!=n{
		panic("配置文件host和port长度不一致")
	}
	var address = []string{"","","","","","","","","",""}
	for k, _ := range hosts {
		address[k] = hosts[k]+":"+ports[k]
	}
	var cli *clientv3.Client
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   address,		//连接地址，多个
		DialTimeout: 5 * time.Second,                    //超时时间
	})
	if err != nil {
		fmt.Printf("failed to connect etcd, err: %v\n", err)
		return err
	}
	fmt.Println("etcd connect success")
	defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	return
}

//检查etcd中conf配置修改和新增
func (t *tailLogMgr)WatchConf(){
	host := viper.GetString("etcd.host")
	port := viper.GetString("etcd.port")
	hosts := strings.Split(host,",")
	ports := strings.Split(port,",")
	n := len(hosts)
	if n>10{
		panic("集群超过规定长度")
	}
	m := len(ports)
	if m!=n{
		panic("配置文件host和port长度不一致")
	}
	var address = []string{"","","","","","","","","",""}
	for k, _ := range hosts {
		address[k] = hosts[k]+":"+ports[k]
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   address,		//连接地址，多个
		DialTimeout: 5 * time.Second,                    //超时时间
	})
	if err != nil {
		fmt.Printf("failed to connect etcd, err: %v\n", err)
	}
	fmt.Println("etcd connect success")
	defer cli.Close()
	// watch
	var logEntryConf []*LogEntry
	rch := cli.Watch(context.Background(), viper.GetString("etcd.key")) // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
			err = json.Unmarshal(ev.Kv.Value, &logEntryConf)
			if err != nil {
				fmt.Printf("Unmarshal etcd value failed, err:%v\n", err)
				continue
			}
			NewConfChan()<-logEntryConf
		}
	}
}

