package objects

import (
	"NOS/nosServer/metadata"
	"net/http"
)
func delete(objectName string, isExist bool, objectInfoMap map[string]string, w http.ResponseWriter)  {
	// delete,说明要删除该object，这里我们在元数据里进行软删除
	if isExist {
		// object存在，对其进行软删除
		if metadata.DeleteObject(objectName, objectInfoMap["sha256_code"]) {
			// 说明删除成功
			WriteLog.Println("delete object", objectName, "is ok")
			w.WriteHeader(200)
			return
		}
		// 说明删除失败
		WriteLog.Println("sorry, delete object", objectName, "is bad")
		w.WriteHeader(417)
		return
	}
	// object存在，返回404
	WriteLog.Println("sorry, object", objectName, "is not exist")
	w.WriteHeader(404)
	return
}

////func delete(url string, w http.ResponseWriter){
////	// 根据URL发起delete请求
////	req, err := http.NewRequest("DELETE", url, nil)
////	if err != nil {
////		// 生成delete请求失败
////		WriteLog.Println("second operation of make delete request is bad, err is", err)
////		w.WriteHeader(406)
////		return
////	}
////	resp, err1 := http.DefaultClient.Do(req)
////	if err1 != nil {
////		// 生成delete请求失败
////		WriteLog.Println("second operation of delete request is bad, err is", err)
////		w.WriteHeader(406)
////		return
////	}
////	// 发起delete请求成功,根据返回的statusCode来处理
////	statusCode := resp.StatusCode
////	if statusCode == 200 {
////		// 说明kvserver执行delete操作成功
////		WriteLog.Println("second operation of delete is ok")
////		w.WriteHeader(200)
////		return
////	}
////	// 说明kvserver执行delete操作成功
////	WriteLog.Println("second operation of delete is bad")
////	w.WriteHeader(404)
////	return
////}
//
