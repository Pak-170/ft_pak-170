package models

type User struct {
	ID			uint `gorm:"primaryKey"`
	Username	string
	Alias		*string
}

func (User) TableName() string {
	return "users"
}