package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main()  {

	for {
		// 对象名称
		objectName := makeString(10)
		// 对象数据
		data := strings.NewReader(makeString(10240))
		// 拼接url
		url := fmt.Sprintf("http://10.10.10.199:9000/%s", objectName)
		req, err := http.NewRequest("PUT", url, data)
		if err != nil {
			fmt.Println("new request is bad, err is", err)
			os.Exit(0)
		}
		// PUT
		startTime := time.Now()
		resp, err1 := http.DefaultClient.Do(req)
		if err1 != nil {
			fmt.Println("do is bad, err is", err1)
		}
		fmt.Println("code is", resp.StatusCode, "run time is", time.Since(startTime))
	}

}

func makeString(LEN int) (nameString string) {
	list := []string{"","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O",
		"P","Q","R","S","T","U","V","W","X","Y","Z"}
	var name string
	for i:=1;i<=LEN;i++ {
		index :=rand.Intn(25)
		str := list[index]
		name += str
	}
	return name
}
