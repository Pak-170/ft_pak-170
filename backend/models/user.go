package models

type User struct {
	ID			uint	`gorm:"primaryKey"`
	Username	string	`json:"username"`
	Alias		*string	`json:"alias"`
	PasswordHash	string	`json:"-"`
}

func (User) TableName() string {
	return "users"
}