package redenvelope

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"math/rand"
	"net/http"
	"time"
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

	//查看uid是否存在
	count := 0
	maxCount := viper.GetInt("snatchMaxCount")
	if user.CheckUserExists() {
		//检查抢红包次数
		count0, err := user.CheckSnatchCount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 100,
				"msg":  "error, 查询抢红包次数出错",
				"data": err,
			})
			return
		}
		count = count0
		if count >= maxCount {
			c.JSON(http.StatusOK, gin.H{
				"code": 2,
				"msg":  "fail, 该用户抢红包次数达到上限",
				"data": nil,
			})
			return
		}
	} else {
		//不存在则添加用户
		if err := user.AddUser(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 100,
				"msg":  "error, 添加新用户失败",
				"data": err,
			})
			return
		}
	}

	//生成随机数
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(100)                            //随机数
	probability := viper.GetInt("snatchProbability") //抢到的概率值 %
	//fmt.Println(probability)
	if num >= probability {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "fail, 未能抢到红包",
			"data": nil,
		})
		return
	}

	//抢到红包
	envelopeID, err := user.Distribute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 100,
			"msg":  "error, 红包表写入失败",
			"data": err,
		})
		return
	}
	data := SuccessSnatch{envelopeID, maxCount, count + 1}
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
	money, err := openre.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 100,
			"msg":  "error, 不存在相应记录或更新数据库失败",
			"data": nil,
		})
		return
	}
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

// WrkTest 压力测试
func WrkTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}
