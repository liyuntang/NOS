package operation

import (
	"net/http"
	"os"
	"strings"
)
var (
	KVSERVER string
)
func PDSearch(w http.ResponseWriter, r *http.Request)  {
	// 拼接文件名称
	filePath := DataDir+"/"+strings.Split(r.URL.EscapedPath(),"/")[2]
	// 查看文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 说明文件不存在
		WriteLog.Println("file", filePath, "is not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 说明文件存在，此时返回该kvserver的endpoint
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(KVSERVER))
	return
}
