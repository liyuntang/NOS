package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
)

func main()  {
	resp, err := http.Get("http://127.0.0.1:9000/picture")
	if err != nil {
		fmt.Println(err)
	}
	file, _ := os.OpenFile("/Users/liyuntang/Desktop/ccc.png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer file.Close()
	buf, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(buf)
	fmt.Println(reflect.TypeOf(buf))
	fmt.Println(len(buf))
	n, _ := file.Write(buf)
	fmt.Println("n is", n)
}


