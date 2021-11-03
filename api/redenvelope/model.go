package redenvelope

import "time"

type User struct {
	UID    *int `json:"uid"`
	Amount int  `json:"amount"`
	IfGet  bool `json:"if_get"`
	//SnatchCount int    `json:"snatch_count"`
}

type RedEnvelope struct {
	EnvelopeID int       `gorm:"primary_key" json:"envelope_id"`
	UID        int       `json:"uid"`
	Opened     bool      `json:"opened"`
	Value      int       `json:"value"`
	SnatchTime time.Time `gorm:"column:snatch_time;default:null"`
	//OpenTime   time.Time `gorm:"column:open_time;default:null"`
}

type OpenRE struct {
	UID        *int `json:"uid"`
	EnvelopeID *int `json:"envelope_id"`
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

type SuccessGet struct {
	Amount       int           `json:"amount"`
	EnvelopeList []*WalletList `json:"envelope_list"`
}
