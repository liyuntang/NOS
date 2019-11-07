package objects

import (
	"net/http"
)

func delete(url string, w http.ResponseWriter){
	// 根据URL发起delete请求
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		// 生成delete请求失败
		WriteLog.Println("second operation of make delete request is bad, err is", err)
		w.WriteHeader(406)
		return
	}
	resp, err1 := http.DefaultClient.Do(req)
	if err1 != nil {
		// 生成delete请求失败
		WriteLog.Println("second operation of delete request is bad, err is", err)
		w.WriteHeader(406)
		return
	}
	// 发起delete请求成功,根据返回的statusCode来处理
	statusCode := resp.StatusCode
	if statusCode == 200 {
		// 说明kvserver执行delete操作成功
		WriteLog.Println("second operation of delete is ok")
		w.WriteHeader(200)
		return
	}
	// 说明kvserver执行delete操作成功
	WriteLog.Println("second operation of delete is bad")
	w.WriteHeader(404)
	return
}

