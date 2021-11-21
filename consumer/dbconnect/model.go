package dbconnect

//type User struct {
//	UID    *int `json:"uid"`
//	Amount int  `json:"amount"`
//	IfGet  bool `json:"if_get"`
//	//SnatchCount int    `json:"snatch_count"`
//}
//

type RedEnvelope struct {
	EnvelopeID *int `gorm:"primary_key" json:"envelope_id"`
	UID        *int `json:"uid"`
	Opened     bool `json:"opened"`
	Value      int  `json:"value"`
	SnatchTime int  `json:"snatch_time"` // 为了避免查询时需要遍历进行转换，直接使用时间戳对应的int类型
	//OpenTime   time.Time `gorm:"column:open_time;default:null"`
}

type SuccessSnatch struct {
	EnvelopeID int `json:"envelope_id"`
	MaxCount   int `json:"max_count"`
	CurCount   int `json:"cur_count"`
}

type SuccessOpen struct {
	Value int `json:"value"`
}

type WalletList struct {
	EnvelopeID int   `json:"envelope_id"`
	Opened     bool  `json:"opened"`
	Value      int   `json:"value,omitempty"`
	SnatchTime int64 `json:"snatch_time"`
}

type Config struct {
	//MaxCount 每个用户最多可抢到的次数
	MaxCount *int `json:"max_count"`

	// Probability 每次抢到红包的概率
	Probability *float64 `json:"probability"`

	// BudgetField 总预算（以分为单位）
	Budget *int `json:"budget"`

	// TotalNumber 总红包数量
	TotalNumber *int `json:"total_number"`

	// MinValue 每个红包的最小金额（以分为单位）
	MinValue *int `json:"min_value"`

	// MaxValue 每个红包的最大金额（以分为单位）
	MaxValue *int `json:"max_value"`
}
