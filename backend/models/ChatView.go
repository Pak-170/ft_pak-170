package models

import "time"

type ChatView struct {
	ChatID		uint		`gorm:"column:chat_id"`
	Username	string		`gorm:"column:username"`
	Msg			string		`gorm:"column:msg"`
	SendTime	time.Time	`gorm:"column:time"`
}

func (ChatView) TableName() string {
	return "chat_view"
}

