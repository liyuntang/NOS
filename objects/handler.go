package objects

import (
	"log"
	"net/http"
	"strings"
)
// 加载一些变量
var (
	DataDir string
	WriteLog *log.Logger
)



func Handler(w http.ResponseWriter, r *http.Request)  {
	method := r.Method
	if strings.ToLower(method) == "get" {
		get(w, r)
		return
	} else if strings.ToLower(method) == "put" {
		put(w, r)
		return
	} else if strings.ToLower(method) == "delete" {
		delete(w, r)
		return
	}else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
