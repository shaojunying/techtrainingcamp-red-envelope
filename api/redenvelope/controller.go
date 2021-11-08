package redenvelope

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// SnatchRedEnvelope 抢红包
func SnatchRedEnvelope(c *gin.Context) {
	var user User
	//匹配参数
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 输入参数有误",
			"data": err,
		})
		return
	}
	if user.UID == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 未获取到uid",
			"data": nil,
		})
		return
	}
	log.Printf("用户%d开始抢红包\n", *user.UID)

	// 每个用户可抢红包数限额
	maxCount := c.GetInt(MaxCountField)
	log.Printf("成功获取最大的可抢红包限额: %d\n", maxCount)
	// 获取当前用户已抢红包的数量
	curCount, err := Mapper.GetRedEnvelops(c, *user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 查询用户已抢红包数失败",
			"data": err,
		})
		return
	}
	log.Printf("成功获取用户%d已抢红包数目: %d\n", *user.UID, maxCount)
	// 判断用户是否超过红包数限额
	if curCount >= maxCount {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 用户已抢红包数达到限额",
			"data": nil,
		})
		return
	}

	// 尝试增加已抢红包数
	curCount, err = Mapper.IncreaseRedEnvelopes(c, *user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 增加已抢红包数失败",
			"data": err,
		})
		return
	}

	// 可能因为并发抢红包的情况，导致用户已抢的红包数超过限额，这时候需要减少已抢红包数（否则配置更新将会出错）
	if curCount > maxCount {
		err = Mapper.DecreaseRedEnvelopes(c, *user.UID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "error, 减少已抢红包数失败",
				"data": err,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 用户已抢红包数超过限额",
			"data": nil,
		})
		return
	}

	// 成功增加了已抢红包数量，生成红包id并添加到set中

	// 生成新红包的id
	envelopeID, err := Mapper.GenerateNewRedEnvelopeId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 生成新红包id失败",
			"data": err,
		})
		return
	}
	err = Mapper.AddRedEnvelopeToUserId(c, *user.UID, envelopeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 添加用户红包id失败",
			"data": err,
		})
		return
	}

	// TODO 将红包、用户信息写入MQ

	data := SuccessSnatch{envelopeID, maxCount, curCount}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
	return
}

// OpenRedEnvelope 拆红包
func OpenRedEnvelope(c *gin.Context) {
	var openre OpenRE
	//匹配参数
	if err := c.ShouldBind(&openre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 输入参数有误",
			"data": err,
		})
		return
	}
	if openre.UID == nil || openre.EnvelopeID == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 未能获取全部参数",
			"data": nil,
		})
		return
	}
	// 判断userId和envelopeId是否匹配
	owned, err := Mapper.CheckIfOwnRedEnvelope(c, *openre.UID, *openre.EnvelopeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 查询用户是否拥有红包失败",
			"data": err,
		})
		return
	}
	if !owned {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 用户未拥有该红包或该红包已被拆开",
			"data": nil,
		})
		return
	}

	// 用户拥有该红包，尝试拆红包
	err = Mapper.RemoveRedEnvelopeForUser(c, *openre.UID, *openre.EnvelopeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 拆红包失败",
			"data": err,
		})
		return
	}

	// TODO 生成红包的金额
	money := -1

	// TODO 将红包id、红包金额写入MQ

	data := SuccessOpen{money}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
	return
}

// GetWalletList 钱包列表
func GetWalletList(c *gin.Context) {
	var user User
	//匹配参数
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 输入参数有误",
			"data": err,
		})
		return
	}
	if user.UID == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 未获取到uid",
			"data": nil,
		})
		return
	}
	if !user.CheckUserExists() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 100,
			"msg":  "error, 用户不存在",
			"data": nil,
		})
		return
	}
	list, err := user.QueryList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 100,
			"msg":  "error, 数据库查询列表失败",
			"data": nil,
		})
		return
	}
	data := &SuccessGet{user.Amount, list}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
	return
}
