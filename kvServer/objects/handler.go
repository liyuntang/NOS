package objects

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 定义几个该包的全局变量
var (
	DataDir string
	WriteLog *log.Logger
)


func Handler(w http.ResponseWriter, r *http.Request)  {
	// 判断请求方法
	method := r.Method

	// 获取object_name
	objectPath := fmt.Sprintf("%s/%s", DataDir, strings.Split(r.URL.EscapedPath(),"/")[1])
	if method == http.MethodGet {
		// get
		get(w, objectPath)
	} else if method == http.MethodPut {
		// put
		put(w, r, objectPath)
	} else {
		// delete
		delete(w, objectPath)
	}

}
