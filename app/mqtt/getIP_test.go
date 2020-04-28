package main

import (
	"fmt"
	"testing"
)

//var get_ip = flag.String("get_ip", "", "external|internal")

func TestIP(t *testing.T) {
	//fmt.Println("Usage of ./getmyip --get_ip=(external|internal)")
	//flag.Parse()
	ex := get_external()
	in := get_internal()
	fmt.Printf("extiP:\t%s\ninIP:\t%s", ex, in)

}

func demo_str() {
	str2 := "hello"
	data2 := []byte(str2)
	fmt.Println(data2)
	str2 = string(data2[:])
	fmt.Println(str2)
}

//func get_external() string {
//	resp, err := http.Get("http://myexternalip.com/raw")
//
//	if err != nil {
//		os.Stderr.WriteString(err.Error())
//		os.Stderr.WriteString("\n")
//		//os.Exit(1)
//		return ""
//	}
//	defer resp.Body.Close()
//	//io.Copy(os.Stdout, resp.Body)
//	//os.Exit(0)
//	s0, err := ioutil.ReadAll(resp.Body)
//	s:=string(s0)
//	return s
//}
//
//func get_internal() string {
//	addrs, err := net.InterfaceAddrs()
//	if err != nil {
//		os.Stderr.WriteString("Oops:" + err.Error())
//		return ""
//	}
//	s := ""
//	for _, a := range addrs {
//		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//			if ipnet.IP.To4() != nil {
//				s += ipnet.IP.String() + "\n"
//			}
//		}
//	}
//	//os.Stdout.WriteString(s)
//	return s
//}
