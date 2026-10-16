package models


type ChatUserRelation struct {
	ChatID	uint `gorm:"column:chat"`
	UserID	uint `gorm:"column:user"`
}

func (ChatUserRelation) TableName() string {
	return "chat_user"
}