package redenvelope

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"math/rand"
	"red_envelope/database"
	"time"
)

// CheckUserExists 查看用户id是否存在
func (u *User) CheckUserExists() bool {
	db := database.GetDB()

	//事务开始
	tx := db.Begin()
	//同步失败则将数据库回退
	defer func() {
		tx.Rollback()
	}()

	res := !tx.Table("user").Where("uid = ?", u.UID).Find(&u).RecordNotFound()
	tx.Commit()
	return res
}

// AddUser 增加新用户
func (u *User) AddUser() error {
	db := database.GetDB()

	//事务开始
	tx := db.Begin()
	//同步失败则将数据库回退
	defer func() {
		tx.Rollback()
	}()

	if err := tx.Table("user").Create(u).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// CheckSnatchCount 查看抢红包次数
func (u *User) CheckSnatchCount() (int, error) {
	db := database.GetDB()

	//事务开始
	tx := db.Begin()
	//同步失败则将数据库回退
	defer func() {
		tx.Rollback()
	}()

	//此处为了减少字段的写入而将从red_envelope表中直接读取记录条数
	count := 0
	if err := tx.Table("red_envelope").Where("uid = ?", u.UID).Count(&count).Error; err != nil {
		return 0, err
	}
	tx.Commit()
	return count, nil
}

// Distribute 分发红包
func (u *User) Distribute() (int, error) {
	db := database.GetDB()

	//事务开始
	tx := db.Begin()
	//同步失败则将数据库回退
	defer func() {
		tx.Rollback()
	}()
	//根据用户是否领到大红包划定金额范围
	min := viper.GetInt("minAmount")
	max := viper.GetInt("maxAmount")
	value := 0     //要分发的金额
	lucky := false //是否能领大红包
	rand.Seed(time.Now().UnixNano())
	if !u.IfGet {
		num := rand.Intn(100)
		if num < 20 {
			lucky = true
		}
	}
	if lucky {
		value = -(max-min)*2/10*rand.Intn(100)/100 + max // 浮动范围限定在max-0.2(max-min)~max
		if err := tx.Table("user").Where("uid = ?", u.UID).Update("if_get", true).Error; err != nil {
			return 0, err
		}
	} else {
		value = (max-min)*2/10*rand.Intn(100)/100 + min // 浮动范围限定在min~min+0.2(max-min)
	}
	//写入红包表
	redEnvelope := RedEnvelope{UID: *u.UID, Value: value}
	if err := tx.Table("red_envelope").Create(&redEnvelope).Error; err != nil {
		return 0, err
	}
	tx.Commit()
	return redEnvelope.EnvelopeID, nil
}

func (o *OpenRE) Open() (int, error) {
	db := database.GetDB()

	//事务开始
	tx := db.Begin()
	//同步失败则将数据库回退
	defer func() {
		tx.Rollback()
	}()

	// 查看有无数据
	var redEnvelope RedEnvelope
	if tx.Table("red_envelope").Where("envelope_id = ? and uid = ?", o.EnvelopeID, o.UID).Find(&redEnvelope).RecordNotFound() {
		return 0, errors.New("wrong message, no record found")
	}
	// 若已拆开就不用再更新数据库
	if redEnvelope.Opened {
		return redEnvelope.Value, nil
	}
	// red_envelope.opened
	if err := tx.Table("red_envelope").Where("envelope_id = ? and uid = ?", o.EnvelopeID, o.UID).Update("opened", true).Error; err != nil {
		return 0, err
	}
	// 给用户账户加上本次红包的金额
	if err := tx.Table("user").Where("uid = ?", o.UID).Update("amount", gorm.Expr("amount + ?", redEnvelope.Value)).Error; err != nil {
		return 0, err
	}
	tx.Commit()
	return redEnvelope.Value, nil
}

// QueryList 获取钱包列表
func (u *User) QueryList() ([]*WalletList, error) {
	db := database.GetDB()

	//事务开始
	tx := db.Begin()
	//同步失败则将数据库回退
	defer func() {
		tx.Rollback()
	}()

	var rs []*RedEnvelope
	if err := tx.Table("red_envelope").Where("uid = ?", u.UID).Find(&rs).Error; err != nil {
		return nil, err
	}
	//需要对未拆开的红包value做隐藏
	var data []*WalletList
	for _, r := range rs {
		tmp := &WalletList{EnvelopeID: r.EnvelopeID, Opened: r.Opened}
		if r.Opened {
			tmp.Value = r.Value
		}
		tmp.SnatchTime = r.SnatchTime.Unix() //时间戳转换
		data = append(data, tmp)
	}
	tx.Commit()
	return data, nil
}
