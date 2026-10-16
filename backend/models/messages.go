package models

import (
	"time"
)


type Message struct {
	ID			uint `gorm:"primaryKey"`
	SenderID	uint
	ChatID		uint
	Msg			string
	Time		time.Time
}

func (Message) TableName() string {
	return "messages"
}