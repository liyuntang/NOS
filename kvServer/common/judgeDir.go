package common

import (
	"log"
	"os"
)

func JudgeDir(dataDir string, logger *log.Logger) {
	// 判读数据目录是否存在，如果不存在则创建
	_, err1 := os.Stat(dataDir)
	if os.IsNotExist(err1) {
		logger.Println("dataDir", dataDir, "is not exist,and now we create it")
		err := os.MkdirAll(dataDir, 0755)
		if err != nil {
			// 说明创建目录失败
			logger.Println("create dataDir", dataDir, "is bad")
			os.Exit(0)
		}
	}
}
