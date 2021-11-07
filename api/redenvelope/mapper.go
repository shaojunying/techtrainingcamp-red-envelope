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
	key := fmt.Sprintf("num_${%d}", userId)
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

//GetRedEnvelopeLimitForUser 获取每位用户抢红包的限额
func (*mapper) GetRedEnvelopeLimitForUser(ctx context.Context) (int, error) {
	key := "max_count"
	rdx := database.GetRdx()
	value, err := rdx.Get(ctx, key).Result()
	if err != nil {
		return -1, errors.New("failed to get max count of red envelopes for every user, err: " + err.Error())
	}
	redEnvelopesLimit, err := strconv.Atoi(value)
	if err != nil {
		return -1, errors.New("failed to parse red envelopes for user, err: " + err.Error())
	}
	return redEnvelopesLimit, nil
}

//IncreaseRedEnvelopes 将id为userId的用户已抢到的红包数目加1, 并返回新的红包数目
func (*mapper) IncreaseRedEnvelopes(ctx context.Context, userId int) (int, error) {
	key := fmt.Sprintf("num_${%d}", userId)
	rdx := database.GetRdx()
	result, err := rdx.Incr(ctx, key).Result()
	if err != nil {
		return -1, errors.New("failed to increase red envelopes for user, err: " + err.Error())
	}
	return int(result), nil
}

//AddRedEnvelopeToUserId 添加红包到id为userId的用户的红包列表中
func (*mapper) AddRedEnvelopeToUserId(ctx context.Context, userId, redEnvelopeId int) error {
	key := fmt.Sprintf("envelopes_${%d}", userId)
	rdx := database.GetRdx()
	err := rdx.SAdd(ctx, key, redEnvelopeId).Err()
	return err
}

//CheckIfOwnRedEnvelope 判断用户是否拥有id为redEnvelopeId的红包
func (*mapper) CheckIfOwnRedEnvelope(ctx context.Context, userId, redEnvelopeId int) (bool, error) {
    key := fmt.Sprintf("envelopes_${%d}", userId)
    rdx := database.GetRdx()
    result, err := rdx.SIsMember(ctx, key, redEnvelopeId).Result()
    if err != nil {
        return false, errors.New("failed to check if user own red envelope, err: " + err.Error())
    }
    return result, nil
}

//OpenRedEnvelope 如果id为redEnvelopeId的红包没有被拆开, 则拆开红包, 并返回true；否则返回false
func (*mapper) OpenRedEnvelope(ctx context.Context, redEnvelopeId int, value int) (bool, error) {
	key := fmt.Sprintf("opened_${%d}", redEnvelopeId)
	rdx := database.GetRdx()
	return rdx.SetNX(ctx, key, value, 0).Result()
}
