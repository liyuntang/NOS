package objects

import (
	"net/http"
	"os"
	"strings"
)

func delete(w http.ResponseWriter, r *http.Request)  {
	// 拼接文件名称
	filePath := DataDir+"/"+strings.Split(r.URL.EscapedPath(),"/")[1]
	// 查看文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 说明文件不存在,提示文件不存在即可
		w.WriteHeader(http.StatusNotFound)
		WriteLog.Println("file", filePath, "is not exist")
		return
	}
	// 说明文件存在，删除该文件
	if err1 := os.Remove(filePath); err1 != nil {
		// 说明删除文件失败
		w.WriteHeader(http.StatusBadRequest)
		WriteLog.Println("delete file", filePath, "is bad")
		return
	}
	// 说明删除文件成功
	w.WriteHeader(http.StatusOK)
	WriteLog.Println("delete file", filePath, "is ok")
	return
}
