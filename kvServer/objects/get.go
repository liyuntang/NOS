package objects

import (
	"io"
	"net/http"
	"os"
)

func get(w http.ResponseWriter, objectPath string)  {
	// 查看文件是否存在
	_, err := os.Stat(objectPath)
	if os.IsNotExist(err) {
		// 说明文件不存在
		WriteLog.Println("file", objectPath, "is not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 说明文件存在，此时读取数据返回给w即可
	file, err1 := os.Open(objectPath)
	defer file.Close()
	if err1 != nil {
		// 说明打开object失败
		WriteLog.Println("get file", objectPath, "is bad")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 说明打开object成功，返回数据
	io.Copy(w, file)
	w.WriteHeader(200)
	return
}
