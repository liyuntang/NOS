package operation

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 加载一些变量
var (
	WriteLog *log.Logger
	DataDir string
)

// handler的处理过程是这样的：
// 1、获取object_name（默认所有的dataDir全部相同，其实不相同也可以）
// 2、通过pdServer向所有的kvServer发布广播（默认所有的kvServer都正常提供服务），查询该object所在的主机，如果该主机存在该object则返回一个消息，不存在则静默
// 3、拿到pdServer返回的主机endpoint，连接到该endpoint进行相应的操作（get、put、delete）
// 4、返回执行结果

func Handler(w http.ResponseWriter, r *http.Request)  {
	// 获取object_name
	objectPath := fmt.Sprintf("%s/%s", DataDir, strings.Split(r.URL.EscapedPath(),"/")[1])

	method := r.Method
	if strings.ToLower(method) == "get" {
		// 要进行get操作
		get(w, objectPath)
	} else if strings.ToLower(method) == "put" {
		// 要进行put操作
		put(w, r, objectPath)
	} else if strings.ToLower(method) == "delete" {
		// 要进行delete操作
		delete(w, objectPath)
	}else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
