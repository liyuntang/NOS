package main

import (
	"NOS/pdServer/search"
	"NOS/pdServer/start"
	"NOS/pdServer/tomlConfig"
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
	search.WriteLog = generalog
	search.KvServer = configration.PdServer.KvServer

	// 根据配置启动一个http server
	endPoint := fmt.Sprintf("%s:%d", configration.PdServer.Address, configration.PdServer.Port)
	generalog.Println("start endPoint", endPoint)
	start.StartPD(endPoint, generalog)
}
