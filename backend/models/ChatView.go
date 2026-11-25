package models

import "time"

type ChatView struct {
	ChatID		uint		`validate:"required" gorm:"column:chat_id"`
	Username	string		`validate:"required" gorm:"column:username"`
	Msg			string		`validate:"required" gorm:"column:msg"`
	SendTime	time.Time	`validate:"required" gorm:"column:time"`
	SenderHashToken	string	`json:"-" gorm:"-"`
}

func (ChatView) TableName() string {
	return "chat_view"
}

