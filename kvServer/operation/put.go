package operation

import (
	"io"
	"net/http"
	"os"
)

func put(w http.ResponseWriter, r *http.Request, objectPath string)  {
	// 查看文件是否存在
	_, err := os.Stat(objectPath)
	if os.IsNotExist(err) {
		// 说明文件不存在,存入文件即可
		file, err := os.OpenFile(objectPath, os.O_CREATE|os.O_WRONLY, 0644)
		defer file.Close()
		if err != nil {
			WriteLog.Println("open file", objectPath, "is bad")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err1 := io.Copy(file, r.Body)
		if err1 != nil {
			// 说明写入失败，
			w.WriteHeader(http.StatusExpectationFailed)
			WriteLog.Println("write to file", objectPath, "is bad")
			return
		}
		// 说明写入成功
		w.WriteHeader(http.StatusOK)
		WriteLog.Println("write file", objectPath, "is ok")
		return
	}
	// 说明文件存在，提示w文件已存在
	w.WriteHeader(http.StatusInternalServerError)
	WriteLog.Println("sorry object", objectPath, "is exist")
	return
}
