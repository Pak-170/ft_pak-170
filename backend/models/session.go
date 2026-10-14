package models

import "time"

type Session struct {
	ID			string `gorm:"primaryKey"`
	UserID		uint
	ExpiresAt	time.Time
}

func (Session) TableName () string {
	return "sessions"
}
