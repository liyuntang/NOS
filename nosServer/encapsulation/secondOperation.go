package encapsulation

import (
	"fmt"
	"log"
	"net/http"
)

var (
	WriteLog   *log.Logger
)

func SecondOperationGet(objectName, kvserver, method string, w http.ResponseWriter, r *http.Request)  {
	// 封装url
	url := fmt.Sprintf("http://%s/%s", kvserver, objectName)
	get(url, w, r)

}

//
func SecondOperationPut(objectName, kvserver, tmpFile string) (isok bool) {
	// 封装url
	url := fmt.Sprintf("http://%s/%s", kvserver, objectName)
	return put(url, tmpFile)
}
