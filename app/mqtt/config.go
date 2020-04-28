package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/fsnotify/fsnotify"
	"github.com/shirou/gopsutil/cpu"
	"github.com/snksoft/crc"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jeek120/cpuid"
	//"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type ConfigInfo struct {
	Title  string `toml:"title"`
	Server struct {
		IP            string `toml:"IP"`
		Port          int    `toml:"port"`
		ConnectionMax int    `toml:"connection_max"`
	} `toml:"server"`
	Client struct {
		Username string `toml:"username"`
		Passwd   string `toml:"passwd"`
		Topic    string `toml:"topic"`
	} `toml:"client"`
	Mysql struct {
		IP     string `toml:"ip"`
		Port   int    `toml:"port"`
		Name   string `toml:"name"`
		Passwd string `toml:"passwd"`
	} `toml:"mysql"`
}

type DeviceInfo struct {
	Title   string `toml:"title"`
	IP      string `toml:"IP"`
	IPInt   string `toml:"IPInt"`
	CPUID   string `toml:"CPUID"`
	CPUInfo string `toml:"CPUInfo"`
	MACID   string `toml:"MACID"`
	DISKID  string `toml:"DISKID"`
	UUID    string `toml:"uuid"`
}
var deviceInfoStr DeviceInfo

func init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	viper.AddConfigPath(path + "/config")
	viper.SetConfigName("config") //指定配置文件的文件名称(不需要制定配置文件的扩展名)
	//viper.AddConfigPath("/etc/appname/")   //设置配置文件的搜索目录
	//viper.AddConfigPath("$HOME/.appname")  // 设置配置文件的搜索目录
	viper.AddConfigPath(".")   // 设置配置文件和可执行二进制文件在用一个目录
	err = viper.ReadInConfig() // 根据以上配置读取加载配置文件
	if err != nil {
		log.Fatal(err) // 读取配置文件失败致命错误
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		//viper配置发生变化了 执行响应的操作
		fmt.Println("Config file changed:", e.Name)
	})
	viper.SetDefault("Server.IP", "127.0.0.1")
	viper.SetDefault("Server.port", 1883)
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	getDeviceInfo()
}

func getCPUID() string {
	ids := [4]uint32{}
	cpuid.Cpuid(&ids, 0)
	//fmt.Printf("%d%d%d%d", ids[0], ids[1], ids[2], ids[3])
	cpustr := fmt.Sprintf("%d%d%d%d", ids[0], ids[1], ids[2], ids[3])

	return cpustr
}
func getNodeNumbyCPUID() int64{
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
	s:=s16%1024
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


func getIP() string {
	resp, err := http.Get("http://myexternalip.com/raw")

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		//os.Exit(1)
		return ""
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
func getDiskID() string {

	var st syscall.Stat_t
	err := syscall.Stat("/dev/disk0", &st)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%+v", st)
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
func getCPUInfo()string {
	cpuInfos, err := cpu.Info()
	if err != nil {
		fmt.Printf("get cpu info failed, err:%v", err)
		return ""
	}
	s:=""
	for _, ci := range cpuInfos {
		//fmt.Println(ci)
		s+=ci.String()
	}
	// CPU使用率
	//for {
	//	percent, _ := cpu.Percent(time.Second, false)
	//	fmt.Printf("cpu percent:%v\n", percent)
	//}
	percent, _ := cpu.Percent(time.Second, false)
	s+=fmt.Sprintf("\npercent:%f",percent)
	return s
}

//var deviceInfoStr map[string]string

func getDeviceInfo()string {
	deviceInfoStr.CPUID=getCPUID()
	deviceInfoStr.MACID=getMACID()
	deviceInfoStr.DISKID=getDiskID()
	deviceInfoStr.CPUInfo=getCPUInfo()
	deviceInfoStr.IP=getIP()
	deviceInfoStr.IPInt=getIPInt()
	NodeNume:=getNodeNumbyCPUID()
	deviceInfoStr.UUID=getUUID(NodeNume)

	s,_  := json.Marshal(deviceInfoStr)
	return string(s)
}