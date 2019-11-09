package encapsulation

import (
	"fmt"
	"io"
	"net/http"
)

func put(url string, reader io.Reader) bool{
	// 根据URL发起put请求
	fmt.Println("start put=====================")
	req, err := http.NewRequest("PUT", url, reader)
	if err != nil {
		// 生成put请求失败
		WriteLog.Println("second operation of make put request is bad, err is", err)
		return false
	}
	resp, err1 := http.DefaultClient.Do(req)
	if err1 != nil {
		// 生成put请求失败
		WriteLog.Println("second operation of put request is bad, err is", err)
		return false
	}
	// 发起put请求成功,根据返回的statusCode来处理
	statusCode := resp.StatusCode
	if statusCode == 200 {
		// 说明kvserver执行put操作成功
		WriteLog.Println("second operation of put is ok")
		return true
	}
	// 说明kvserver执行put操作失败
	WriteLog.Println("second operation of put is bad, beacuse kvserver return data is bad")
	return false
}

