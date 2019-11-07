package etcd

import (
	"NOS/nosServer/tomlConfig"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"log"
	"time"
)

func EtcdGet(config tomlConfig.ETCD, logg *log.Logger) []string {
	// 声明一个存放返回信息的slice
	getSlice := []string{}
	// 连接etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:config.EtcdServers,
		DialTimeout:time.Duration(config.EtcdTimeOut)*time.Second,
	})
	if err != nil{
		// 说明连接etcd报错
		logg.Println("connect to etcd", config.EtcdServers, "is bad")
		return nil
	}
	// 说明连接etcd成功
	logg.Println("connect to etcd", config.EtcdServers, "is ok")
	// 获取ctx
	ctx, _ := context.WithCancel(context.TODO())
	// 执行get操作
	resp, err1 := cli.Get(ctx, config.EtcdDir, clientv3.WithPrefix())
	if err1 != nil {
		logg.Println("get from etcd", config.EtcdServers, "is bad")
		return nil
	}
	for _, kvs := range resp.Kvs{
		getSlice = append(getSlice, string(kvs.Value))
	}
	return getSlice
}
