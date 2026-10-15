package models

import "time"

type ChatView struct {
	Name		string
	Username	string
	Msg			string
	Send_time	time.Time
}

func (ChatView) TableName() string {
	return "chat_view"
}

