package etcd

import (
	"NOS/kvServer/tomlConfig"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"log"
	"time"
)

func Put(config tomlConfig.ETCD, endPoint string, logger *log.Logger)  {
	// 连接etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:config.EtcdServers,
		DialTimeout:time.Duration(config.EtcdTimeOut)*time.Second,
	})
	if err != nil{
		// 说明连接etcd报错
		logger.Println("connect to etcd", config.EtcdServers, "is bad")
		return
	}
	// 说明连接etcd成功
	//logger.Println("connect to etcd", config.EtcdServers, "is ok")
	// 获取ctx
	ctx, _ := context.WithCancel(context.TODO())

	// 生成租约
	lease := clientv3.NewLease(cli)
	lgr, err := lease.Grant(context.TODO(), config.Lease)
	if err != nil {
		logger.Println("set lease is bad, err is", err)
	}
	_, err1 := cli.Put(ctx, config.EtcdDir, endPoint, clientv3.WithLease(lgr.ID))
	if err1 != nil {
		logger.Println("sorry, put is bad, err is", err1)
	}
	//logger.Println("put is ok", )
}
