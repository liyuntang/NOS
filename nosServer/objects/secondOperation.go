package objects

import (
	"fmt"
	"net/http"
	"strings"
)

func secondOperation(objectName, kvserver, method string, w http.ResponseWriter, r *http.Request)  {
	// 封装url
	url := fmt.Sprintf("http://%s/%s", kvserver, objectName)
	// 根据不通的method发起请求
	if strings.ToLower(method) == "get" {
		// get
		get(url, w)
	} else if strings.ToLower(method) == "put" {
		// put
		put(url, w, r)
	} else {
		// delete
		delete(url, w)
	}

}
