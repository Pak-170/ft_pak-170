package models

import "time"

type Session struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement:false"`
	UserID    uint
	ExpiresAt time.Time
}

func (Session) TableName() string {
	return "sessions"
}
