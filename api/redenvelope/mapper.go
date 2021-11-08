package redenvelope

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"red_envelope/database"
	"strconv"
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
func (*mapper) RemoveRedEnvelopeForUser(ctx context.Context, userId, redEnvelopeId int) error {
	key := fmt.Sprintf(SetOfRedEnvelopePerUserKey, userId)
	rdx := database.GetRdx()
	err := rdx.SRem(ctx, key, redEnvelopeId).Err()
	return err
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

//GenerateNewRedEnvelopeId 生成新的红包id
func (*mapper) GenerateNewRedEnvelopeId(ctx context.Context) (int, error) {
	key := NumberOfEnvelopesForAllUserKey
	rdx := database.GetRdx()
	result, err := rdx.Incr(ctx, key).Result()
	if err != nil {
		return -1, errors.New("failed to generate new red envelope id, err: " + err.Error())
	}
	return int(result), nil
}

//GetConfigParameters 获取配置参数
func (*mapper) GetConfigParameters(ctx context.Context) (*Config, error) {
	rdx := database.GetRdx()
	configMap, err := rdx.HGetAll(ctx, ConfigKey).Result()
	fmt.Println(configMap)
	if err != nil {
		return nil, err
	}
	config := Config{}
	maxCount, err := strconv.Atoi(configMap[MaxCountField])
	if err != nil {
		return nil, err
	}
	config.MaxCount = maxCount
	probability, err := strconv.ParseFloat(configMap[ProbabilityField], 64)
	if err != nil {
        return nil, err
    }
	config.Probability = probability
	budget, err := strconv.Atoi(configMap[BudgetField])
	if err != nil {
		return nil, err
	}
	config.Budget = budget
	totalNumber , err := strconv.Atoi(configMap[TotalNumberField])
	if err != nil {
        return nil, err
    }
	config.TotalNumber = totalNumber
	minValue, err := strconv.Atoi(configMap[MinValueField])
	if err != nil {
        return nil, err
    }
	config.MinValue = minValue
	maxValue, err := strconv.Atoi(configMap[MaxValueField])
	if err != nil {
        return nil, err
    }
	config.MaxValue = maxValue
	return &config, nil
}
