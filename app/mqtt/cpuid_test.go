package main

import (
	"fmt"
	cpuid "github.com/jeek120/cpuid"
	"testing"
	//"time"
)
//import "github.com/shirou/gopsutil/cpu"
func TestCPUID(t *testing.T) {

	ids := [4]uint32{}
	cpuid.Cpuid(&ids, 0)
	//fmt.Printf("%d%d%d%d", ids[0], ids[1], ids[2], ids[3])
	s1:=getCPUInfo()
	s2:=getCPUID()
	s3:=getNodeNumbyCPUID()

	fmt.Printf("CPU type:%s\n\nCPU ID:\t%s\nCPU Hash:\t%d\n",s1,s2,s3)
}

//// cpu info
//func getCpuInfo() {
//	cpuInfos, err := cpu.Info()
//	if err != nil {
//		fmt.Printf("get cpu info failed, err:%v", err)
//	}
//	for _, ci := range cpuInfos {
//		fmt.Println(ci)
//	}
//	// CPU使用率
//	for {
//		percent, _ := cpu.Percent(time.Second, false)
//		fmt.Printf("cpu percent:%v\n", percent)
//	}
//}