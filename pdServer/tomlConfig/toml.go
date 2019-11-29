package tomlConfig

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"sync"
)

// 声明几个变量
var (
	conf *PDS
	once sync.Once
)

func TomlConfig(configFile string) *PDS {
	// 检测configFile的状态，并获取configFile的绝对路径
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		fmt.Println("configFile", configFile, "is not exist")
		os.Exit(0)
	}
	// 说明configFile存在，获取其绝对路径信息
	absPath, err := filepath.Abs(configFile)
	if err != nil {
		fmt.Println("get abs path of configFile", configFile, "is bad")
		os.Exit(0)
	}

	// 单例模式解析配置
	once.Do(func() {
		_, err := toml.DecodeFile(absPath, &conf)
		if err != nil {
			fmt.Println("toml configFile", configFile, "is bad, err is", err)
			os.Exit(0)
		}

	})
	return conf
}
