package models

type UsersInChatView struct {
	ChatID		uint
	Name		string
	UserID		uint
	Username	string
	Alias		*string
}

func (UsersInChatView) TableName() string {
	return "user_chat_view"
}

