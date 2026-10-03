package models

import "time"

type Message struct {
	ID			uint		`gorm:"primaryKey" json:"id"`
	SenderID	uint		`json:"sender_id"`
	ChatID		uint		`json:"chat_id"`
	Msg			string		`json:"msg"`
	SendTime	time.Time	`gorm:"autoCreateTime" json:"send_time"`
}

func (Message) TableName() string {
	return "messages"
}