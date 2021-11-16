package redenvelope

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

func HandleSnatchOK(c *gin.Context, code int, uid int, data *SuccessSnatch) {
	var msg string
	switch code {
	case 0:
		msg = "uid正确，成功抢到红包"
	case 1:
		msg = "uid正确，没能抢到红包"
	case 2:
		msg = "uid正确，用户抢红包次数到限额"
	case 3:
		msg = "红包已全部发完"
	}
	log.Printf("%v的用户正常响应：%v\n", uid, msg)
	if data == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": *data,
	})
}

func HandleOpenOK(c *gin.Context, code int, uid int, envelopeid int, data *SuccessOpen) {
	var msg string
	switch code {
	case 0:
		msg = "成功拆开红包"
	case 4:
		msg = "不能拆未拥有的或者已拆开的红包"
	}
	log.Printf("%v的用户拆%v的红包正常响应：%v\n", uid, envelopeid, msg)
	if data == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": *data,
	})
}

func HandleERR(c *gin.Context, code int, err error) {
	var msg string
	httpcode := http.StatusInternalServerError
	switch code {
	case 101, 102:
		msg = "请求参数有误"
		httpcode = http.StatusBadRequest
	case 103:
		msg = "请求过于频繁"
		httpcode = http.StatusBadRequest
	case 201, 202:
		msg = "查询时有误"
	case 301, 302, 303:
		msg = "写入时有误"
	case 401, 402:
		msg = "其他错误"
	}
	log.Printf("%v，错误码：%v，具体错误为：%v\n", msg, code, err)
	c.JSON(httpcode, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}

// SnatchRedEnvelope 抢红包
func SnatchRedEnvelope(c *gin.Context) {
	var r RedEnvelope
	//匹配参数
	err := c.ShouldBindBodyWith(&r, binding.JSON)
	if err != nil {
		HandleERR(c, 101, err)
		return
	}

	if r.UID == nil {
		HandleERR(c, 102, errors.New("user.UID is nil"))
		return
	}
	log.Printf("用户 %d 开始抢红包\n", *r.UID)

	// 只有一定概率抢到红包
	rand.Seed(time.Now().UnixNano())
	num := rand.Float64()                         //随机数
	probability := c.GetFloat64(ProbabilityField) //抢到的概率值 %
	log.Printf("用户 %d 抽到了随机数 %f，小于 %f 可以抢到红包", *r.UID, num, probability)
	if num > probability {
		HandleSnatchOK(c, 1, *r.UID, nil)
		return
	}

	maxCount := c.GetInt(MaxCountField)
	log.Printf("成功获取最大的可抢红包限额: %d\n", maxCount)


	out1 := make(chan Result)
	out2 := make(chan Result)
	go func() {
		numberOfEnvelopesForALlUser, err := Mapper.GetNumberOfEnvelopesForALlUser(c)
		result := Result{val: numberOfEnvelopesForALlUser, err: err}
		out1 <- result
	}()
	go func() {
		redEnvelops, err := Mapper.GetRedEnvelops(c, *r.UID)
		result := Result{val: redEnvelops, err: err}
		out2 <- result
	}()
	result := <-out1
	if result.err != nil {
		HandleERR(c, 202, result.err)
		return
	}
	numberOfEnvelopesForALlUser := result.val
	log.Printf("成功获取系统已发红包总数: %d\n", numberOfEnvelopesForALlUser)

	result = <-out2
	if result.err != nil {
		HandleERR(c, 201, result.err)
		return
	}
	curCount := result.val
	log.Printf("成功获取用户 %d 已抢红包数目: %d\n", *r.UID, curCount)

	// 判断用户是否超过红包数限额
	if curCount >= maxCount {
		HandleSnatchOK(c, 2, *r.UID, nil)
		return
	}

	// 判断系统是否超过红包数限额
	if numberOfEnvelopesForALlUser >= c.GetInt(TotalNumberField) {
		HandleSnatchOK(c, 3, *r.UID, nil)
		return
	}


	out3 := make(chan Result)
	out4 := make(chan Result)
	out5 := make(chan Result)
	go func() {
		// 尝试增加已发红包数
		numberOfEnvelopesForAllUser, err := Mapper.IncreaseNumberOfEnvelopesForAllUser(c)
		result := Result{val: numberOfEnvelopesForAllUser, err: err}
		out3 <- result
	}()
	go func() {
		// 尝试增加已抢红包数
		curCount, err = Mapper.IncreaseRedEnvelopes(c, *r.UID)
		result := Result{val: curCount, err: err}
		out4 <- result
	}()
	go func() {
		// 生成新红包的id
		envelopeID, err := Mapper.IncreaseCurEnvelopeId(c)
		result := Result{val: envelopeID, err: err}
		out5 <- result
	}()

	result = <- out3
	if result.err != nil {
		HandleERR(c, 301, result.err)
		return
	}
	numberOfEnvelopesForAllUser := result.val
	log.Printf("成功增加系统已发红包总数: %d\n", numberOfEnvelopesForAllUser)

	result = <- out4
	if result.err != nil {
		HandleERR(c, 302, result.err)
		return
	}
	curCount = result.val

	result = <- out5
	if result.err != nil {
		HandleERR(c, 401, result.err)
		return
	}
	envelopeID := result.val
	log.Printf("成功为用户 %d 生成红包 %d\n", *r.UID, envelopeID)

	// 判断增加之后是否超额
	if numberOfEnvelopesForAllUser >= c.GetInt(TotalNumberField) {
		// TODO 逻辑错误
		// 递减刚刚增加的红包
		err := Mapper.DecreaseOpenedEnvelopes(c)
		if err != nil {
			log.Printf("撤销对系统已发红包总数的自增失败")
		}
		log.Printf("增加完已抢红包数，用户 %d 已抢红包数超过限额，尝试取消上一步操作\n", *r.UID)
		err = Mapper.DecreaseRedEnvelopes(c, *r.UID)
		if err != nil {
			log.Printf("撤销 增加用户 %d 已抢红包数失败\n", *r.UID)
		}
		err = Mapper.DecreaseNumberOfEnvelopesForAllUser(c)
		if err != nil {
			log.Printf("撤销 增加已抢红包总数失败\n")
		}
		log.Printf("取消用户 %d 已抢红包数成功\n", *r.UID)
		HandleSnatchOK(c, 3, *r.UID, nil)
		return
	}

	// 可能因为并发抢红包的情况，导致用户已抢的红包数超过限额，这时候需要减少已抢红包数（否则配置更新将会出错）
	if curCount > maxCount {
		log.Printf("增加完已抢红包数，用户 %d 已抢红包数超过限额，尝试取消上一步操作\n", *r.UID)
		err = Mapper.DecreaseRedEnvelopes(c, *r.UID)
		if err != nil {
			log.Printf("撤销 增加用户 %d 已抢红包数失败\n", *r.UID)
		}
		err = Mapper.DecreaseNumberOfEnvelopesForAllUser(c)
		if err != nil {
			log.Printf("撤销 增加已抢红包总数失败\n")
		}
		log.Printf("取消用户 %d 已抢红包数成功\n", *r.UID)
		HandleSnatchOK(c, 2, *r.UID, nil)
		return
	}

	// 成功增加了已抢红包数量，生成红包id并添加到set中
	log.Printf("成功增加了用户 %d 已抢红包数量，准备生成红包id并添加到set中\n", *r.UID)

	err = Mapper.AddRedEnvelopeToUserId(c, *r.UID, envelopeID)
	if err != nil {
		HandleERR(c, 303, err)
		return
	}
	log.Printf("成功为用户 %d 添加红包 %d\n", *r.UID, envelopeID)

	// 将红包、用户信息写入MQ
	err = SnatchHistoryToMQ(*r.UID, envelopeID)
	if err != nil {
		HandleERR(c, 402, err)
		// 回滚操作，丢弃请求。
		log.Println("MQ not working... Rollback & Return")
		// 撤销上面的redis操作
		_, err := Mapper.RemoveRedEnvelopeForUser(c, *r.UID, envelopeID)
		if err != nil {
			log.Printf("删除用户 %d 的红包 %d 失败\n", *r.UID, envelopeID)
		}
		err = Mapper.DecreaseRedEnvelopes(c, *r.UID)
		if err != nil {
			log.Printf("减少用户 %d 抢到的红包数\n", *r.UID)
		}
		err = Mapper.DecreaseNumberOfEnvelopesForAllUser(c)
		if err != nil {
			log.Printf("减少已发放的红包总数失败\n")
		}
		return
	}

	data := SuccessSnatch{envelopeID, maxCount, curCount}
	HandleSnatchOK(c, 0, *r.UID, &data)
}

// OpenRedEnvelope 拆红包
func OpenRedEnvelope(c *gin.Context) {
	var r RedEnvelope
	//匹配参数
	if err := c.ShouldBindBodyWith(&r, binding.JSON); err != nil {
		HandleERR(c, 101, err)
		return
	}
	if r.UID == nil || r.EnvelopeID == nil {
		HandleERR(c, 102, errors.New("UID or EnvelopeID is nil"))
		return
	}
	// 用户拥有该红包，尝试拆红包
	success, err := Mapper.RemoveRedEnvelopeForUser(c, *r.UID, *r.EnvelopeID)
	if err != nil {
		HandleERR(c, 304, err)
		return
	}
	if !success {
		// 用户没有这个红包 或者 红包已经被拆开
		HandleOpenOK(c, 4, *r.UID, *r.EnvelopeID, nil)
		return
	}

	// 生成红包的金额
	openedEnvelopes, err := Mapper.GetOpenedEnvelopes(c)
	if err != nil {
		HandleERR(c, 204, err)
		return
	}
	spentBudget, err := Mapper.GetSpentBudget(c)
	if err != nil {
		HandleERR(c, 205, err)
		return
	}
	money := GenerateRedEnvelopeValue(c.GetInt(BudgetField)-spentBudget,
		c.GetInt(TotalNumberField)-openedEnvelopes, c.GetInt(MaxValueField), c.GetInt(MinValueField))
	_, err = Mapper.IncreaseSpentBudget(c, money)
	if err != nil {
		HandleERR(c, 305, err)
		return
	}
	_, err = Mapper.IncreaseOpenedEnvelopes(c)
	if err != nil {
		HandleERR(c, 306, err)
		return
	}

	// 将红包id、红包金额写入MQ
	err = OpenValueToMQ(*r.EnvelopeID, money)
	if err != nil {
		HandleERR(c, 402, err)
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
		err = Mapper.AddRedEnvelopeToUserId(c, *r.UID, *r.EnvelopeID)
		if err != nil {
			log.Printf("将红包 %d 放入用户 %d 的红包集合失败\n", *r.EnvelopeID, *r.UID)
		}
		return
	}

	HandleOpenOK(c, 0, *r.UID, *r.EnvelopeID, &SuccessOpen{money})
}

// GetWalletList 钱包列表
func GetWalletList(c *gin.Context) {
	var r RedEnvelope
	//匹配参数
	if err := c.ShouldBindBodyWith(&r, binding.JSON); err != nil {
		HandleERR(c, 101, err)
		return
	}
	if r.UID == nil {
		HandleERR(c, 102, errors.New("user.UID is nil"))
		return
	}
	list, err := r.QueryListSql()
	if list != nil {
		log.Printf("%d获取到的红包列表有：\n", *r.UID)
		for _, pWalletList := range list["envelope_list"].([]*WalletList) {
			log.Printf("%+v\n", *pWalletList)
		}
	}

	if err != nil {
		HandleERR(c, 206, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": list,
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
	finalConfig, _ := Mapper.GetConfigParameters(c)
	configMap := make(map[string]interface{})
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
		finalConfig.Budget = config.Budget
	}
	if config.TotalNumber != nil {
		c.Set(TotalNumberField, *config.TotalNumber)
		configMap[TotalNumberField] = *config.TotalNumber
		finalConfig.TotalNumber = config.TotalNumber
	}
	if config.MinValue != nil {
		c.Set(MinValueField, *config.MinValue)
		configMap[MinValueField] = *config.MinValue
		finalConfig.MinValue = config.MinValue
	}
	if config.MaxValue != nil {
		c.Set(MaxValueField, *config.MaxValue)
		configMap[MaxValueField] = *config.MaxValue
	}
	if finalConfig.Budget != nil && finalConfig.TotalNumber != nil && finalConfig.MinValue != nil && *finalConfig.TotalNumber**finalConfig.MinValue > *finalConfig.Budget {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "预算不足以发出最小金额的红包",
		})
		return
	} else {
		if err := Mapper.SetConfigParameters(c, configMap); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "设置红包全局配置失败",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

func CatchKoi(c *gin.Context) {
	ks, err := QueryKoiListSql()
	if err != nil {
		HandleERR(c, 501, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "抓锦鲤成功",
		"data": ks,
	})
}

// WrkTest 压力测试
func WrkTest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

type Result struct {
	val int
	err error
}
