package snatch

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
	IfOpen     bool      `json:"if_open"`
	Money      int       `json:"money"`
	OpenTime   time.Time `gorm:"column:open_time;default:null"`
}
