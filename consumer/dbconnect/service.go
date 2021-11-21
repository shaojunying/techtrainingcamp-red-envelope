package dbconnect

import (
	"errors"
	"math/rand"
	"redpacket/consumer-demo/database"
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

// 生成红包金额
func GenerateRedEnvelopeValue(remainValue, remainAmount, maxValue, minValue int) int {
	rand.Seed(time.Now().UnixNano())
	averageValue := remainValue / remainAmount
	for i := 0; i < 5; i++ {
		value := int(float64(averageValue) * rand.ExpFloat64())
		if value >= minValue && value <= maxValue {
			return value
		}
	}
	return averageValue
}
