package models

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Log struct {
	Id       int       `json:"id,omitempty" db:"id"`
	Level    string    `json:"level,omitempty" db:"level"`         // Error/Success/Info
	MemberId string    `json:"member_id,omitempty" db:"member_id"` // member id
	Action   string    `json:"action,omitempty" db:"action"`       // func name
	Text     string    `json:"text,omitempty" db:"text"`
	Created  time.Time `json:"created,omitempty" db:"created"`
}

func ToLog(db *sqlx.DB, val *Log) {
	_, err := db.Exec("INSERT INTO logs (level, member_id, action, text, created) VALUES ($1, $2, $3, $4, $5)",
		val.Level, val.MemberId, val.Action, val.Text, time.Now())
	if err != nil {
		log.Println("ToLog()::", err)
	}
}
