package objects

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request)  {
	// 获取object_name
	gFile := fmt.Sprintf("%s/%s.gz", DataDir, strings.Split(r.URL.EscapedPath(),"/")[1])
	fmt.Println("gFile is", gFile)

	// 将数据转换成[]byte类型
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteLog.Println("read data is bad, err is", err)
		w.WriteHeader(500)
		return
	}
	// 将数据写入gFile
	if writeFile(gFile, buf) {
		WriteLog.Println("write data to gFile", gFile, "is ok")
		w.WriteHeader(200)
		return
	}
	WriteLog.Println("write data to gFile", gFile, "is bad")
	w.WriteHeader(500)
	return
}

func writeFile(gFile string, data []byte) bool {
	// 打开gFile
	file, err := os.Create(gFile)
	defer file.Close()
	if err != nil {
		return false
	}
	// 生成一个writer
	gWriter := gzip.NewWriter(file)
	defer gWriter.Close()
	// 写入数据
	n, err1 := gWriter.Write(data)
	fmt.Println("n is", n)
	if err1 != nil {
		return false
	}
	return true
}
