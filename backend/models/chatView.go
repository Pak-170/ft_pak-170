package models

import "time"

type ChatView struct {
	ID			int
	ChatID		int
	Name		string
	Username	string
	Msg			string
	SendTime	time.Time
}

func (ChatView) TableName() string {
	return "chat_view"
}

