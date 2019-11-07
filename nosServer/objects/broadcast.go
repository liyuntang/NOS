package objects

import (
	"fmt"
	"log"
	"net/http"
)

func broadcast(objectName, kvserver string, ch chan string, logger *log.Logger) {
	// 重新封装请求
	url := fmt.Sprintf("http://%s/search/%s", kvserver, objectName)
	// 用search对kvserver发起get请求，并接收statusCode
	resp, err := http.Get(url)
	if err != nil {
		// 说明访问kvserver报错
		logger.Println("connect to kvserver", kvserver, "is bad, err is", err)
	}
	// 说明访问kvserver成功
	// 根据statusCode判读是否将kvserver写入channel
	if resp.StatusCode == 200 {
		// 说明该objectName存在该kvserver中，将kvserver放入到channel里
		ch <- kvserver
	}

}
