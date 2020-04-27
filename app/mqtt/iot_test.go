package main

import (
	"fmt"
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
	nums := 10
	waitGroup := sync.WaitGroup{}

	for i := 0; i < nums; i++ {
		fmt.Printf("publish client num : %d \n", i)
		waitGroup.Add(1)
		time.Sleep(3 * time.Millisecond)
		//调用连接和发布消息
		go mqttConnPubMsgTask(i, &waitGroup)
		//订阅
		go mqttConnSubMsgTask(i, &waitGroup)
	}

	waitGroup.Wait()
}
