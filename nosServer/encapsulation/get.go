package encapsulation

import (
	"io/ioutil"
	"net/http"
)

func get(url string, w http.ResponseWriter, r *http.Request){
	// 根据URL发起get请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// 生成get请求失败
		WriteLog.Println("second operation of make get request is bad, err is", err)
		w.WriteHeader(406)
		return
	}
	resp, err1 := http.DefaultClient.Do(req)
	if err1 != nil {
		// 生成get请求失败
		WriteLog.Println("second operation of get request is bad, err is", err)
		w.WriteHeader(406)
		return
	}
	// 发起get请求成功,根据返回的statusCode来处理
	statusCode := resp.StatusCode
	if statusCode == 200 {
		// 说明kvserver返回数据,读取返回的数据
		buf, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			// 说明读取数据失败
			WriteLog.Println("second operation of get is bad, beacuse kvserver return data is ok,but read it is bad, err is", err2)
			w.WriteHeader(404)
			return
		}
		// 说明读取数据成功，返回数据
		WriteLog.Println("second operation of get is ok")
		w.WriteHeader(200)
		w.Write(buf)
		return
	}
	// 说明kvserver没有返回数据
	WriteLog.Println("second operation of get is bad, beacuse kvserver return data is bad")
	w.WriteHeader(404)
	return
}
