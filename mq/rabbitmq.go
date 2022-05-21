package mq

import (
	"dousheng/config"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var AmqpClient *amqp.Connection

func InitAmqp(config *config.RabbitMQConfig) {
	var err error
	if AmqpClient, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", config.User, config.Pwd, config.Host, config.Vhost)); err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("amqp inti success")
}
