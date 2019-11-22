package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main()  {
	resp, err := http.Get("http://127.0.0.1:9000/aaa.file")
	if err != nil {
		fmt.Println("get is bad, err is", err)
		os.Exit(0)
	}
	fmt.Println("get is ok")
	buf, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println("read resp data is bad, err is", err1)
		os.Exit(0)
	}
	fmt.Println(string(buf))
}


