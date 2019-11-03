package start

import (
	"NOS/kvServer/operation"
	"log"
	"net/http"
)

// kvServer接收到的请求有两种：1、pdserver发起的object定位请求      2、httpServer发起的get、put、delete等操作请求

func OpenServer(endPoint string, logger *log.Logger)  {
	http.HandleFunc("/pd/", operation.PDSearch)
	http.HandleFunc("/", operation.Handler)
	http.ListenAndServe(endPoint, nil)
}
