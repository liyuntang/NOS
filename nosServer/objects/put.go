package objects

import (
	"net/http"
)

func put(url string, w http.ResponseWriter, r *http.Request){
	// 根据URL发起put请求
	req, err := http.NewRequest("PUT", url, r.Body)
	if err != nil {
		// 生成put请求失败
		WriteLog.Println("second operation of make put request is bad, err is", err)
		w.WriteHeader(406)
		return
	}
	resp, err1 := http.DefaultClient.Do(req)
	if err1 != nil {
		// 生成put请求失败
		WriteLog.Println("second operation of put request is bad, err is", err)
		w.WriteHeader(406)
		return
	}
	// 发起put请求成功,根据返回的statusCode来处理
	statusCode := resp.StatusCode
	if statusCode == 200 {
		// 说明kvserver执行put操作成功
		WriteLog.Println("second operation of put is ok")
		w.WriteHeader(200)
		return
	}
	// 说明kvserver执行put操作失败
	WriteLog.Println("second operation of put is bad, beacuse kvserver return data is bad")
	w.WriteHeader(400)
	return
}

