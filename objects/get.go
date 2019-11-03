package objects

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request)  {
	// 拼接文件名称
	filePath := DataDir+"/"+strings.Split(r.URL.EscapedPath(),"/")[1]
	// 查看文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 说明文件不存在
		WriteLog.Println("file", filePath, "is not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 说明文件存在，此时读取数据返回给w即可
	file, err1 := os.Open(filePath)
	defer file.Close()
	if err1 != nil {
		WriteLog.Println("get file", filePath, "is bad")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, file)
}
