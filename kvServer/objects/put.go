package objects

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func put(w http.ResponseWriter, r *http.Request, objectPath string)  {
	defer r.Body.Close()
	gFile := fmt.Sprintf("%s.gz", objectPath)
	_, err := os.Stat(gFile)
	if os.IsNotExist(err) {
		// 说明文件不存在,存入文件即可
		file, err := os.OpenFile(gFile, os.O_CREATE|os.O_WRONLY, 0644)
		defer file.Close()
		if err != nil {
			WriteLog.Println("open file", gFile, "is bad")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// 将数据压缩存入文件
		buf, _ := ioutil.ReadAll(r.Body)
		if gzipWrite(file, buf) {
			// 说明写入数据成功
			w.WriteHeader(http.StatusOK)
			WriteLog.Println("write file", gFile, "is ok")
			return
		}
		// 说明写入数据失败
			w.WriteHeader(http.StatusExpectationFailed)
			WriteLog.Println("write to file", gFile, "is bad")
			return
		//_, err1 := io.Copy(file, r.Body)
		//if err1 != nil {
		//	// 说明写入失败，
		//	w.WriteHeader(http.StatusExpectationFailed)
		//	WriteLog.Println("write to file", objectPath, "is bad")
		//	return
		//}
		//// 说明写入成功
		//w.WriteHeader(http.StatusOK)
		//WriteLog.Println("write file", objectPath, "is ok")
		//return
	}
	// 说明文件存在，提示w文件已存在
	w.WriteHeader(http.StatusInternalServerError)
	WriteLog.Println("sorry object", objectPath, "is exist")
	return
}

func gzipWrite(file *os.File, data []byte) bool {
	// 初始化write
	gWrite := gzip.NewWriter(file)
	// 写入数据
	_, err1 := gWrite.Write(data)
	if err1 != nil {
		return false
	}
	return true
}
