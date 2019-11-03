package main

import (
	"NOS/common"
	"NOS/httpServer"
	"NOS/objects"
	"NOS/tomlConfig"
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
	objects.DataDir = configration.System.DataDir
	objects.WriteLog = generalog

	// 判断dataDir是否存在，如果不存在则创建该目录，如果存在则判断是否为空，如果不为空应当提示用户
	common.JudgeDir(configration.System.DataDir, generalog)

	// 根据配置启动一个http server
	endPoint := fmt.Sprintf("%s:%d", configration.System.Address, configration.System.Port)
	httpServer.OpenServer(endPoint, generalog)
}



