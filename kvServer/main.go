package main

import (
	"NOS/kvServer/operation"
	"NOS/kvServer/start"
	"NOS/kvServer/tomlConfig"
	"NOS/kvServer/common"
	"flag"
	"fmt"
	"os"
)

var (
	config string
)

func init()  {
	flag.StringVar(&config, "c", "", "configration file")
}


func main()  {
	flag.Parse()
	if len(os.Args) <= 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}
	// 解析配置
	configration := tomlConfig.TomlConfig(config)

	// 初始化日志
	generalog := tomlConfig.InitLogger()

	// 为http初始化一些变量
	operation.WriteLog = generalog
	operation.DataDir = configration.KvServer.DataDir

	// 判断dataDir是否存在，如果不存在则创建该目录，如果存在则判断是否为空，如果不为空应当提示用户
	common.JudgeDir(configration.KvServer.DataDir, generalog)

	// 传递一个参数
	operation.KVSERVER = fmt.Sprintf("%s:%d", configration.KvServer.Address, configration.KvServer.Port)
	// 根据配置启动一个http server
	fmt.Println("start endPoint is", operation.KVSERVER)
	start.OpenServer(operation.KVSERVER, generalog)
}



