package main

// place it into mqtt directory, it will be ok
import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const testNum = 20

func main() {
	clientNum := flag.Uint64("clientNum", 10, "client nums")
	nums := int(*clientNum)

	msg := [testNum][]string{}
	waitGroup := &sync.WaitGroup{}
	//fmt.Printf("org Address:%p,%p",&msg,&msg[0])
	deviceInfoStr.Msg = "hello"

	for i := 0; i < nums; i++ {
		fmt.Printf("publish client num : %d \n", i)
		waitGroup.Add(2) // should be 2 since it run with two go routine
		time.Sleep(3 * time.Millisecond)
		//调用连接和发布消息
		deviceInfoStr.Key = strconv.Itoa(i)
		payload, _ := json.Marshal(deviceInfoStr)
		go mqttConnPubMsgTask(i, string(payload), waitGroup)
		//订阅
		go mqttConnSubMsgTask(i, &(msg[i]), waitGroup)
	}

	waitGroup.Wait()
	fmt.Printf("finish test, %+v ", msg)
}
