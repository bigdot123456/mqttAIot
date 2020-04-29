package main

import (
	"encoding/json"
	"time"

	//"flag"
	"fmt"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"sync"
)

//const testNum = 20
func main() {

	//aChan := make(chan int, 1)
	var i int64 = 0

	for {
		n := viper.GetInt("client.TimeIntervel")
		ticker := time.NewTicker(time.Duration(n * 1e9))

		select {
		case <-ticker.C:
			//fmt.Printf("ticked at %v\n", time.Now())
			sendOnePub()
		}
		i++
		fmt.Printf("Run No.%d publish with %d seconds interval...\n", i, n)
	}
}

func main0() {
	c := cron.New()
	c.AddFunc("* * * * * *", sendOnePub)
	c.Start()
	select {}
}

func sendOnePub() {
	nums := viper.GetInt("client.msgRepeatNum")

	waitGroup := &sync.WaitGroup{}
	deviceInfoStr.Key = "0"
	payload, _ := json.Marshal(deviceInfoStr)
	waitGroup.Add(1)
	mqttConnPubMsgTask(0, string(payload), waitGroup)
	cpuStr := getCPUID()
	for i := 0; i < nums; i++ {
		waitGroup.Add(1)
		text := fmt.Sprintf("{\"Key\":\"%s%d\"}", cpuStr, i)
		go mqttConnPubMsgTask(i+1, text, waitGroup)
	}
	waitGroup.Wait()
	//fmt.Println("Finish one test!")
}
