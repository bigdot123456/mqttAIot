package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

/***
* 创建客户端连接
 */
func TestIoT(t *testing.T) {
	//clientNum := flag.Uint64("clientNum", 30000, "client nums")
	//flag.Parse()
	//nums := int(*clientNum)
	nums := 20
	msg := [20][]string{}
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
