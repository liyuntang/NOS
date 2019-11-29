package objects

import (
	"NOS/nosServer/metadata"
	"fmt"
	"net/http"
	"strings"
)
func delete(w http.ResponseWriter, r *http.Request) {
	// 获取对象名称
	objectName := fmt.Sprintf(strings.Split(r.URL.EscapedPath(), "/")[1])
	// 确认object是否存在如果存在则返回true以及一个存放了object信息的map,如果不存在则返回false以及一个空map
	isok, objectInfoMap := metadata.ObjectISOK(objectName)

	// delete,说明要删除该object，这里我们在元数据里进行软删除
	if isok {
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
	// object不存在，返回404
	WriteLog.Println("sorry, object", objectName, "is not exist")
	w.WriteHeader(404)
	return
}
