package common

import (
	"fmt"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

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

func GetEtcdClient() *clientv3.Client{
	return EtcdClient
}
