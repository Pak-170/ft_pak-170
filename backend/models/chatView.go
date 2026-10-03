package models

import "time"

type ChatView struct {
	ID			uint
	ChatID		uint
	Name		string
	Username	string
	Msg			string
	SendTime	time.Time
}

func (ChatView) TableName() string {
	return "chat_view"
}

