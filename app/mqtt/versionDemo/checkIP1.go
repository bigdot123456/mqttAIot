package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func main() {
	resp, _ := http.Get("http://www.baidu.com")
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	addr, _ := net.InterfaceAddrs()
	fmt.Printf("内网：%s", addr)

}
