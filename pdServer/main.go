package main

import (
	"NOS/pdServer/common"
	"NOS/pdServer/objects"
	"NOS/pdServer/tomlConfig"
	"flag"
	"fmt"
	"net/http"
	"os"
)

// 定义一些变量
var (
	config string
)

// flag

func init()  {
	flag.StringVar(&config, "c", "", "pd service config file")
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
	pdlog := common.CreateLog(configration.System.LogFile)

	// 启动一个http作为接入层，对外提供服务
	httpEndPoint := fmt.Sprintf("%s:%d", configration.System.Address, configration.System.Port)
	http.HandleFunc("/", objects.Handler)
	http.ListenAndServe(httpEndPoint, nil)

}
