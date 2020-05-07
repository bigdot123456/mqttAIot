package main

import (
	"encoding/json"
	"strconv"
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
	i := 0
	//设置连接参数
	client := mqttConnTask(100)

	topic := viper.GetString("client.pubtopic")
	//fmt.Printf("\nPublisher %d Disconnected\n",taskId)

	for {
		n := viper.GetInt("client.TimeIntervel")
		if n == 0 {
			n = 1
			fmt.Printf("Error in reading config files ./config/config.toml with TimeIntervel, Please check it! \n")
		}
		ticker := time.NewTicker(time.Duration(n * 1e9))

		select {
		case <-ticker.C:
			//fmt.Printf("ticked at %v\n", time.Now())
			go sendOnePub()
			if client != nil {
				go func() {
					hash := MAChash(deviceInfoStr.CPUID)
					token := client.Publish(topic+"/N/"+deviceInfoStr.CPUID+"/"+GitCommit, 0, false, hash+strconv.Itoa(int(i)))
					fmt.Printf("Hash result is: %s\n", hash)

					token.Wait()
				}()

			}
		}
		i++
		x := viper.GetInt("testtimes")
		if x != 0 {
			if i > x {
				fmt.Println("Finish test!")
				break
			}
		}
		fmt.Printf("No.%d Miner activity with %d seconds interval...\n", i, n)
	}
	client.Disconnect(250)
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
