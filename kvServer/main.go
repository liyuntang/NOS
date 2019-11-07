package main

import (
	"NOS/kvServer/common"
	"NOS/kvServer/etcd"
	"NOS/kvServer/objects"
	"NOS/kvServer/tomlConfig"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

// 定义一些变量
var (
	config string
)

// flag

func init()  {
	flag.StringVar(&config, "c", "", "nos service config file")
}


func main()  {
	flag.Parse()
	argsLen := len(os.Args)
	if argsLen <= 1 {
		flag.PrintDefaults()
	}

	// 解析配置文件信息
	configration := tomlConfig.TomlConfig(config)

	// 初始化日志
	kvlog := common.CreateLog(configration.System.LogFile)

	// 初始化数据存储目录
	common.JudgeDir(configration.System.DataDir, kvlog)

	// 赋值几个变量
	objects.WriteLog = kvlog
	objects.DataDir = configration.System.DataDir

	// 持续向etcd服务注册自己
	endPoint := fmt.Sprintf("%s:%d", configration.System.Address, configration.System.Port)
	go func() {
		for {
			// 这个地方要注意注册的频率，一定要小于lease
			etcd.Put(configration.Etcd, endPoint, kvlog)
			time.Sleep(1*time.Second)
		}
	}()


	etcd.Put(configration.Etcd, endPoint, kvlog)

	// 启动一个http作为数据服务层，该数据服务层监听两个路径"/search/"、"/"
	// "/search/"表示发起的是广播请求
	// "/"表示发起的是正常的get、delete、put操作
	httpEndPoint := fmt.Sprintf("%s:%d", configration.System.Address, configration.System.Port)
	http.HandleFunc("/search/", objects.Search)
	http.HandleFunc("/", objects.Handler)
	http.ListenAndServe(httpEndPoint, nil)

}
