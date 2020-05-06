package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/jaypipes/ghw"
	"github.com/jeek120/cpuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/snksoft/crc"
	"runtime"
	"strings"
	"time"
	//"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"io/ioutil"
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
	Key     string `toml:"key"`
	CPUID   string `toml:"CPUID"`
	Title   string `toml:"title"`
	IP      string `toml:"IP"`
	IPInt   string `toml:"IPInt"`
	CPUInfo string `toml:"CPUInfo"`
	MACID   string `toml:"MACID"`
	DISKID  string `toml:"DISKID"`
	UUID    string `toml:"uuid"`
	Msg     string `toml:"msg"`
	OS      string `toml:"os"`
}

var deviceInfoStr DeviceInfo
var printVersion bool

//需要赋值的变量
var (
	//Version 项目版本信息
	Version = ""
	//GoVersion Go版本信息
	GoVersion = ""
	//GitTag Tag id
	GitTag = ""
	//GitCommit git提交commmit id
	GitCommit = ""
	//BuildTime 构建时间
	BuildTime = ""
	//Author 作者
	Author = ""
)

func init() {
	flag.BoolVar(&printVersion, "version", false, "print program build version")
	flag.Parse()
	if printVersion {
		PrintVersion()
	}
	fmt.Printf("Start MAC miner Monitor Process...!\n")
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
		//log.Fatal(err) // 读取配置文件失败致命错误
		fmt.Printf("Config file doesn't exist, Use default setting! Miner config maybe wrong!")
	}
	viper.WatchConfig()
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MAC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.OnConfigChange(func(e fsnotify.Event) {
		//viper配置发生变化了 执行响应的操作
		fmt.Println("Config file changed:", e.Name)
	})
	viper.SetDefault("title", "MAC V0 Miner")
	viper.SetDefault("Server.IP", "39.99.160.245")
	viper.SetDefault("Server.port", 1883)

	viper.SetDefault("client.username", "userA")
	viper.SetDefault("client.passwd", "userfast")
	viper.SetDefault("client.subtopic", "mtopic")
	viper.SetDefault("client.pubtopic", "mtopic")
	viper.SetDefault("client.store", ":memory:")
	viper.SetDefault("client.ID", "MACID")
	viper.SetDefault("client.msgRepeatNum", 1)
	viper.SetDefault("client.TimeIntervel", 10)

	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	s := getDeviceInfo()
	WriteWithIoutil("SysInfo.json", s)
}

//PrintVersion 输出版本信息
func PrintVersion() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Go Version: %s\n", GoVersion)
	fmt.Printf("HASH: %s\n", GitCommit)
	fmt.Printf("Tag: %s\n", GitTag)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Author: %s\n", Author)
	os.Exit(0)
}

func getCPUID() string {
	ids := [4]uint32{}
	cpuid.Cpuid(&ids, 0)
	//fmt.Printf("%d%d%d%d", ids[0], ids[1], ids[2], ids[3])
	cpustr := fmt.Sprintf("%d%d%d%d", ids[0], ids[1], ids[2], ids[3])

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
	deviceInfoStr.Msg = GitCommit
	deviceInfoStr.Key = MAChash(deviceInfoStr.CPUInfo)
	s, _ := json.Marshal(deviceInfoStr)

	return string(s)
}

func WriteWithIoutil(name, content string) {
	data := []byte(content)
	if ioutil.WriteFile(name, data, 0644) == nil {
		fmt.Println("\nWrite Device Info file %s:\n", name, content)
	}
}
