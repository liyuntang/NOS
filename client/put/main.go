package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// 定义几个变量
var (
	httpServer string
	dstName string
	files string	// 上传文件的名称
	mode string	// 上传模式，支持随机上传及指定文件上传
	size int	// 如果mode为随机上传的话需要指定随机文件大小
	count int	// 如果mode为随机上传的话需要指定上传的次数
)

func init()  {
	flag.StringVar(&httpServer, "h", "", "存储服务器地址")
	flag.StringVar(&dstName, "n", "", "文件在存储服务器的存储名称，默认同本地文件名")
	flag.StringVar(&files, "f", "", "选择上次的文件，支持单个文件及多文件上次，比如：file01/file01,file02,file03")
	flag.StringVar(&mode, "m", "random", "上传模式，支持随机上传(random)及指定文件上传(file),默认为random")
	flag.IntVar(&size, "s", 1024000, "如果mode为随机上传的话需要指定随机文件大小,单位字节，默认1MB")
	flag.IntVar(&count, "c", 100, "如果mode为随机上传的话需要指定上传的次数,默认100次")
}



func main()  {
	flag.Parse()
	// 判读用户输入是否合法
	dowhat()
	// 走起
	zouqi()
}




func zouqi()  {
	if mode == "file" {
		// 说明是上传指定文件
		startTime := time.Now()
		// 将files接收的参数转化为slice
		fileSlice := fileToSlice()
		// 根据文件获取sha256_code值、文件大小及添加head，并将数据转化为[]byte类型，这些操作似乎都是通用的
		for _, file := range fileSlice {
			// 获取文件大小
			fileSize := getFileSize(file)
			if fileSize == -1 {
				continue
			}
			// 计算sha256_code值
			sha256_code := MakeSha256(dataToByte(file))
			// put
			if filePut(file, fileSize, sha256_code) {
				fmt.Println("put file", file, "is ok, time is", time.Since(startTime))
			} else {
				fmt.Println("put file", file, "is bad, time is", time.Since(startTime))
			}
		}
	} else {
		// 说明要随机上传
		for i:=1;i<=count;i++ {
			startTime := time.Now()
			// 获取文件名称
			objectName := makeString(3)
			fmt.Println("开始上传第", i, "个对象，对象名称为", objectName)
			// 根据指定的文件大小生成数据
			data := makeString(size)
			// 计算sha256_code值
			sha256_code := MakeSha256([]byte(data))
			// put
			if objectPut(objectName, strings.NewReader(data), size, sha256_code) {
				fmt.Println("上传", objectName, "到服务器成功, time is", time.Since(startTime))
			} else {
				fmt.Println("上传", objectName, "到服务器失败 time is", time.Since(startTime))
			}
		}

	}
}

// object put接口
func objectPut(objectName string, data *strings.Reader, fileSize int, sha256_code string) (isok bool) {
	url := fmt.Sprintf("http://%s/%s", httpServer, objectName)
	req, err1 := http.NewRequest("PUT", url, data)
	if err1 != nil {
		fmt.Println("new request is bad, err is", err1)
		return false
	}
	req.Header.Add("FileSize", strconv.Itoa(fileSize))
	req.Header.Add("SncryptionMethod", "sha256")
	req.Header.Add("Sha256_code", sha256_code)
	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		fmt.Println("client do is bad, err is", err2)
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}
// fileput接口
func filePut(file string, fileSize int, sha256_code string) (isok bool) {
	if dstName == "" {
		_, file := path.Split(file)
		dstName = file
	}
	fmt.Println("dstName is", dstName)
	url := fmt.Sprintf("http://%s/%s", httpServer, dstName)
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("打开文件失败", file)
		return false
	}
	req, err1 := http.NewRequest("PUT", url, f)
	if err1 != nil {
		fmt.Println("new request is bad, err is", err1)
		return false
	}

	req.Header.Add("fileSize", strconv.Itoa(fileSize))
	req.Header.Add("sha256_code", sha256_code)
	req.Header.Set("SncryptionMethod", "sha256")
	resp, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		fmt.Println("client do is bad, err is", err2)
		return false
	}
	if resp.StatusCode != 200 {
		fmt.Println("put object", file, "is bad")
		return false
	}
	return true
}

// 计算sha256_code值
func MakeSha256(buf []byte) string {
	h := sha256.New()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 将数据转化为[]byte类型
func dataToByte(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("打开文件失败", file)
		return nil
	}
	buf, err1 := ioutil.ReadAll(f)
	if err1 != nil {
		fmt.Println("读取文件失败", file)
		return nil
	}
	return buf
}

// 获取文件大小
func getFileSize(file string) int {
	fileInfo, err := os.Stat(file)
	if err != nil {
		fmt.Println("获取文件状态失败", file)
		return -1
	}

	n, err1 := strconv.Atoi(strconv.FormatInt(fileInfo.Size(), 10))
	if err1 != nil {
		fmt.Println("get file size is bad, err is", err1)
		return -1
	}
	return n
}


func fileToSlice() []string {
	fileSlice := []string{}
	for _, file := range strings.Split(files, ",") {
		if len(strings.Trim(file, " ")) != 0 {
			// 说明文件名称符合要求，下面需要判读文件是否存在
			_, err := os.Stat(file)
			if !os.IsNotExist(err) {
				// 说明文件存在
				fileSlice = append(fileSlice, file)
			} else {
				// 说明文件不存在，提示用户
				fmt.Println("对不起，文件", file, "不存在，请核对文件")
				os.Exit(0)
			}

		}
	}
	return fileSlice
}



// 判读用户输入是否合法
func dowhat()  {
	if len(os.Args) <= 1 {
		fmt.Println("对不起，请指定运行参数")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if mode == "random" || mode == "file" {
		if mode == "file" {
			// 需要判读用户是否输入的文件
			str := strings.Trim(files, " ")
			str2 := strings.Trim(str, ",")
			str3 := strings.Trim(str2, " ")
			if len(str) == 0 || len(str2) == 0 || len(str3) == 0 {
				// 说明用户输入的file的名字为空格，不符合要求
				fmt.Println("输入如正确的文件名称")
				os.Exit(0)
			}
		}
	} else {
		fmt.Println("对不起，mode参数不符合要求，mode仅支持random、file两个选项")
		os.Exit(0)
	}
}

func makeString(LEN int) (nameString string) {
	list := []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O",
		"P","Q","R","S","T","U","V","W","X","Y","Z","1","2","3","4","5","6","7","8","9","0"}
	var name string
	for i:=1;i<=LEN;i++ {
		index :=rand.Intn(len(list))
		str := list[index]
		name += str
	}
	return name
}
