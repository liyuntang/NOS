package objects

import (
	"NOS/nosServer/etcd"
	"NOS/nosServer/metadata"
	"NOS/nosServer/tomlConfig"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// 定义几个该包的全局变量
var (
	WriteLog   *log.Logger
	EtcdServer tomlConfig.ETCD
)

// 该函数为接入层的处理函数，处理流程如下：
// 1、判断method
//		如果method==get则到metada表里查询该objectName
//			如果该objectName存在，则拿到该objectName对应data的sha256_code值，并以该值为数据存储的名称向所有kvserver进行广播查找object所在的机器
//			如果该objectName不存在，则直接返回404
//		如果method==delete则到metada表里查询该objectName
//			如果该objectName存在，则判断is_del的值
//				如果该值为0(未删除)则将其设置为1(已删除)
//				如果该值为1则直接返回404
//			如果该objectName不存在，则直接返回404
//		如果method==put则到metada表里查询该objectName
//			如果objectName存在且is_del=0则表示该object存在，返回500(文件已存在)
//			如果objectName存在且is_del=1或者objectName为空则表示该object不存在，此时需要计算其sha256的值，并以该值为objectName，存入随机选中的kvserver

// 2、广播过程如下：
// 		1）从etcd中获取kvserver信息
// 		2）获取objectName对应的sha256_code，并封装成search请求，并发方式向kvserver进行广播，
//		3）如果该kvserver存在该object的数据则将自己的IP地址放入channel中，如果没有则不做任何回应，该版本中最多有且只有一个kvserver会返回一个ok状态
// 		4）设置一个select，如果超过1s还没有从channel中取到数据，则返回400
//		5）对kvserver发起secondOperation操作
// 3、该系统只支持get、put、delete操作，其他操作据返回500报错

// 4、object存在的判断条件： object_name不为空且is_del=0

// 5、该系统加入了元数据服务，其设计如下：
//		1）元数据服务选择tidb作为存储媒介，也可以选择其他的媒介，该版本支持mysql、pg、tidb等数据库，如果选择es的话需要额外开发
//		2）元数据信息存储在db_platform.nos_metadata表中
//		3）对象名称与对象数据存储的文件名称解耦，数据存储的文件名称是对象数据的sha256加密后的字符串

func Handler(w http.ResponseWriter, r *http.Request) {
	// 判断请求方法
	method := r.Method
	// 获取objectName
	objectName := fmt.Sprintf(strings.Split(r.URL.EscapedPath(), "/")[1])
	fmt.Println("method is", method, "objectName is", objectName)
	// 根据系统设定我们只支持get、put、delete都是可允许的正常需求
	if strings.ToLower(method) == "get" || strings.ToLower(method) == "put" || strings.ToLower(method) == "delete" {
		// 说明客户端发起的请求是我们可允许的，下面对object进行判断
		// 不管是get、put还是delete都需要确认object是否存在
		// 如果存在则返回true以及一个存放了object信息的map
		// 如果不存在则返回false以及一个空map
		isok, objectInfoMap := metadata.ObjectISOK(objectName)
		if strings.ToLower(method) == "get" {
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
				// 此时我们已经拿到了kvserver地址，下面进行secondOperation操作即可
				secondOperation(sha256Code, kvserver, method, w, r)
			}
			// 说明object不存在，直接返回404
			WriteLog.Println("object", objectName, "is not exist")
			w.WriteHeader(404)
			return
		} else if strings.ToLower(method) == "put" {
			if isok {
				// 说明该object存在，此时返回
				WriteLog.Println("sorry, object", objectName, "is exist")
				w.WriteHeader(400)
				return
			}
			// 说明该object不存在，此时需要计算sha256_code值，并以此为objectName进行存储
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				// 读取客户端发送的数据失败
				WriteLog.Println("read data from client is bad")
				w.WriteHeader(400)
				return
			}
			// 读取客户端发送的数据成功，计算sha256_code值
			sha256Code := metadata.MakeSha256(buf)
			fmt.Println(sha256Code)

			// 从etcd中获取kvserver信息
			kvservers := etcd.EtcdGet(EtcdServer, WriteLog)
			// 如果method为put操作则随机抽取一台kvserver执行put操作即可
			kvserver := kvservers[rand.Intn(len(kvservers))]
			// 封装put操作
			secondOperation(sha256Code, kvserver, method, w, r)
			// 如果数据落盘成功则需要将objectName、sha256Code记录到元数据里
			if metadata.InsertObject(objectName, sha256Code) {
				// 记录元数据成功，则返回200
				WriteLog.Println("put object", objectName, "is ok")
				w.WriteHeader(200)
				return
			}
			// 记录元数据失败，失败返回417
			WriteLog.Println("put object", objectName, "is bad")
			w.WriteHeader(417)
			return
		} else {
			// delete,说明要删除该object，这里我们在元数据里进行软删除
			if isok {
				// object存在，对其进行软删除
				if metadata.DeleteObject(objectName, objectInfoMap["sha256_code"]) {
					// 说明删除成功
					WriteLog.Println("delete object", objectName, "is ok")
					w.WriteHeader(200)
					return
				}
				// 说明删除失败
				WriteLog.Println("sorry, delete object", objectName, "is bad")
				w.WriteHeader(417)
				return
			}
			// object存在，返回404
			WriteLog.Println("sorry, object", objectName, "is not exist")
			w.WriteHeader(404)
			return
		}


	} else {
		// 该系统只支持get、put、delete三个操作，其他的全部报错
		WriteLog.Println("sorry, method", method, "is not allow, we only suport get/put/delete")
		w.WriteHeader(405)
		return
	}
}

