package main

import (
	"fmt"
	"testing"
)

func TestNetID(t *testing.T) {
	//interfaces, err :=  net.Interfaces()
	//if err != nil {
	//	panic("Poor soul, here is what you got: " + err.Error())
	//}
	//
	//for _, inter := range interfaces {
	//	fmt.Println(inter.Name, inter.HardwareAddr)
	//}
	//var st syscall.Stat_t
	//err = syscall.Stat("/dev/disk0", &st)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v", st)
	s1:=getDiskID()
	s2:=getMACID()
	fmt.Println(s1,s2)
}
