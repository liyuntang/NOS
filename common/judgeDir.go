package common

import (
	"io/ioutil"
	"log"
	"os"
)

func JudgeDir(dir string, logger *log.Logger)  {
	// 判断dir是否存在
	fileInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// 说明dir不存在
		logger.Println("dir", dir, "is not exist,and now we create it")
		createDir(dir, logger)
		return
	}
	// 说明dir存在，现在需要判读dir是否为目录，如果是目录的话是否为空
	if fileInfo.IsDir() {
		// 说明是目录，下面判断该目录是否为空,如果不为空则需要提示用户
		res, err := ioutil.ReadDir(dir)
		if err != nil {
			logger.Println("list dir", dir, "is bad, err is", err)
			os.Exit(0)
		}
		if len(res) != 0 {
			logger.Println("[warning] dir", dir, "is not null")
		}
	}
	// 说明不是目录，需要创建该dir
	createDir(dir, logger)
}

func createDir(dir string, logger *log.Logger)  {
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Println("create dir", dir, "is bad, err is", err)
		os.Exit(0)
	}
	logger.Println("create dir", dir, "is ok")
}