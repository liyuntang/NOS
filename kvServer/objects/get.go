package objects

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request)  {
	// 获取object_name
	gFile := fmt.Sprintf("%s/%s.gz", DataDir, strings.Split(r.URL.EscapedPath(),"/")[1])
	// 查看文件是否存在
	_, err := os.Stat(gFile)
	if os.IsNotExist(err) {
		// 说明文件不存在
		WriteLog.Println("file", gFile, "is not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 说明文件存在，此时读取数据返回给w即可
	file, err1 := os.Open(gFile)
	defer file.Close()
	if err1 != nil {
		// 说明打开object失败
		WriteLog.Println("get file", gFile, "is bad")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 初始化一个reader
	gReader, err2 := gzip.NewReader(file)
	defer gReader.Close()
	if err2 != nil {
		WriteLog.Println("new gizp reader of", gFile, "is bad")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 说明打开object成功，返回数据
	io.Copy(w, gReader)
	w.WriteHeader(200)
	return
}
