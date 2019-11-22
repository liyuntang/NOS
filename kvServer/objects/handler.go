package objects

import (
	"log"
	"net/http"
)

// 定义几个该包的全局变量
var (
	DataDir string
	WriteLog *log.Logger
)

// 由于我们有元数据服务，所有的delete操作均在元数据里做软删除，所以数据服务层只支持get、put操作
func Handler(w http.ResponseWriter, r *http.Request)  {
	// 判断请求方法
	method := r.Method

	if method == http.MethodGet {
		// get
		get(w, r)
	} else {
		// put
		put(w, r)
	}

}
