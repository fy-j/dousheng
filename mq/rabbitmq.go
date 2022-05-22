package mq

import (
	"dousheng/config"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var (
	AmqpClient *amqp.Connection
	//发布视频消息队列管道
	PublishChannel *amqp.Channel
	//发布视频消息队列
	PublishQueue amqp.Queue
)

func InitAmqp(config *config.RabbitMQConfig) {
	var err error
	if AmqpClient, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", config.User, config.Pwd, config.Host, config.Vhost)); err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("amqp inti success")
	PublishChannel, err = AmqpClient.Channel()
	CreatePublishExchange()
	CreatePublishQueue()
	BindingPublish()
}

//创建了一个名字叫publishQueue的消息队列
func CreatePublishQueue() {
	var err error
	PublishQueue, err = PublishChannel.QueueDeclare(
		"publishQueue", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatalln(err)
		fmt.Println("创建队列失败")
		return
	} else {
		fmt.Println("create messageQueue success")
	}
}

//创建了一个叫publishExchange的直连交换机
func CreatePublishExchange() {
	err := PublishChannel.ExchangeDeclare("publishExchange", "direct", true, false, false, true, nil)
	if err != nil {
		fmt.Println("exchange create fill")
	} else {
		fmt.Println("exchange create success")
	}
}

//绑定消息队列和交换机，绑定键为publish
func BindingPublish() {
	err := PublishChannel.QueueBind("publishQueue", "publish", "publishExchange", false, nil)
	if err != nil {
		fmt.Println("publish banding filled")
	} else {
		fmt.Println("publish banding success")
	}
}
