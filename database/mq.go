package database

import (
	"errors"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

var mq rocketmq.Producer

func InitMQ() error {
	endPoint := []string{"http://MQ_INST_8149062485579066312_2634790851.cn-beijing.rocketmq-internal.ivolces.com:24009"}
	rlog.SetLogger(nil)
	credentials := primitive.Credentials{
		AccessKey: "ZgybYsmKUsSEYrfExwbryNa1",
		SecretKey: "vdtXqPm1cm6njun1zumNTWzk",
	}
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(endPoint),
		producer.WithGroupName("GID_001"),
		producer.WithCredentials(credentials),
		producer.WithNamespace("MQ_INST_8149062485579066312_2634790851"),
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
