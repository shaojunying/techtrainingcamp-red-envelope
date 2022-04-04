package database

import (
	"errors"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var mq rocketmq.Producer

func InitMQ() error {
	endPoint := []string{"http://rmqnamesrv:9876"}
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(endPoint),
	)
	if err != nil {
		return errors.New("init producer error: " + err.Error())
	}
	err = p.Start()
	if err != nil {
		return errors.New("start producer error: " + err.Error())
	}
	mq = p
	return nil
}

func GetMQ() rocketmq.Producer {
	return mq
}

func CloseMQ() {
	mq.Shutdown()
}
