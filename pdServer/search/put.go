package search

import (
	"fmt"
	"math/rand"
	"net/http"
)

// 说明httpServer发过来的是put请求，put请求的话只需要返回给client一个可用的kvserver地址即可
func Put(w http.ResponseWriter, r *http.Request)  {
	 fmt.Println("yaoxi you are put")
	 // 从KvServer中随机获取一个可用的kvserver
	 num := rand.Intn(len(KvServer))
	 kvserver := KvServer[num]
	 // 这里需要判断，返回的kvserver是否正常，如果kvserver的长度为0表示没有可用的kvserver，不为0表示返回kvserver正常
	 if len(kvserver) == 0 {
	 	// 表示没有可用的kvserver
	 	w.WriteHeader(http.StatusNotAcceptable)
		 WriteLog.Println("search kvserver is bad, all kvservers is unusable")
		 return
	 }
	 w.WriteHeader(http.StatusOK)
	 w.Write([]byte(kvserver))
	 WriteLog.Println("send kvserver", kvserver, "to client is ok")
	 return
}
