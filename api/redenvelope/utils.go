package redenvelope

//ConfigKey redis中配置参数的key
var ConfigKey = "config"

// NumberOfRedEnvelopePerUserKey redis中当前用户已抢红包数
var NumberOfRedEnvelopePerUserKey = "num_%d"

//SetOfRedEnvelopePerUserKey redis中当前用户已抢红包集合
var SetOfRedEnvelopePerUserKey = "envelopes_%d"

//NumberOfEnvelopesForAllUserKey redis中所有用户已抢红包数
var NumberOfEnvelopesForAllUserKey = "number_of_envelopes"

//CurEnvelopeIdKey 目前已生成的最大红包id
var CurEnvelopeIdKey = "cur_envelope_id"

//OpenedEnvelopesKey 目前已拆开的红包数目
var OpenedEnvelopesKey = "opened_envelopes"

//SpentBudgetKey 目前已花费的预算
var SpentBudgetKey = "spent_budget"

// LastRequestTimeKey 当前用户上一次请求的时间戳
var LastRequestTimeKey = "last_request_time_%d"

//MaxCountField 每个用户最多可抢到的次数
var MaxCountField = "max_count"

// ProbabilityField 每次抢到红包的概率
var ProbabilityField = "probability"

// BudgetField 总预算（以分为单位）
var BudgetField = "budget"

// TotalNumberField 总红包数量
var TotalNumberField = "total_number"

// MinValueField 每个红包的最小金额（以分为单位）
var MinValueField = "min_value"

// MaxValueField 每个红包的最大金额（以分为单位）
var MaxValueField = "max_value"
