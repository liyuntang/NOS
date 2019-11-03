package operation

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

// 加载一些变量
var (
	WriteLog *log.Logger
	PdServer []string
)

// handler的处理过程是这样的：
// 1、获取object_name（默认所有的dataDir全部相同，其实不相同也可以）
// 2、从pdServer中随机选择一台pd，向该pd发起search请求
// 3、通过pd向所有的kvServer发布广播（默认所有的kvServer都正常提供服务），查询该object所在的主机，如果该主机存在该object则返回一个ok消息，不存在则静默或者返回not found信息
// 4、拿到pd返回的kvserver地址，连接到该kvserver进行相应的操作（get、put、delete）
// 5、返回执行结果

func Handler(w http.ResponseWriter, r *http.Request)  {
	// 获取请求方法
	method := r.Method
	fmt.Println("method is", method)
	// 获取object_name
	objectName := fmt.Sprintf(strings.Split(r.URL.EscapedPath(),"/")[1])
	fmt.Println("objectName is", objectName, "method is", r.Method)

	// 从pdServer中随机选择一台pd,通过该pd向所有的kvServer发布广播
	pdserver := PdServer[rand.Intn(len(PdServer))]
	pdUrl := fmt.Sprintf("http://%s/search/%s", pdserver, objectName)
	fmt.Println("开始广播", pdUrl)
	reps, err := http.Get(pdUrl)
	if err != nil {
		// 说明向pd发起search操作失败
		w.WriteHeader(http.StatusNotAcceptable)
		WriteLog.Println("pd", pdserver, "is not accept")
		return
	}
	// 说明向pd发起search操作成功，此时pd返回的结果有三种：1、正常的kvserver地址  2、kvserver不可达  3、method不对
	if reps.StatusCode == 200 {
		// 说明获取object所在的kvserver很成功，拿到该kvserver地址后发起get、put、delte操作即可
		kvMessage, err := ioutil.ReadAll(reps.Body)
		reps.Body.Close()
		if err != nil {
			WriteLog.Println("get kvMessage from reps is bad")
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		// 对kvserver发起get/put/delete请求
		fmt.Println("kvserver endpoint is", string(kvMessage))
		// 拼接url
		kvUrl := fmt.Sprintf("http://%s/%s", kvMessage, objectName)
		if strings.ToLower(method) == "get"{
			// get object
			reps, err := http.Get(kvUrl)
			if err != nil {
				// 说明从kvserver获取object数据失败
				WriteLog.Println("get object from kvserver", kvMessage, "is bad")
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			// 说明从kvserver获取object数据成功，判断statucode
			if reps.StatusCode == 200 {
				// 说明返回数据成功
				objectData, err := ioutil.ReadAll(reps.Body)
				reps.Body.Close()
				if err != nil {
					WriteLog.Println("get object from reps is bad")
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
				w.WriteHeader(http.StatusOK)
				//fmt.Println("开始给httpServer返回kvServer地址", kvMessage, string(kvMessage))
				w.Write(objectData)
				return
			}

		} else if strings.ToLower(method) == "put" {
			// put object

		} else if strings.ToLower(method) == "delete" {
			// delete object

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			WriteLog.Println("method is not allowd")
			return
		}


	} else if reps.StatusCode == 404 {
		w.WriteHeader(404)
		return
	}else {
		// 说明发起的不是get操作，返回报错
		w.WriteHeader(http.StatusMethodNotAllowed)
		WriteLog.Println("method is not allowd")
		return
	}



}
