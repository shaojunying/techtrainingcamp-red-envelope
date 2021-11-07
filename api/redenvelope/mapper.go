package redenvelope

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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

//RemoveRedEnvelopeForUser 将红包从用户的红包列表中移除
func (*mapper) RemoveRedEnvelopeForUser(ctx context.Context, userId, redEnvelopeId int) error {
	key := fmt.Sprintf("envelopes_${%d}", userId)
	rdx := database.GetRdx()
	err := rdx.SRem(ctx, key, redEnvelopeId).Err()
	return err
}

//DecreaseRedEnvelopes 将id为userId的用户已抢到的红包数目减1, 并返回新的红包数目
func (m *mapper) DecreaseRedEnvelopes(c *gin.Context, userId int) error {
	key := fmt.Sprintf("num_${%d}", userId)
	rdx := database.GetRdx()
	_, err := rdx.Decr(c, key).Result()
	if err != nil {
		return errors.New("failed to decrease red envelopes for user, err: " + err.Error())
	}
	return nil
}

//GenerateNewRedEnvelopeId 生成新的红包id
func (*mapper) GenerateNewRedEnvelopeId(c *gin.Context) (int, error) {
	key := "red_envelope_id"
	rdx := database.GetRdx()
	result, err := rdx.Incr(c, key).Result()
	if err != nil {
		return -1, errors.New("failed to generate new red envelope id, err: " + err.Error())
	}
	return int(result), nil
}
