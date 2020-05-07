package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/jaypipes/ghw"
	"github.com/jeek120/cpuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/snksoft/crc"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

type DeviceInfo struct {
	Key     string `toml:"key"`
	UUID    string `toml:"uuid"`
	CPUID   string `toml:"CPUID"`
	Title   string `toml:"title"`
	IP      string `toml:"IP"`
	IPInt   string `toml:"IPInt"`
	CPUInfo string `toml:"CPUInfo"`
	MACID   string `toml:"MACID"`
	DISKID  string `toml:"DISKID"`
	Msg     string `toml:"msg"`
	OS      string `toml:"os"`
}

var deviceInfoStr DeviceInfo
var printVersion bool

func getCPUID() string {
	ids := [4]uint32{}
	cpuid.Cpuid(&ids, 0)
	//fmt.Printf("%d%d%d%d", ids[0], ids[1], ids[2], ids[3])
	cpustr := fmt.Sprintf("%d%d%d%d", ids[0], ids[1], ids[2], ids[3])

	interfaces, err := net.Interfaces()
	if err != nil {
		panic("MACID error: " + err.Error())
	}

	for _, inter := range interfaces {
		cpustr += fmt.Sprintf("_%s", inter.HardwareAddr)
	}
	cpustr = strings.Replace(cpustr, ":", "", -1)
	return cpustr
}

func getNodeNumbyCPUID() int64 {
	ids := [4]uint32{}
	cpuid.Cpuid(&ids, 0)
	cpustr := fmt.Sprintf("%d%d%d%d", ids[0], ids[1], ids[2], ids[3])
	signByte := []byte(cpustr)
	//hash := md5.New()
	//hash.Write(signByte)
	//r:=hash.Sum(nil)
	//sign := hex.EncodeToString(r)
	//s64, _ := strconv.ParseInt(sign, 10, 64)
	//
	//s:=s64%1024

	hash := crc.NewHash(crc.XMODEM)
	s16 := hash.CalculateCRC([]byte(signByte))
	s := s16 % 1024
	return int64(s)

	//// string到int
	//int, err := strconv.Atoi(string)
	//
	//// string到int64
	//int64, err := strconv.ParseInt(string, 10, 64)
	//
	//// int到string
	//string := strconv.Itoa(int)
	//
	//// int64到string
	//string := strconv.FormatInt(int64,10)

}

func getIP2() string {
	conn, err := net.Dial("udp", "baidu.com:80")
	//conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer conn.Close()
	newIP := strings.Split(conn.LocalAddr().String(), ":")[0]
	fmt.Printf("Machine real IP is %s\n", newIP)
	return newIP
}

func getIP1() string {
	//resp, err := http.Get("http://myexternalip.com/raw")
	resp, err := http.Get("http://ifconfig.io/ip")

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		//os.Exit(1)
		newIP := getIP0()
		return newIP
	}
	defer resp.Body.Close()
	//io.Copy(os.Stdout, resp.Body)
	//os.Exit(0)
	s0, err := ioutil.ReadAll(resp.Body)
	s := string(s0)
	return s
}

func getIPInt() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops:" + err.Error())
		return ""
	}
	s := ""
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				s += ipnet.IP.String() + "\n"
			}
		}
	}
	//os.Stdout.WriteString(s)
	return s
}
func getMACID() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("MACID FAIL: " + err.Error())
	}
	s := ""
	for _, inter := range interfaces {
		//fmt.Println(inter.Name, inter.HardwareAddr)
		s += fmt.Sprintf("Name:\t%s  Address:\t%s\n", inter.Name, inter.HardwareAddr)
	}
	return s
}

//func getDiskID1() string {
//
//	var st syscall.Stat_t
//	err := syscall.Stat("/dev/disk0", &st)
//	if err != nil {
//		panic(err)
//	}
//	return fmt.Sprintf("%+v", st)
//}

func getDiskID() string {
	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
		return "Error diskID"
	}

	//fmt.Printf("%v\n", block)
	s := ""
	for _, disk := range block.Disks {
		s += fmt.Sprintf(" %v\n", disk)
		//for _, part := range disk.Partitions {
		//	fmt.Printf("  %v\n", part)
		//}
	}
	return s
}

//import "github.com/satori/go.uuid"

func getUUID(NodeID int64) string {
	node, err := snowflake.NewNode(NodeID)
	if err != nil {
		fmt.Println(err)
		u1 := uuid.Must(uuid.New(), nil).String()
		//fmt.Printf("UUID: %s", u1)
		return u1
	}
	id := node.Generate().String()
	//fmt.Println("id is:", id)

	return id
}

// cpu info
func getCPUInfo() string {
	cpuInfos, err := cpu.Info()
	if err != nil {
		fmt.Printf("get cpu info failed, err:%v", err)
		return ""
	}
	s := ""
	for _, ci := range cpuInfos {
		//fmt.Println(ci)
		s += ci.String()
	}
	// CPU使用率
	//for {
	//	percent, _ := cpu.Percent(time.Second, false)
	//	fmt.Printf("cpu percent:%v\n", percent)
	//}
	percent, _ := cpu.Percent(time.Second, false)
	s += fmt.Sprintf("\npercent:%f", percent)
	return s
}

// Get preferred outbound ip of this machine
func getIP0() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Printf("get IP0 info failed, err:%v", err)
		return "0.0.0.0"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

//var deviceInfoStr map[string]string

func getDeviceInfo() string {
	deviceInfoStr.Title = viper.GetString("title")
	deviceInfoStr.CPUID = getCPUID()
	deviceInfoStr.MACID = getMACID()
	deviceInfoStr.DISKID = getDiskID()
	deviceInfoStr.CPUInfo = getCPUInfo()
	deviceInfoStr.IP = getIP1()
	deviceInfoStr.IPInt = getIPInt()
	NodeNum := getNodeNumbyCPUID()
	deviceInfoStr.UUID = getUUID(NodeNum)
	deviceInfoStr.OS = runtime.GOOS
	deviceInfoStr.Msg = MAChash(deviceInfoStr.CPUInfo)

	deviceInfoStr.Key = pubkey
	s, _ := json.Marshal(deviceInfoStr)

	return string(s)
}
