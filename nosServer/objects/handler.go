package objects

import (
	"NOS/nosServer/metadata"
	"NOS/nosServer/tomlConfig"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 定义几个该包的全局变量
var (
	WriteLog   *log.Logger
	EtcdServer tomlConfig.ETCD
	TmpDir string
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
//		如果method==put则首先判读put请求是否符合要求，在符合要求的情况下到metada表里查询该objectName
//			如果objectName存在且is_del=0则表示该object存在，返回400(文件已存在)
//			如果objectName存在且is_del=1或者objectName为空则表示该object不存在，此时需要如下操作：
// 			1、检查客户端发起的put操作是否符合要求，要求如下：head中必须要有对象的大小、对象的加密方式（该版本只支持sha256），对象的sha256值，否则直接返回报错
//			2、如果客户端发起的put操作是否符合要求则将对象临时存放在nos组件的tmp目录下
// 			3、在本地计算其sha256的值，并与客户端发送过来的sha256值进行对不，如果不相同，则返回报错
// 			4、如果sha256值相同则以该值为objectName，存入随机选中的kvserver

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
	// 根据系统设定我们只支持get、put、delete都是可允许的正常需求
	if strings.ToLower(method) == "get" || strings.ToLower(method) == "put" || strings.ToLower(method) == "delete" {
		// 获取objectName
		objectName := fmt.Sprintf(strings.Split(r.URL.EscapedPath(), "/")[1])
		// 说明客户端发起的请求是我们可允许的，下面对object进行判断
		// 不管是get、put还是delete都需要确认object是否存在
		// 如果存在则返回true以及一个存放了object信息的map
		// 如果不存在则返回false以及一个空map
		isok, objectInfoMap := metadata.ObjectISOK(objectName)
		// 下面针对object在不同method情况下进行处理
		if strings.ToLower(method) == "get" {
			// get
			get(objectName, isok, objectInfoMap, w, r)
		} else if strings.ToLower(method) == "put" {
			// 说明文件不存在，进行put操作
			put(objectName, objectInfoMap, w, r)
		} else {
			// delete
			delete(objectName, isok, objectInfoMap, w)
		}
	} else {
		// 该系统只支持get、put、delete三个操作，其他的全部报错
		WriteLog.Println("sorry, method", method, "is not allow, we only suport get/put/delete")
		w.WriteHeader(405)
		return
	}
}

