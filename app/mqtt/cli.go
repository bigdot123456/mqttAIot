package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func viperconfigDefault() {
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
	viper.SetDefault("testtimes", 0)

	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
}
func viperConfig() {
	viperconfigDefault()
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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("MAC")

	//x:=os.Getenv("MAC_TESTTIMES")
	//y:=viper.GetInt("testtimes")
	//
	////
	////viper.SetEnvPrefix("Baz")
	////os.Setenv("BAZ_BAR", "13")
	////z:=viper.Get("bar")
	//
	//fmt.Printf(x,y)

	viper.OnConfigChange(func(e fsnotify.Event) {
		//viper配置发生变化了 执行响应的操作
		fmt.Println("Config file changed:", e.Name)
	})

}
