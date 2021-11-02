package snatch

import (
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

	res := tx.Table("user").Where("uid = ?", u.UID).Find(&u).RecordNotFound()
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
	money := 0     //要分发的金额
	lucky := false //是否能领大红包
	rand.Seed(time.Now().UnixNano())
	if !u.IfGet {
		num := rand.Intn(100)
		if num < 20 {
			lucky = true
		}
	}
	if lucky {
		money = -(max-min)*2/10*rand.Intn(100)/100 + max // 浮动范围限定在max-0.2(max-min)~max
		if err := tx.Table("user").Where("uid = ?", u.UID).Update("if_get", true).Error; err != nil {
			return 0, err
		}
	} else {
		money = (max-min)*2/10*rand.Intn(100)/100 + min // 浮动范围限定在min~min+0.2(max-min)
	}
	//写入红包表
	redEnvelope := RedEnvelope{UID: *u.UID, Money: money}
	if err := tx.Table("red_envelope").Create(&redEnvelope).Error; err != nil {
		return 0, err
	}
	tx.Commit()
	return redEnvelope.EnvelopeID, nil
}
