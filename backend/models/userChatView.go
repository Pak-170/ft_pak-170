package models

type UsersInChatView struct {
	ChatId		int
	Name		string
	User		int
	Username	string
	Alias		string
}

func (UsersInChatView) TableName() string {
	return "user_chat_view"
}

