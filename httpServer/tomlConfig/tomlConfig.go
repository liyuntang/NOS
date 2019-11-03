package tomlConfig

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"sync"
)

var (
	conf *HTTPSERVER
	once sync.Once
)

func TomlConfig(configFile string) *HTTPSERVER {
	// 获取配置文件的绝对路径
	absPath, err := filepath.Abs(configFile)
	if err != nil {
		fmt.Println("get abs of configFile", configFile, "is bad")
		os.Exit(0)
	}

	once.Do(func() {
		_, err := toml.DecodeFile(absPath, &conf)
		if err != nil {
			fmt.Println("toml decode file", absPath, "is bad")
			os.Exit(0)
		}
	})
	return conf
}