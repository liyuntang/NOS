package objects

import (
	"NOS/nosServer/encapsulation"
	"NOS/nosServer/metadata"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

// put操作检查流程如下：
//	1、head中是否有"Filesize", "Sha256_code", "Sncryptionmethod"这三个标志
//	2、判读object在metadata中是否存在，该判断主要通过objectName进行的，没有对sha256判断
//	3、将object存放到tmp目录中，计算其size、sha256值与head中的值进行判断
//
//
//
//
// put操作完整流程如下：
//	1、检查客户端发起的put操作是否符合要求，如果不符合要求则直接返回报错，要求如下：
// 		1)head中必须要有对象的大小
// 		2)head中必须要有对象的加密方式（该版本只支持sha256）
// 		3)head中必须要有对象的sha256值，
//	2、检查objectName是否存在，判断标准是：
//		1）如果bojectName不存在，直接执行put操作
//		2）如果bojectName存在且is_del=1则表示该该object不存在，直接执行put操作
//		3）如果bojectName存在且is_del=0则表示该该object存在，此时返回400错误信息
//	3、put操作如下：
//		1）想将object临时存放到nos组件的tmp目录下，
// 		2）计算其sha256值，与head中的sha256值进行对比，如果相同则将该object转存到kvserver中，如果不同则返回报错
//		3) 计算其数据大小，如果与head中的size不同的话，则返回报错，如果相同则转存到kvserver中
// 由于objectName与数据存储的名字是解耦的，所以即使object满足上面不存在的条件，我们也要判断数据存储的名字是否存在，也就是sha256_code值存不存在
// 如果sha256_code值存在，则put流程只需要再nos_metadta表中添加一行记录即可
// 如果sha256_code值不存在，则直接执行put操作
// 也就是说针对put操作的情况我们不仅要判断isExist，还要判断sha256_code的值

// put流程分两个操作：(优先将数据存入kvserver，然后在记录metadata，metadata的优先级小于数据存储)
//	1、数据服务层存入数据
//	2、metadata记录元数据
// 只有这两个操作同时成功才表示该object的put成功

//func put(objectName string, isExist bool, objectInfoMap map[string]string, w http.ResponseWriter, data []byte)  {
func put(objectName string, objectInfoMap map[string]string, w http.ResponseWriter, r *http.Request) {
	// 判断用户的head设置是否符合要求
	for _, key := range []string{"Filesize", "Sha256_code", "Sncryptionmethod"} {
		_, isok := r.Header[key]
		if !isok {
			WriteLog.Println("sorry, head is bad")
			w.WriteHeader(400)
		}
	}
	// 至此，说明request发送的head完全符合我们的要求，下面判读object是否已经存在
	// 不管是get、put还是delete都需要确认object是否存在
	// 如果存在则返回true以及一个存放了object信息的map
	// 如果不存在则返回false以及一个空map
	isok, objectInfoMap := metadata.ObjectISOK(objectName)
	if isok {
		// 说明该object存在，此时返回400
		WriteLog.Println("sorry, object", objectName, "is exist")
		w.WriteHeader(400)
		return
	}
	// 至此，说明用户上传的object不存在，进入真正的put流程，流程如下：
	//	1、将用户上传的object存放到tmp目录下
	//	2、计算其size，并与head.Filesize进行对比，如果相同则进行第3步，如果不同则报错
	//	3、计算其sha256_code值，并与head.Sha256_code进行对比，如果相同则将该object转存到kvserver中，否则返回报错

	// 将输入存入tmp目录下：
	sha256 := r.Header["Sha256_code"][0]
	tmpFile := fmt.Sprintf("%s/%s", TmpDir, sha256)
	num, isok := writeToTmpFile(tmpFile, r)
	if isok {
		// 说明数据写入tmp目录成功,此时可以通过返回的num来验证文件大小
		if num == r.Header["Filesize"][0] {
			// 说明文件大小相同，接下来验证sha256_code
			if !judgeSha256(sha256) {
				// 说明sha256_code不同
				WriteLog.Println("sorry, object", objectName, "sha256Code is bad")
				w.WriteHeader(400)
				return
			}
		} else {
			// 说明文件大小不同，此时需要删除tmp文件，并且返回报错
			deleteFile(tmpFile)
			WriteLog.Println("sorry, object", objectName, "size is bad")
			w.WriteHeader(400)
			return
		}
	} else {
		// 说明写入tmp目录失败,其他的就不用验证了
		WriteLog.Println("sorry, write object", objectName, "to tmp is bad")
		w.WriteHeader(400)
		return
	}
	// 至此，object校验工作全部完成，将该boject转正（存入到kvserver）即可，转正流程如下：
	//	1、将object入库
	//	2、将object数据存入kvserver
	// 此时需要根据data计算sha256_code值，并以此为objectName进行存储
	// 说明该object可能不存在，此时需要判断sha256_code的值，如果sha256_code的值为空，则表示该object真的不存在，如果不为空则只需要在nos_metadta表中添加一行记录即可
	if metadata.Sha256CodeISOK(sha256) {
		// 说明sha256_code存在，在nos_metadta表中添加一行记录即可,此时还要删除tmp文件
		deleteFile(tmpFile) // 删除成功与否无所谓，因为每个write都是truncate
		if metadata.InsertObject(objectName, sha256) {
			// 说明在nos_metadta表中添加一行记录成功
			WriteLog.Println("put object", objectName, "is ok")
			w.WriteHeader(200)
			return
		}
		// 在nos_metadta表中添加一行记录失败
		WriteLog.Println("sorry, put object", objectName, "is bad")
		w.WriteHeader(400)
		return
	}
	// 说明sha256_code值不存在，则直接进行转正操作，该过程如下：
	// 	1、从etcd中获取kvserver信息
	//	2、将tmp/sha256转存到kvserver
	//	3、写入metadata信息
	// 注意这个地方我们先转存object，然后在记录metadata，在最坏的情况下我们宁可metadata没有记录也要把object存入kvserver，

	// 1、根据所设置的副本集的数量（max_replicas）从etcd中获取kvserver信息
	//kvservers := etcd.EtcdGet(EtcdServer, WriteLog)
	kvservers := []string{"10.10.30.202:9100", "10.10.10.69:9100", "10.10.30.202:9100"}
	if len(kvservers) == 0 {
		// 说明没有从etcd中取到kvserver，此时直接报错
		WriteLog.Println("get kvserver is bad")
		w.WriteHeader(400)
		return
	}
	kvServerSlice := []string{}
	for len(kvServerSlice) < MaxReplicas {
		index := rand.Intn(len(kvservers))
		num := kvservers[index]
		if isExsit(num, kvServerSlice) {
			// 说明该kvserver已经存在与切片中
			fmt.Println("kvServerSlice is", kvServerSlice, "num is", num)
			continue
		}
		kvServerSlice = append(kvServerSlice, num)
	}

	// 遍历kvServerSlice，执行put操作即可
	for _, kvserver := range kvServerSlice {
		if !encapsulation.SecondOperationPut(sha256, kvserver, tmpFile) {
			// 说明某一台kvserver执行失败，返回报错，已存储的数据暂时不删除
			deleteFile(tmpFile)
			WriteLog.Println("put object", objectName, "is bad")
			w.WriteHeader(400)
			return
		}
	}
	// 说明数据层存储object成功，此时需要将objectName、sha256Code记录到元数据里
	deleteFile(tmpFile)
	if metadata.InsertObject(objectName, sha256) {
		// 记录元数据成功，则返回200，同时删除tmpFile
		WriteLog.Println("put object", objectName, "is ok")
		w.WriteHeader(200)
		return
	}
	// 记录元数据失败，失败返回417
	WriteLog.Println("put object", objectName, "is bad")
	w.WriteHeader(417)
	return
	}

func judgeSha256(sha256Code string) bool {
	filePath := fmt.Sprintf("%s/%s", TmpDir, sha256Code)

	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		WriteLog.Println("read file", filePath, "is bad, err is", err)
		return false
	}
	code := metadata.MakeSha256(buf)
	if code != sha256Code {
		return false
	}
	return true
}


func deleteFile(fileName string) {
	os.Remove(fileName)
}

func writeToTmpFile(fileName string, r *http.Request) (writeBytes string, isok bool) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		WriteLog.Println("open file", fileName, "is bad")
		return "-1", false
	}
	n, err1 := io.Copy(file, r.Body)
	if err1 != nil {
		// 说明写入失败，
		WriteLog.Println("write to file", fileName, "is bad")
		return "-1", false
	}
	// 说明写入成功
	WriteLog.Println("write file", fileName, "is ok")
	return strconv.FormatInt(n, 10), true

}

func isExsit(iterm string, kvSlice []string) bool {
	for _, value := range kvSlice {
		if value == iterm {
			return true
		}
	}
	return false
}
