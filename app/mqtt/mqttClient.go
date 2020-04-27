package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"sync"
	"time"
)

//创建全局mqtt publish消息处理 handler
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Pub Client Topic : %s \n", msg.Topic())
	fmt.Printf("Pub Client msg : %s \n", msg.Payload())
}

//创建全局mqtt sub消息处理 handler
var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Sub Client Topic : %s \n", msg.Topic())
	fmt.Printf("Sub Client msg : %s \n", msg.Payload())
}

//连接失败数
var failNums = 0

func getMqttConn(taskId int) *mqtt.ClientOptions {
	IDstr := fmt.Sprintf("%d", taskId)
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://"+viper.GetString("server.IP")+":"+viper.GetString("server.port"))

	opts.SetClientID(viper.GetString("client.ID") + IDstr)
	opts.SetUsername(viper.GetString("client.username"))
	opts.SetPassword(viper.GetString("client.passwd"))
	opts.SetStore(mqtt.NewFileStore(viper.GetString("client.store")))
	return opts
}

/***
*
* 连接任务和发布消息方法
 */

func mqttConnPubMsgTask(taskId int, waitGroup *sync.WaitGroup) {
	//defer waitGroup.Done()
	//设置连接参数
	opts := getMqttConn(taskId)
	opts.SetClientID(fmt.Sprintf("client%d_%d_%d",taskId,rand.Intn(1000),time.Now().Unix()))

	topic := viper.GetString("client.topic")
	client := mqtt.NewClient(opts)
	//客户端连接判断
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Sample Publisher Started")
	payload:="publish msg"+getCPUIDStr()
	for i := 0; i < viper.GetInt("client.msgnum") ; i++ {
		fmt.Printf("---- doing publish ID:%d round %d ----\n",taskId,i)
		text:=fmt.Sprintf("Round %d:%s",i,payload)
		token := client.Publish(topic, 0, false, text)
		token.Wait()
	}

	client.Disconnect(250)
	fmt.Printf("\nPublisher %d Disconnected",taskId)
}

/***
*
* 连接任务和消息订阅方法
 */
func mqttConnSubMsgTask(taskId int, waitGroup *sync.WaitGroup) {
	//defer waitGroup.Done()
	//设置连接参数
	receiveCount := 0
	opts := getMqttConn(taskId)
	opts.SetClientID(fmt.Sprintf("client%d_%d_%d",taskId,rand.Intn(1000),time.Now().Unix()))
	//设置客户端ID
	//opts.SetClientID(fmt.Sprintf("go Subscribe client example： %d-%d", taskId, time.Now().Unix()))
	//设置连接超时
	//opts.SetConnectTimeout(time.Duration(60) * time.Second)
	//创建客户端连接
	topic := viper.GetString("client.topic")

	choke := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for receiveCount < viper.GetInt("client.msgnum") {
		incoming := <-choke
		fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		receiveCount++
	}

	client.Disconnect(250)
	fmt.Println("[Sub] task is ok")

}
