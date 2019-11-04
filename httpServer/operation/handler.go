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
// 这里要对用户发起的请求操作做排查，如果发起的是get、delete操作则需要pd对所有的kvserver进行广播，如果发起的是put请求，则只需要pd返回一个可用的kvserver即可
func Handler(w http.ResponseWriter, r *http.Request)  {
	// 获取请求方法
	method := r.Method
	fmt.Println("method is", method)
	// 对请求方法做判断，如果发起的是get、delete操作则需要pd对所有的kvserver进行广播，如果发起的是put请求，则只需要pd返回一个可用的kvserver即可
	if strings.ToLower(method) == "get" || strings.ToLower(method) == "delete"{
		// 说明发起的是get、delete操作,此时需要pd对所有的kvserver进行广播
		// 获取object_name
		objectName := fmt.Sprintf(strings.Split(r.URL.EscapedPath(),"/")[1])
		// 从pdServer中随机选择一台pd,通过该pd向所有的kvServer发布广播
		pdserver := PdServer[rand.Intn(len(PdServer))]
		pdUrl := fmt.Sprintf("http://%s/search/%s", pdserver, objectName)
		reps, err := http.Get(pdUrl)
		defer reps.Body.Close()
		if err != nil {
			// 说明向pd发起search操作失败
			w.WriteHeader(http.StatusNotAcceptable)
			WriteLog.Println("pd", pdserver, "is not accept")
			return
		}
		fmt.Println("pd code is", reps.StatusCode)
		// 说明向pd发起search操作成功，此时pd返回的结果有两种：1、正常的kvserver地址  2、kvserver不可达
		if reps.StatusCode == 200 {
			// 说明获取object所在的kvserver很成功，拿到该kvserver地址后发起get、put、delte操作即可
			kvMessage, err := ioutil.ReadAll(reps.Body)
			if err != nil {
				WriteLog.Println("get kvMessage from reps is bad")
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			// 对kvserver发起get/put/delete请求
			//fmt.Println("kvserver endpoint is", string(kvMessage))
			// 拼接url
			kvUrl := fmt.Sprintf("http://%s/%s", kvMessage, objectName)
			if strings.ToLower(method) == "get"{
				// get object
				fmt.Println("kvUrl is", kvUrl)
				reps, err := http.Get(kvUrl)
				if err != nil {
					// 说明从kvserver获取object数据失败
					WriteLog.Println("get object from kvserver", kvMessage, "is bad")
					w.WriteHeader(http.StatusNotAcceptable)
					return
				}
				// 说明从kvserver获取object数据成功，判断statucode
				// statucode有两种情况：1、object存在 2object不存在
				if reps.StatusCode == 200 {
					// object存在
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
				} else if reps.StatusCode == 404 {
					// object不存在
					WriteLog.Println("object", objectName, "is not exist")
					w.WriteHeader(404)
				}

			} else if strings.ToLower(method) == "delete" {
				// put object
				req, err := http.NewRequest("delete", kvUrl, nil)
				if err != nil {
					// 说明项kvserver发起delete请求操作失败
					WriteLog.Println("new request of delete is bad")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				// 说明项kvserver发起delete请求操作成功
				resp, err1 := http.DefaultClient.Do(req)
				defer resp.Body.Close()
				if err1 != nil {
					WriteLog.Println("run default client is bad, err is", err1)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				// 根据code来判读delete操作是否成功
				// code有三个状态：1、删除成功  2、文件存在但删除失败  3、文件不存在
				code := resp.StatusCode
				if code == 200 {
					// 删除成功
					WriteLog.Println("delete object", objectName, "is ok")
					w.WriteHeader(http.StatusOK)

				} else if code == 400 {
					// 文件存在但删除失败
					WriteLog.Println("object", objectName, "is exist but delete it is bad")
					w.WriteHeader(http.StatusBadRequest)

				}

			}
		} else if reps.StatusCode == 404 {
			// object不存在
			WriteLog.Println("object", objectName, "is not exist")
			w.WriteHeader(404)
			return
		}else {
			// 说明发起的不是get操作，返回报错
			w.WriteHeader(http.StatusMethodNotAllowed)
			WriteLog.Println("method is not allowd")
			return
		}
	} else if strings.ToLower(method) == "put" {
		// 说明发起的是put操作，此时只需要pd返回一个可用的kvserver即可
		fmt.Println("heheh你发起的是put操作")
		// 从pdServer中随机选择一台pd,通过该pd向所有的kvServer发布广播
		pdserver := PdServer[rand.Intn(len(PdServer))]
		pdUrl := fmt.Sprintf("http://%s/put/", pdserver, )
		fmt.Println("开始广播", pdUrl)
		resp, err := http.Get(pdUrl)
		if err != nil {
			// 说明请求pd不通
			WriteLog.Println("sorry, pdserver", pdserver, "is unreachable")
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		// 说明请求pd正常
		defer resp.Body.Close()
		fmt.Println("code is", resp.StatusCode)
		// 根据resp.StatusCode的值进行判断操作
		code := resp.StatusCode
		if code == 200 {
			// 说明返回kvserver地址正常，此时根据返回的kvserver发起put操作
			// 解析kvserver地址
			buf, _ := ioutil.ReadAll(resp.Body)
			kvserver := string(buf)
			// 获取objectName
			objectName := fmt.Sprintf(strings.Split(r.URL.EscapedPath(),"/")[1])
			// 拼接url
			kvUrl := fmt.Sprintf("http://%s/%s", kvserver, objectName)
			// 对kvserver发起put操作请求
			req, err1 := http.NewRequest("PUT", kvUrl, r.Body)
			if err1 != nil {
				// 说明对kvserver发起put操作失败
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			// 说明对kvserver发起put操作成功
			resp, err2 := http.DefaultClient.Do(req)
			if err2 != nil {
				// 说明报错了
				w.WriteHeader(http.StatusNotAcceptable)
				WriteLog.Println("http defaultclient do is bad")
				return
			}
			defer resp.Body.Close()
			code := resp.StatusCode
			if code == 200 {
				// 说明put操作成功
				w.WriteHeader(http.StatusOK)
				WriteLog.Println("put object is ok")
				return
			} else if code == 500 {
				// 说明object文件已经存在，此次操作失败
				w.WriteHeader(http.StatusExpectationFailed)
				WriteLog.Println("sorry object", objectName, "has exist")
				return
			} else if code == 417{
				// 说明写入数据失败
				w.WriteHeader(http.StatusExpectationFailed)
				WriteLog.Println("sorry cp object data to file is bad")
				return
			}


		} else if code == 406 {
			// 说明没有可用的kvserver，也就是说所有的kvserver都宕机了，给用户返回提示
			w.WriteHeader(http.StatusNotAcceptable)
			WriteLog.Println("get kvserver from pd is bad")
			return
		}
	} else {
		// 该程序只支持get、put、delete三个对象操作，其他的操作一律报错
		WriteLog.Println("sorry, your object method", method, "is not allow,we only suport get/put/delete object")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}





















}
