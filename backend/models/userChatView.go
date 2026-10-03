package models

type UsersInChatView struct {
	ChatID		int
	Name		string
	UserID		int
	Username	string
	Alias		*string
}

func (UsersInChatView) TableName() string {
	return "user_chat_view"
}

