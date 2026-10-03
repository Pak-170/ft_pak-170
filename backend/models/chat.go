package models

type Chat struct {
	ID		uint	`gorm:"primaryKey" json:"id"`
	Name	string	`json:"name"`
	OwnerID	uint	`json:"owner_id"`
}

func (Chat) TableName() string {
	return "chats"
}

// Intermediate table: which users belong to which chat
type ChatUser struct {
	ChatID uint `gorm:"primaryKey"`
	UserID uint `gorm:"primaryKey"`
}

func (ChatUser) TableName() string {
	return "chat_user"
}