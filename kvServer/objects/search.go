package objects

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Search(w http.ResponseWriter, r *http.Request)  {
	// 获取objectName
	objectName := fmt.Sprintf(strings.Split(r.URL.EscapedPath(),"/")[2])
	// 拼接全路径
	gFile := fmt.Sprintf("%s.gz", objectName)
	objectPath := fmt.Sprintf("%s/%s", DataDir, gFile)
	// 查看objectName状态
	// 查看文件是否存在
	_, err := os.Stat(objectPath)
	if os.IsNotExist(err) {
		// 说明文件不存在,提示文件不存在即可
		WriteLog.Println("file", objectPath, "is not exist")
		w.WriteHeader(404)
	}
	// 说明文件存在
	//WriteLog.Println("search file", objectPath, "is ok")
	w.WriteHeader(200)
}
