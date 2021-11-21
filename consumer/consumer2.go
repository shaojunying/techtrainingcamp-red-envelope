package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"redpacket/consumer-demo/config"
	"redpacket/consumer-demo/database"
	"redpacket/consumer-demo/dbconnect"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	config.InitConf()
	db := database.InitDB()
	defer db.Close()
	endPoint := []string{"http://MQ_INST_8149062485579066312_2634790851.cn-beijing.rocketmq-internal.ivolces.com:24009"}
	credentials := primitive.Credentials{
		AccessKey: "ZgybYsmKUsSEYrfExwbryNa1",
		SecretKey: "vdtXqPm1cm6njun1zumNTWzk",
	}
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(endPoint),
		consumer.WithGroupName("GID_002"),
		consumer.WithCredentials(credentials),
		consumer.WithNamespace("MQ_INST_8149062485579066312_2634790851"),
	)
	if err != nil {
		log.Printf("init producer error: %s", err.Error())
		os.Exit(1)
	}
	err = c.Subscribe("open_value", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			if msg.ReconsumeTimes > 2 {
				continue
			}
			message_body_str := string(msg.Message.Body)
			message_body := strings.Split(message_body_str[1:len(message_body_str)-1], ",")
			pid, _ := strconv.Atoi(message_body[0])
			val, _ := strconv.Atoi(message_body[1])
			log.Printf("Subscribe Callback: %v \n", msg)
			if pid < 0 || val < 0 {
				log.Println("Illegal Data. Dry Run. Not submitted to SQL.")
			} else {
				tmp := dbconnect.RedEnvelope{EnvelopeID: &pid, Value: val}
				err := tmp.OpenSql()
				if err != nil {
					log.Println("Submit to SQL failed.")
					return consumer.ConsumeRetryLater, errors.New(err.Error())
				}
			}
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		log.Printf("start producer error: %s\n", err.Error())
	}
	err = c.Start()
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		log.Printf("Shutdown Consumer error: %s\n", err.Error())
	}
}

// Example
// 2021/11/09 12:11:17 Subscribe Callback: [[Message=[topic=MQ_INST_8149062485579066312_2634790851%snatch_history, body={1,-1}, Flag=0, properties=map[CLUSTER:onlinecluster CONSUME_START_TIME:1636459877399 MAX_OFFSET:9063 MIN_OFFSET:0 UNIQ_KEY:AC15203E0001000000002bd062880001], TransactionId=], MsgId=AC15203E0001000000002bd062880001, OffsetMsgId=AC12DA33000078BF00000393553D7F05,QueueId=1, StoreSize=214, QueueOffset=9062, SysFlag=0, BornTimestamp=1636459877396, BornHost=172.31.200.1:62418, StoreTimestamp=1636459877398, StoreHost=172.18.218.51:30911, CommitLogOffset=3931325169413, BodyCRC=2130576017, ReconsumeTimes=0, PreparedTransactionOffset=0]]
