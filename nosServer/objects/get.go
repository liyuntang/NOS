package objects

import (
	"NOS/nosServer/encapsulation"
	"NOS/nosServer/etcd"
	"NOS/nosServer/metadata"
	"fmt"
	"net/http"
	"strings"
	"time"
)
// 6228 4118 2302 8739 463
func get(w http.ResponseWriter, r *http.Request) {
	// 获取对象名称
	objectName := fmt.Sprintf(strings.Split(r.URL.EscapedPath(), "/")[1])
	// 确认object是否存在如果存在则返回true以及一个存放了object信息的map,如果不存在则返回false以及一个空map
	isok, objectInfoMap := metadata.ObjectISOK(objectName)
	if isok {
		// 说明object存在,此时需要做两个操作：
		// 1、从objectInfoMap中获取sha256_code值，并以此为objectName到kvserver中广播
		// 2、从etcd中获取kvserver信息，然后对所有kvserver进行广播
		// 获取sha256_code值
		sha256Code := objectInfoMap["sha256_code"]
		// 声明一个变量用于标示kvserver地址
		var kvserver string
		// 从etcd中获取kvserver信息
		kvservers := etcd.EtcdGet(EtcdServer, WriteLog)
		if len(kvservers) == 0 {
			// 说明从etcd获取kv地址失败，返回400
			WriteLog.Println("get kvserver from etcd is bad")
			w.WriteHeader(400)
			return
		}
		// 说明从etcd获取kv地址成功，对其进行广播
		// 定义一个channel，如果该kvserver存在objectName，则将kvserver放入channel中，开启一个select如果从channel中读取到了数据则连接该kvserver
		ch := make(chan string, 1)
		// 发起新请求，如果超过1s中还没有从channel中获取到数据则超时，返回客户端信息
		for _, kvIp := range kvservers {
			go broadcast(sha256Code, kvIp, ch)
		}
		// 从channel中读取数据，如果超过1s还没有数据则超时
		select {
		case kvAddress := <-ch:
			kvserver = kvAddress
		case <-time.After(1 * time.Second):
			// 说明kv返回超时，或者该objectName不存在
			WriteLog.Println("broadcast object", objectName, "is time out")
			w.WriteHeader(408)
			return
		}
		// 此时我们已经拿到了kvserver地址，下面对请求进行封装操作
		encapsulation.SecondOperationGet(sha256Code, kvserver, "get", w, r)
		return
	}
	// 说明object不存在，直接返回404
	WriteLog.Println("object", objectName, "is not exist")
	w.WriteHeader(404)
	return
}
