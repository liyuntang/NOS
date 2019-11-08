package objects

import (
	"NOS/nosServer/etcd"
	"NOS/nosServer/tomlConfig"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// 定义几个该包的全局变量
var (
	WriteLog *log.Logger
	EtcdServer tomlConfig.ETCD
)

// 该函数为接入层的处理函数，处理流程如下：
// 1、判断method，如果method为get、delete则需要对所有kvserver进行广播查找object所在的机器，如果为put则随机抽取一台kvserver执行put操作即可
// 2、从etcd中获取kvserver信息
// 3、获取objectName，并封装成search请求
// 4、如果method为get、delete则对所有kvserver进行广播，如果存在该object，则返回ok状态，该版本中最多有且只有一个kvserver会返回一个ok状态
//	 如果返回的状态为404，则提示用户404
// 5、如果method为put操作则随机抽取一台kvserver执行put操作即可

func Handler(w http.ResponseWriter, r *http.Request)  {
	// 判断请求方法
	method := r.Method
	// get、put、delete都是可允许的正常需求，由于他们有三个共同的操作所以他们放在了一块，三个共同的操作为：
	// 1、从etcd中获取kvserver信息    2、获取objectName		3、封装请求发送到kvserver(get、delete为search请求，put为put请求)
	if method == http.MethodGet || method == http.MethodDelete || method == http.MethodPut{
		// 声明一个变量用于标示kvserver地址
		var kvserver string

		// 从etcd中获取kvserver信息
		kvservers := etcd.EtcdGet(EtcdServer, WriteLog)
		if len(kvservers) == 0 {
			// 说明从etcd获取kv地址失败，
			WriteLog.Println("get kvserver from etcd is bad")
			w.WriteHeader(400)
			return
		}
		// 获取objectName
		objectName := fmt.Sprintf(strings.Split(r.URL.EscapedPath(),"/")[1])

		if method == http.MethodGet || method == http.MethodDelete {
			// 如果method为get、delete则需要对所有kvserver进行广播查找object所在的机器
			// 定义一个channel，如果该kvserver存在objectName，则将kvserver放入channel中，开启一个select如果从channel中读取到了数据则连接该kvserver
			ch := make(chan string, 1)
			// 发起新请求，如果超过1s中还没有从channel中获取到数据则超时，返回客户端信息
			for _, kvIp := range kvservers {
				go broadcast(objectName, kvIp, ch)
			}
			// 从channel中读取数据，如果超过1s还没有数据则超时
			select {
			case kvAddress := <- ch:
				kvserver = kvAddress
			case <- time.After(1*time.Second):
				// 说明kv返回超时，或者该objectName不存在
				WriteLog.Println("object", objectName, "is not exist")
				w.WriteHeader(404)
				return
			}
		} else {
			// 如果method为put操作则随机抽取一台kvserver执行put操作即可
			kvserver = kvservers[rand.Intn(len(kvservers))]
		}

		// 到此为止，我们已经拿到kvserver的地址，这里我们开始封装用户的请求，这个操作适用于get、put、delete请求
		secondOperation(objectName, kvserver, method, w, r)
	} else {
		// 该系统只支持get、put、delete三个操作，其他的全部报错
		WriteLog.Println("sorry, method", method, "is not allow, we only suport get/put/delete")
		w.WriteHeader(405)
		return
	}

}
