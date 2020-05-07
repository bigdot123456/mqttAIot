package main

import (
	"flag"
	"fmt"

	"io/ioutil"
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
	pubkey = ""
)

func init() {
	flag.BoolVar(&printVersion, "version", false, "print program build version")
	flag.Parse()
	if printVersion {
		PrintVersion()
	}
	fmt.Printf("Start MAC miner Monitor Process...!\n")

	viperConfig()

	GenRsaKey(2048)
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

func WriteWithIoutil(name, content string) {
	data := []byte(content)
	if ioutil.WriteFile(name, data, 0644) == nil {
		fmt.Printf("\nWrite Device Info file %s:\n %s\n", name, content)
	}
}
