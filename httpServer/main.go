package main

import (
	"NOS/httpServer/operation"
	"NOS/httpServer/start"
	"NOS/httpServer/tomlConfig"
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
	operation.PdServer = configration.HttpServer.PdServer

	// 根据配置启动一个http server
	endPoint := fmt.Sprintf("%s:%d", configration.HttpServer.Address, configration.HttpServer.Port)
	start.OpenServer(endPoint, generalog)
}



