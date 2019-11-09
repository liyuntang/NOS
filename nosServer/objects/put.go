package objects

import (
	"NOS/nosServer/encapsulation"
	"NOS/nosServer/etcd"
	"NOS/nosServer/metadata"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
)

// put操作完整流程如下：
//	1、如果bojectName不存在，直接执行put操作
//	2、如果bojectName存在且is_del=0则表示该该object存在，此时返回400错误信息
//	3、如果bojectName存在且is_del=1则表示该该object不存在，直接执行put操作
// 由于objectName与数据存储的名字是解耦的，所以即使object满足上面不存在的条件，我们也要判断数据存储的名字是否存在，也就是sha256_code值存不存在
// 如果sha256_code值存在，则put流程只需要再nos_metadta表中添加一行记录即可
// 如果sha256_code值不存在，则直接执行put操作
// 也就是说针对put操作的情况我们不仅要判断isExist，还要判断sha256_code的值


// put流程分两个操作：
//	1、数据服务层存入数据
//	2、metadata记录元数据
// 只有这两个操作同时成功才表示该object的put成功

//func put(objectName string, isExist bool, objectInfoMap map[string]string, w http.ResponseWriter, data []byte)  {
func put(objectName string, isExist bool, objectInfoMap map[string]string, w http.ResponseWriter, r *http.Request)  {
	if isExist {
		// 说明该object存在，此时返回400
		WriteLog.Println("sorry, object", objectName, "is exist")
		w.WriteHeader(400)
		return
	}

	// 此时需要根据data计算sha256_code值，并以此为objectName进行存储
	// 当我们再往后传递r变量时，发现一个很严重的问题，r.Body的值没了，why?
	// 这个地方为了解决r.Body值的问题，我们暂时引用一个文件保存该值
	data, _ := ioutil.ReadAll(r.Body)
	tmpReader := tmpFile(data)
	sha256Code := metadata.MakeSha256(data)
	//fmt.Println("+++++++++++++++++++", string(data))
	//fmt.Println("-------------------", sha256Code)
	// 说明该object可能不存在，此时需要判断sha256_code的值，如果sha256_code的值为空，则表示该object真的不存在，如果不为空则只需要在nos_metadta表中添加一行记录即可
	if metadata.Sha256CodeISOK(sha256Code) {
		// 说明sha256_code存在，在nos_metadta表中添加一行记录即可
		if metadata.InsertObject(objectName, sha256Code) {
			// 说明在nos_metadta表中添加一行记录成功
			WriteLog.Println("put object", objectName, "is ok")
			w.WriteHeader(200)
			return
		}
		// 在nos_metadta表中添加一行记录失败
		WriteLog.Println("sorry, put object", objectName, "is bad")
		w.WriteHeader(400)
	}
	// 说明sha256_code值不存在，则直接执行put操作
	// 从etcd中获取kvserver信息
	kvservers := etcd.EtcdGet(EtcdServer, WriteLog)
	// 如果method为put操作则随机抽取一台kvserver执行put操作即可
	kvserver := kvservers[rand.Intn(len(kvservers))]
	// 封装put操作,这个地方put操作应该有个返回值，
	if encapsulation.SecondOperationPut(sha256Code, kvserver, tmpReader) {
		// 说明数据层存储object成功，此时需要将objectName、sha256Code记录到元数据里
		if metadata.InsertObject(objectName, sha256Code) {
			// 记录元数据成功，则返回200
			WriteLog.Println("put object", objectName, "is ok")
			w.WriteHeader(200)
			return
		}
		// 记录元数据失败，失败返回417
		WriteLog.Println("put object", objectName, "is bad")
		w.WriteHeader(417)
		return
	}
	// 说明数据层存储object失败,也不用写元数据了，直接返回报错
	WriteLog.Println("put object", objectName, "is bad")
	w.WriteHeader(400)
	return

}

func tmpFile(data []byte) io.Reader {
	file, _ := os.OpenFile("tmp.file", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	_, err := file.Write(data)
	if err != nil {
		return nil
	}
	file.Close()
	f, _ := os.Open("tmp.file")
	return f


}
