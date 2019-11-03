package operation

import (
	"net/http"
	"os"
)

func delete(w http.ResponseWriter, objectPath string)  {
	// 查看文件是否存在
	_, err := os.Stat(objectPath)
	if os.IsNotExist(err) {
		// 说明文件不存在,提示文件不存在即可
		w.WriteHeader(http.StatusNotFound)
		WriteLog.Println("file", objectPath, "is not exist")
		return
	}
	// 说明文件存在，删除该文件
	if err1 := os.Remove(objectPath); err1 != nil {
		// 说明删除文件失败
		w.WriteHeader(http.StatusBadRequest)
		WriteLog.Println("delete file", objectPath, "is bad")
		return
	}
	// 说明删除文件成功
	w.WriteHeader(http.StatusOK)
	WriteLog.Println("delete file", objectPath, "is ok")
	return
}
