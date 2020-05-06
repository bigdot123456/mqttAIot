package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
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
	opts.AddBroker("tcp://" + viper.GetString("server.IP") + ":" + viper.GetString("server.port"))

	opts.SetClientID(viper.GetString("client.ID") + IDstr)
	opts.SetUsername(viper.GetString("client.username"))
	opts.SetPassword(viper.GetString("client.passwd"))
	store := viper.GetString("client.store")
	if store != ":memory:" {
		opts.SetStore(mqtt.NewFileStore(store))
	}

	return opts
}

/***
*
* 连接任务和发布消息方法
 */

func mqttConnPubMsgTaskOrg(taskId int, payload string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	//设置连接参数
	opts := getMqttConn(taskId)
	opts.SetClientID(fmt.Sprintf("client%d_%d_%d", taskId, rand.Intn(1000), time.Now().Unix()))

	topic := viper.GetString("client.pubtopic")
	client := mqtt.NewClient(opts)
	//客户端连接判断
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		fmt.Printf("%d:Can't Connect Server....\n", time.Now().Unix())
		//panic(token.Error())
		return
	}
	//fmt.Println("Sample Publisher Started")

	//for i := 0; i < viper.GetInt("client.msgRepeatNum") ; i++ {
	//fmt.Printf("---- doing publish ID:%d round %d ----\n",taskId,i)
	//text:=fmt.Sprintf("ID%d No.%d:%s_",taskId,i,payload)
	token := client.Publish(topic+"/"+strconv.Itoa(taskId)+"/"+deviceInfoStr.CPUID, 0, false, payload)
	token.Wait()
	//}

	client.Disconnect(250)
	//fmt.Printf("\nPublisher %d Disconnected\n",taskId)
}

func mqttConnPubMsgTask(taskId int, payload string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	//设置连接参数
	client := mqttConnTask(taskId)
	if client == nil {
		fmt.Printf("ticked at %v:Please check network connection!\n", time.Now())
		return
	}
	topic := viper.GetString("client.pubtopic")
	token := client.Publish(topic+"/"+strconv.Itoa(taskId)+"/"+deviceInfoStr.CPUID, 0, false, payload)
	token.Wait()
	client.Disconnect(250)
	//fmt.Printf("\nPublisher %d Disconnected\n",taskId)
}

func mqttConnTask(taskId int) mqtt.Client {
	//设置连接参数
	opts := getMqttConn(taskId)
	opts.SetClientID(fmt.Sprintf("client%d_%d_%d", taskId, rand.Intn(1000), time.Now().Unix()))

	client := mqtt.NewClient(opts)
	//客户端连接判断
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		fmt.Printf("%d:Can't Connect Server....\n", time.Now().Unix())
		//panic(token.Error())
		return nil
	}
	return client
}

/***
*
* 连接任务和消息订阅方法
 */
func mqttConnSubMsgTask(taskId int, playload *[]string, waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()
	//设置连接参数
	receiveCount := 0
	opts := getMqttConn(taskId)
	opts.SetClientID(fmt.Sprintf("client%d_%d_%d", taskId, rand.Intn(1000), time.Now().Unix()))
	//设置客户端ID
	//opts.SetClientID(fmt.Sprintf("go Subscribe client example： %d-%d", taskId, time.Now().Unix()))
	//设置连接超时
	//opts.SetConnectTimeout(time.Duration(60) * time.Second)
	//创建客户端连接
	topic := viper.GetString("client.subtopic")

	choke := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		fmt.Printf("%d:Can't Connect Server....\n", time.Now().Unix())
		//panic(token.Error())
		return
	}

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		//os.Exit(1)
		return
	}
	msg := make([]string, 0)
	for receiveCount < viper.GetInt("client.msgnum") {
		incoming := <-choke
		msg = append(msg, incoming[1])
		//fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], msg)
		receiveCount++
	}
	*playload = msg
	client.Disconnect(250)
	//fmt.Printf("[Sub] task %d msg:\t%v\nAddr:%p\n",taskId,playload,playload)
	//return playload

}
