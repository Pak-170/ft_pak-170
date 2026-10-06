package models

type Chat struct {
	ID			uint `gorm:"primaryKey"`
	Name		string
	Owner		uint
}

func (Chat) TableName() string {
	return "chats"
}