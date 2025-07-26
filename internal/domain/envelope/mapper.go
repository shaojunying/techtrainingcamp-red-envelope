package envelope

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"red_envelope/internal/infrastructure/database"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var Mapper = newMapper()

func newMapper() *mapper {
	return &mapper{}
}

type mapper struct{}

//GetRedEnvelops ByUserId 获取id为userId的用户已抢到的红包数目
func (*mapper) GetRedEnvelops(ctx context.Context, userId int) (int, error) {
	key := fmt.Sprintf(NumberOfRedEnvelopePerUserKey, userId)
	rdx := database.GetRdx()
	value, err := rdx.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		// 获取红包数量失败
		return -1, errors.New("get redis key error: " + err.Error())
	}
	if err == redis.Nil {
		// 当前用户第一次抢红包
		return 0, nil
	}
	redEnvelopes, err := strconv.Atoi(value)
	if err != nil {
		return -1, errors.New("failed to parse red envelopes for user, err: " + err.Error())
	}
	return redEnvelopes, nil
}

//IncreaseRedEnvelopes 将id为userId的用户已抢到的红包数目加1, 并返回新的红包数目
func (*mapper) IncreaseRedEnvelopes(ctx context.Context, userId int) (int, error) {
	key := fmt.Sprintf(NumberOfRedEnvelopePerUserKey, userId)
	rdx := database.GetRdx()
	result, err := rdx.Incr(ctx, key).Result()
	if err != nil {
		return -1, errors.New("failed to increase red envelopes for user, err: " + err.Error())
	}
	return int(result), nil
}

//DecreaseRedEnvelopes 将id为userId的用户已抢到的红包数目减1, 并返回新的红包数目
func (m *mapper) DecreaseRedEnvelopes(ctx context.Context, userId int) error {
	key := fmt.Sprintf(NumberOfRedEnvelopePerUserKey, userId)
	rdx := database.GetRdx()
	_, err := rdx.Decr(ctx, key).Result()
	if err != nil {
		return errors.New("failed to decrease red envelopes for user, err: " + err.Error())
	}
	return nil
}

//AddRedEnvelopeToUserId 添加红包到id为userId的用户的红包列表中
func (*mapper) AddRedEnvelopeToUserId(ctx context.Context, userId, redEnvelopeId int) error {
	key := fmt.Sprintf(SetOfRedEnvelopePerUserKey, userId)
	rdx := database.GetRdx()
	err := rdx.SAdd(ctx, key, redEnvelopeId).Err()
	return err
}

//CheckIfOwnRedEnvelope 判断用户是否拥有id为redEnvelopeId的红包
func (*mapper) CheckIfOwnRedEnvelope(ctx context.Context, userId, redEnvelopeId int) (bool, error) {
	key := fmt.Sprintf(SetOfRedEnvelopePerUserKey, userId)
	rdx := database.GetRdx()
	result, err := rdx.SIsMember(ctx, key, redEnvelopeId).Result()
	if err != nil {
		return false, errors.New("failed to check if user own red envelope, err: " + err.Error())
	}
	return result, nil
}

//RemoveRedEnvelopeForUser 将红包从用户的红包列表中移除
func (*mapper) RemoveRedEnvelopeForUser(ctx context.Context, userId, redEnvelopeId int) (bool, error) {
	key := fmt.Sprintf(SetOfRedEnvelopePerUserKey, userId)
	rdx := database.GetRdx()
	result, err := rdx.SRem(ctx, key, redEnvelopeId).Result()
	if err != nil {
		return false, errors.New("failed to remove red envelope for user, err: " + err.Error())
    }
	return int(result) != 0, err
}

//IncreaseCurEnvelopeId 自增目前最大的红包id，并返回自增之后的结果
func (*mapper) IncreaseCurEnvelopeId(ctx context.Context) (int, error) {
	key := CurEnvelopeIdKey
	rdx := database.GetRdx()
	result, err := rdx.Incr(ctx, key).Result()
	if err != nil {
		return -1, errors.New("failed to generate new red envelope id, err: " + err.Error())
	}
	return int(result), nil
}

//IncreaseNumberOfEnvelopesForAllUser 增加已发红包总数
func (*mapper) IncreaseNumberOfEnvelopesForAllUser(ctx context.Context) (int, error) {
	key := NumberOfEnvelopesForAllUserKey
	rdx := database.GetRdx()
	result, err := rdx.Incr(ctx, key).Result()
	if err != nil {
		return -1, errors.New("failed to increase number of envelopes, err: " + err.Error())
	}
	return int(result), nil
}

//GetNumberOfEnvelopesForALlUser 获取已发红包总数
func (*mapper) GetNumberOfEnvelopesForALlUser(ctx context.Context) (int, error) {
	key := NumberOfEnvelopesForAllUserKey
	rdx := database.GetRdx()
	result, err := rdx.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		// 获取出错
		return -1, errors.New("failed to get number of envelopes, err: " + err.Error())
	}

	if err == redis.Nil {
		// 初始值为0
		return 0, nil
	}

	numberOfEnvelopes, err := strconv.Atoi(result)
	if err != nil {
		return -1, errors.New("failed to parse number of envelopes, err: " + err.Error())
	}
	return numberOfEnvelopes, nil
}

//DecreaseNumberOfEnvelopesForAllUser 减少已发红包总数
func (*mapper) DecreaseNumberOfEnvelopesForAllUser(ctx context.Context) error {
	key := NumberOfEnvelopesForAllUserKey
	rdx := database.GetRdx()
	err := rdx.Decr(ctx, key).Err()
	return err
}

//GetOpenedEnvelopes 获取已拆开的红包数量
func (*mapper) GetOpenedEnvelopes(ctx context.Context) (int, error) {
	value, err := database.GetRdx().Get(ctx, OpenedEnvelopesKey).Result()
	if err != nil && err != redis.Nil {
		return -1, errors.New("failed to get opened envelopes, err: " + err.Error())
	}
	if err == redis.Nil {
		return 0, nil
	}
	openedEnvelopes, err := strconv.Atoi(value)
	if err != nil {
		return -1, errors.New("failed to parse opened envelopes, err: " + err.Error())
	}
	return openedEnvelopes, nil
}

func (*mapper) IncreaseOpenedEnvelopes(ctx context.Context) (int, error) {
	result, err := database.GetRdx().Incr(ctx, OpenedEnvelopesKey).Result()
	return int(result), err
}

func (*mapper) DecreaseOpenedEnvelopes(ctx context.Context) error {
	return database.GetRdx().Decr(ctx, OpenedEnvelopesKey).Err()
}

func (*mapper) GetSpentBudget(ctx context.Context) (int, error) {
	value, err := database.GetRdx().Get(ctx, SpentBudgetKey).Result()
	if err != nil && err != redis.Nil {
		return -1, errors.New("failed to get spent budget, err: " + err.Error())
	}
	if err == redis.Nil {
		return 0, nil
	}
	spentBudget, err := strconv.Atoi(value)
	if err != nil {
		return -1, errors.New("failed to parse spent budget, err: " + err.Error())
	}
	return spentBudget, nil
}

func (*mapper) IncreaseSpentBudget(ctx context.Context, amount int) (int, error) {
	result, err := database.GetRdx().IncrBy(ctx, SpentBudgetKey, int64(amount)).Result()
	return int(result), err
}
func (*mapper) DecreaseSpentBudget(ctx context.Context, amount int) error {
	return database.GetRdx().DecrBy(ctx, SpentBudgetKey, int64(amount)).Err()
}

//GetConfigParameters 获取配置参数
func (*mapper) GetConfigParameters(ctx context.Context) (*Config, error) {
	rdx := database.GetRdx()
	configMap, err := rdx.HGetAll(ctx, ConfigKey).Result()
	if err != nil {
		return nil, err
	}
	config := Config{}
	maxCount, err := strconv.Atoi(configMap[MaxCountField])
	if err != nil {
		return nil, err
	}
	config.MaxCount = &maxCount
	probability, err := strconv.ParseFloat(configMap[ProbabilityField], 64)
	if err != nil {
		return nil, err
	}
	config.Probability = &probability
	budget, err := strconv.Atoi(configMap[BudgetField])
	if err != nil {
		return nil, err
	}
	config.Budget = &budget
	totalNumber, err := strconv.Atoi(configMap[TotalNumberField])
	if err != nil {
		return nil, err
	}
	config.TotalNumber = &totalNumber
	minValue, err := strconv.Atoi(configMap[MinValueField])
	if err != nil {
		return nil, err
	}
	config.MinValue = &minValue
	maxValue, err := strconv.Atoi(configMap[MaxValueField])
	if err != nil {
		return nil, err
	}
	config.MaxValue = &maxValue
	return &config, nil
}

// SetConfigParameters 设置配置参数
func (*mapper) SetConfigParameters(ctx context.Context, configMap map[string]interface{}) error {
	rdx := database.GetRdx()
	err := rdx.HSet(ctx, ConfigKey, configMap).Err()
	return err
}

func (m *mapper) GetLastRequestTime(ctx context.Context, uid int) (int64, error) {
	key := fmt.Sprintf(LastRequestTimeKey, uid)
	rdx := database.GetRdx()
	value, err := rdx.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return -1, err
	}
	if err == redis.Nil {
		return 0, nil
	}
	lastRequestTime, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
        return -1, err
    }
	return lastRequestTime, nil
}

func (m *mapper) UpdateLastRequestTime(c *gin.Context, uid int, unix int64, milliseconds int64) error {
	key := fmt.Sprintf(LastRequestTimeKey, uid)
	rdx := database.GetRdx()
	err := rdx.Set(c, key, unix, time.Millisecond*time.Duration(milliseconds)).Err()
	return err
}
