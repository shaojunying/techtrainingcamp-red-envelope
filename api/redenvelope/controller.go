package redenvelope

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
)

// SnatchRedEnvelope 抢红包
func SnatchRedEnvelope(c *gin.Context) {
	var user User
	//匹配参数
	if err := c.ShouldBind(&user); err != nil {
		log.Printf("输入参数有误:%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  "error, 输入参数有误",
			"data": err,
		})
		return
	}
	if user.UID == nil {
		log.Printf("未获取到uid\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  "error, 未获取到uid",
			"data": nil,
		})
		return
	}
	log.Printf("用户 %d 开始抢红包\n", *user.UID)

	maxCount := c.GetInt(MaxCountField)
	log.Printf("成功获取最大的可抢红包限额: %d\n", maxCount)
	// 获取当前用户已抢红包的数量
	curCount, err := Mapper.GetRedEnvelops(c, *user.UID)
	if err != nil {
		log.Printf("查询用户 %d 已抢红包数失败\n", *user.UID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 查询用户已抢红包数失败",
			"data": err,
		})
		return
	}
	log.Printf("成功获取用户 %d 已抢红包数目: %d\n", *user.UID, curCount)
	// 判断用户是否超过红包数限额
	if curCount >= maxCount {
		log.Printf("用户 %d 已抢红包数达到限额\n", *user.UID)
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "error, 用户已抢红包数达到限额",
			"data": nil,
		})
		return
	}

	// 获取系统已发红包总数
	numberOfEnvelopesForALlUser, err := Mapper.GetNumberOfEnvelopesForALlUser(c)
	if err != nil {
		log.Printf("查询系统已发红包总数失败\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 查询用户已发红包总数失败",
			"data": err,
		})
		return
	}
	log.Printf("成功获取系统已发红包总数: %d\n", numberOfEnvelopesForALlUser)
	// 判断系统是否超过红包数限额
	if numberOfEnvelopesForALlUser >= c.GetInt(TotalNumberField) {
		log.Printf("系统已发红包总数达到限额\n")
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "error, 系统已发红包总数达到限额",
			"data": nil,
		})
		return
	}

	// 用户有一定概率抢不到红包
	probability := c.GetFloat64(ProbabilityField)
	// 生成一个0到1的随机数
	random := rand.Float64()
	if random > probability {
		// 没抢到红包
		log.Printf("用户 %d 没抢到红包\n", *user.UID)
		c.JSON(http.StatusOK, gin.H{
            "code": 500,
            "msg":  "error, 用户没抢到红包",
            "data": nil,
        })
		return
	}

	// 尝试增加已发红包数
	numberOfEnvelopesForAllUser, err := Mapper.IncreaseNumberOfEnvelopesForAllUser(c)
	if err != nil {
		log.Printf("增加系统已发红包总数失败\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 增加系统已发红包总数失败",
			"data": err,
		})
		return
	}
	log.Printf("成功增加系统已发红包总数: %d\n", numberOfEnvelopesForAllUser)
	// 判断增加之后是否超额
	if numberOfEnvelopesForAllUser >= c.GetInt(TotalNumberField) {
		log.Printf("系统已发红包总数达到限额\n")
		// 递减刚刚增加的红包
		err := Mapper.DecreaseOpenedEnvelopes(c)
		if err != nil {
			log.Printf("撤销对系统已发红包总数的自增失败")
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "error, 系统已发红包总数达到限额",
			"data": nil,
		})
		return
	}

	// 尝试增加已抢红包数
	log.Printf("尝试增加用户 %d 已抢红包个数\n", *user.UID)
	curCount, err = Mapper.IncreaseRedEnvelopes(c, *user.UID)
	if err != nil {
		log.Printf("增加用户 %d 已抢红包个数失败\n", *user.UID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 增加已抢红包数失败",
			"data": err,
		})
		return
	}

	// 可能因为并发抢红包的情况，导致用户已抢的红包数超过限额，这时候需要减少已抢红包数（否则配置更新将会出错）
	if curCount > maxCount {
		log.Printf("增加完已抢红包数，用户 %d 已抢红包数超过限额，尝试取消上一步操作\n", *user.UID)
		err = Mapper.DecreaseRedEnvelopes(c, *user.UID)
		if err != nil {
			log.Printf("撤销 增加用户 %d 已抢红包数失败\n", *user.UID)
		}
		err = Mapper.DecreaseNumberOfEnvelopesForAllUser(c)
		if err != nil {
			log.Printf("撤销 增加已抢红包总数失败\n")
		}
		log.Printf("取消用户 %d 已抢红包数成功\n", *user.UID)
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "error, 用户已抢红包数超过限额",
			"data": nil,
		})
		return
	}

	log.Printf("成功增加了用户 %d 已抢红包数量，准备生成红包id并添加到set中\n", *user.UID)
	// 成功增加了已抢红包数量，生成红包id并添加到set中

	// 生成新红包的id
	envelopeID, err := Mapper.IncreaseCurEnvelopeId(c)
	if err != nil {
		log.Printf("生成新红包id失败\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 生成新红包id失败",
			"data": err,
		})
		return
	}
	log.Printf("成功为用户 %d 生成红包 %d\n", *user.UID, envelopeID)
	err = Mapper.AddRedEnvelopeToUserId(c, *user.UID, envelopeID)
	if err != nil {
		log.Printf("为用户 %d 添加红包 %d 失败\n", *user.UID, envelopeID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 添加用户红包id失败",
			"data": err,
		})
		return
	}
	log.Printf("成功为用户 %d 添加红包 %d\n", *user.UID, envelopeID)

	// 将红包、用户信息写入MQ
	err = SnatchHistoryToMQ(*user.UID, envelopeID)
	if err != nil {
		// 回滚操作，丢弃请求。
		log.Println("MQ not working... Rollback & Return")
		// 撤销上面的redis操作
		err := Mapper.RemoveRedEnvelopeForUser(c, *user.UID, envelopeID)
		if err != nil {
			log.Printf("删除用户 %d 的红包 %d 失败\n", *user.UID, envelopeID)
		}
		err = Mapper.DecreaseRedEnvelopes(c, *user.UID)
		if err != nil {
			log.Printf("减少用户 %d 抢到的红包数\n", *user.UID)
		}
		err = Mapper.DecreaseNumberOfEnvelopesForAllUser(c)
		if err != nil {
			log.Printf("减少已发放的红包总数失败\n")
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 向消息队列发送消息失败",
			"data": err,
		})
		return
	}

	data := SuccessSnatch{envelopeID, maxCount, curCount}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
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
		log.Printf("查询用户 %d 是否拥有红包 %d 失败\n", *openre.UID, *openre.EnvelopeID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 查询用户是否拥有红包失败",
			"data": err,
		})
		return
	}
	if !owned {
		log.Printf("用户 %d 没有拥有红包 %d\n", *openre.UID, *openre.EnvelopeID)
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "error, 用户未拥有该红包或该红包已被拆开",
			"data": nil,
		})
		return
	}

	// 用户拥有该红包，尝试拆红包
	err = Mapper.RemoveRedEnvelopeForUser(c, *openre.UID, *openre.EnvelopeID)
	if err != nil {
		log.Printf("为用户 %d 拆红包 %d 失败\n", *openre.UID, *openre.EnvelopeID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 拆红包失败",
			"data": err,
		})
		return
	}

	// 生成红包的金额
	openedEnvelopes, err := Mapper.GetOpenedEnvelopes(c)
	if err != nil {
		log.Printf("获取已拆开红包个数失败\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 获取已拆开红包个数失败",
			"data": err,
		})
		return
	}
	spentBudget, err := Mapper.GetSpentBudget(c)
	if err != nil {
		log.Printf("获取已花费预算失败\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 获取已花费预算失败",
			"data": err,
		})
		return
	}
	money := GenerateRedEnvelopeValue(c.GetInt(BudgetField)-spentBudget,
		c.GetInt(TotalNumberField)-openedEnvelopes, c.GetInt(MaxValueField), c.GetInt(MinValueField))
	_, err = Mapper.IncreaseSpentBudget(c, money)
	if err != nil {
		log.Printf("增加已花费预算失败\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 增加已花费预算失败",
			"data": err,
		})
		return
	}
	_, err = Mapper.IncreaseOpenedEnvelopes(c)
	if err != nil {
		log.Printf("增加已拆红包数失败\n")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 增加已拆红包数失败",
			"data": err,
		})
		return
	}

	// 将红包id、红包金额写入MQ
	err = OpenValueToMQ(*openre.UID, money)
	if err != nil {
		// 回滚操作，丢弃请求。
		log.Println("MQ not working... Rollback & Return")
		// 撤销上面的redis操作
		err := Mapper.DecreaseOpenedEnvelopes(c)
		if err != nil {
			log.Printf("减少已拆红包数失败\n")
		}
		err = Mapper.DecreaseSpentBudget(c, money)
		if err != nil {
			log.Printf("减少已花费预算失败\n")
		}
		err = Mapper.AddRedEnvelopeToUserId(c, *openre.UID, *openre.EnvelopeID)
		if err != nil {
			log.Printf("将红包 %d 放入用户 %d 的红包集合失败\n", *openre.EnvelopeID, *openre.UID)
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 向消息队列发送消息失败",
			"data": err,
		})
		return
	}

	data := SuccessOpen{money}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
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
}

// SetConfig 设置红包全局配置
func SetRedEnvelopeConfig(c *gin.Context) {
	var config Config
	//匹配参数
	if err := c.ShouldBind(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "error, 输入参数有误",
			"data": err,
		})
		return
	}
	var configMap map[string]interface{}
	if config.MaxCount != nil {
		c.Set(MaxCountField, *config.MaxCount)
		configMap[MaxCountField] = *config.MaxCount
	}
	if config.Probability != nil {
		c.Set(ProbabilityField, *config.Probability)
		configMap[ProbabilityField] = *config.Probability
	}
	if config.Budget != nil {
		c.Set(BudgetField, *config.Budget)
		configMap[BudgetField] = *config.Budget
	}
	if config.TotalNumber != nil {
		c.Set(TotalNumberField, *config.TotalNumber)
		configMap[TotalNumberField] = *config.TotalNumber
	}
	if config.MinValue != nil {
		c.Set(MinValueField, *config.MinValue)
		configMap[MinValueField] = *config.MinValue
	}
	if config.MaxValue != nil {
		c.Set(MaxValueField, *config.MaxValue)
		configMap[MaxValueField] = *config.MaxValue
	}
	err := Mapper.SetConfigParameters(c, configMap)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "设置红包全局配置失败",
		})
		return
	}
}

// WrkTest 压力测试
func WrkTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}
