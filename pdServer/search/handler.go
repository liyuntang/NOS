package search

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// 加载一些变量
var (
	WriteLog *log.Logger
	KvServer []string
)

// handler的处理过程是这样的：
// 1、获取object_name（默认所有的dataDir全部相同，其实不相同也可以）
// 2、通过pdServer向所有的kvServer发布广播（默认所有的kvServer都正常提供服务），查询该object所在的主机，如果该主机存在该object则返回一个消息，不存在则静默
// 3、拿到pdServer返回的主机endpoint，连接到该endpoint进行相应的操作（get、put、delete）
// 4、返回执行结果

func Seacher(w http.ResponseWriter, r *http.Request) {
	// search操作一定是get操作,如果是其他的操作，则抱回报错
	// 这个地方开启多线程广播
	method := r.Method
	if strings.ToLower(method) == "get" {

		// 获取objectName
		objectName := strings.Split(r.URL.EscapedPath(), "/")[2]
		// 这个地方应该用协程并发操作的，但是我们这里先单线程运行，
		for _, kvserver := range KvServer {
			fmt.Println("我要开始广播了", method, kvserver)
			// 拼接object请求kvServer里的全路径,这里我们对pd的定位是介于httpServer与kvServer中间的一个桥梁，所以pd发起的查询要有一个标志
			kvName := fmt.Sprintf("http://%s/pd/%s", kvserver, objectName)
			fmt.Println("开始到kvServer里get请求", kvName)
			reps, err := http.Get(kvName)
			if err != nil {
				// 说明请求kvserver报错，返回服务不可达
				w.WriteHeader(http.StatusNotAcceptable)
				WriteLog.Println("sorry, kvserver", kvserver, "is not acceptable")
				return
			}
			fmt.Println("kvserver的返回code是", reps.StatusCode)
			fmt.Println("kvserver的返回header是", reps.Header)
			fmt.Println("kvserver的返回body是", reps.Body)
			fmt.Println("kvserver的返回status是", reps.Status)

			// 说明请求kvserver正常，此时需要根据返回的statuscode进行判断，如果为200则表示正常，将kvserver地址返回httpServer，否则返回报错
			if reps.StatusCode == 200 {
				// 说明获取到了object存放的kvserver，返回该kvserver地址
				kvMessage, err := ioutil.ReadAll(reps.Body)
				reps.Body.Close()
				if err != nil {
					WriteLog.Println("get kvMessage from reps is bad")
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
				w.WriteHeader(http.StatusOK)
				//fmt.Println("开始给httpServer返回kvServer地址", kvMessage, string(kvMessage))
				w.Write(kvMessage)
				return

			} else if reps.StatusCode == 404 {
				w.WriteHeader(404)
				return
			}else {
				// 说明发起的不是get操作，返回报错
				w.WriteHeader(http.StatusMethodNotAllowed)
				WriteLog.Println("method is not allowd")
				return
			}
		}
	}
}