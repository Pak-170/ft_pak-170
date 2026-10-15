package models

import "time"

type ChatView struct {
	chat_id		uint
	Username	string
	Msg			string
	Send_time	time.Time
}

func (ChatView) TableName() string {
	return "chat_view"
}

