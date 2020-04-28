package main

import (
	"encoding/json"
	//"flag"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"sync"
	"time"
)

const testNum = 20

func main() {
	//clientNum := flag.Uint64("clientNum", 10, "client nums")
	//nums := int(*clientNum)
	nums := viper.GetInt("client.msgRepeatNum")
	TimeIntervel := viper.GetInt("client.TimeIntervel")
	waitGroup := &sync.WaitGroup{}
	//fmt.Printf("org Address:%p,%p",&msg,&msg[0])
	deviceInfoStr.Msg = "hello"
	mytime := TimeIntervel * 10e9
	//新建计时器，120秒以后触发，go触发计时器的方法比较特别，就是在计时器的channel中发送值
	//tick := time.NewTicker(time.Duration(TimeIntervel))
	tick := time.NewTicker(time.Duration(mytime))
	Num := int64(0)
	for {
		select {
		//此处在等待channel中的信号，因此执行此段代码时会阻塞120秒
		case <-tick.C:
			//调用连接和发布消息
			for i := 0; i < nums; i++ {
				waitGroup.Add(1)
				deviceInfoStr.Key = strconv.Itoa(i)
				payload, _ := json.Marshal(deviceInfoStr)
				go mqttConnPubMsgTask(i, string(payload), waitGroup)
			}
			waitGroup.Wait()
		}
		Num += 1
		fmt.Printf("Send No.%d Message to host... \n", Num)
	}

}
