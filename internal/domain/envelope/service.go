package envelope

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"log"
	"math/rand"
	"red_envelope/internal/infrastructure/database"
	"time"
)

// SnatchSql 分发红包sql
func (r *RedEnvelope) SnatchSql() error {
	db := database.GetDB()

	//事务开始
	tx := db.Begin()
	//同步失败则将数据库回退
	defer func() {
		tx.Rollback()
	}()

	//写入红包表
	if err := tx.Table("red_envelope").Create(&r).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// OpenSql 拆红包sql
func (r *RedEnvelope) OpenSql() error {
	db := database.GetDB()

	//事务开始
	tx := db.Begin()
	//同步失败则将数据库回退
	defer func() {
		tx.Rollback()
	}()

	r.Opened = true
	if err := tx.Table("red_envelope").Where("envelope_id = ?", r.EnvelopeID).Update(&r).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// QueryListSql 获取钱包列表
func (r *RedEnvelope) QueryListSql() (map[string]interface{}, error) {
	db := database.GetDB()

	//事务开始
	tx := db.Begin()
	//同步失败则将数据库回退
	defer func() {
		tx.Rollback()
	}()

	var list []*WalletList
	// 获取列表
	if err := tx.Table("red_envelope").Where("uid = ?", r.UID).Find(&list).Error; err != nil {
		return nil, err
	}
	tx.Commit()

	//若list为空 则返回错误
	if len(list) == 0 {
		return nil, errors.New("该用户尚未抢过红包")
	}
	// 由于红包金额可能存在null值，gorm求和会出现报错,因此采用遍历获取用户总金额
	amount := 0
	for _, l := range list {
		amount += l.Value
	}
	data := make(map[string]interface{})
	data["amount"] = amount
	data["envelope_list"] = list
	//fmt.Println(data)
	return data, nil
}

// QueryKoiListSql 数据库获取锦鲤列表
func QueryKoiListSql() ([]*Koi, error) {
	db := database.GetDB()

	//事务开始
	tx := db.Begin()
	//同步失败则将数据库回退
	defer func() {
		tx.Rollback()
	}()

	var ks []*Koi
	if err := tx.Table("red_envelope").Raw("select uid, sum(value) as amount from red_envelope group by uid order by amount desc limit 10").Find(&ks).Error; err != nil {
		return nil, err
	}

	tx.Commit()
	return ks, nil
}

// 生成红包金额
func GenerateRedEnvelopeValue(remainValue, remainAmount, maxValue, minValue int) int {
	rand.Seed(time.Now().UnixNano())
	averageValue := remainValue / remainAmount
	value := int(float64(averageValue) * rand.ExpFloat64())
	if value < minValue {
		return minValue
	} else if value > maxValue {
		return maxValue
	} else {
		return value
	}
}

func SnatchHistoryToMQ(uid int, pid int) error {
	data_to_be_sent := fmt.Sprintf("{%d,%d}", uid, pid)
	mq := database.GetMQ()
	result, err := mq.SendSync(context.Background(), primitive.NewMessage("snatch_history", []byte(data_to_be_sent)))
	if err != nil {
		return errors.New("MQ produce error: " + err.Error())
	} else {
		log.Println("MQ produce success: " + result.String())
	}
	return nil
}

func OpenValueToMQ(pid int, value int) error {
	data_to_be_sent := fmt.Sprintf("[%d,%d]", pid, value)
	mq := database.GetMQ()
	result, err := mq.SendSync(context.Background(), primitive.NewMessage("open_value", []byte(data_to_be_sent)))
	if err != nil {
		return errors.New("MQ produce error: " + err.Error())
	} else {
		log.Println("MQ produce success: " + result.String())
	}
	return nil
}
