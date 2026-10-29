package models

import "encoding/json"

type User struct {
	ID				uint	`gorm:"primaryKey"`
	Username		string	`gorm:"username" json:"username"`
	Alias			*string	`gorm:"alias" json:"alias"`
	PasswordHash	string	`gorm:"password_hash" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		InfoType		string	`json:"infoType"`
		*Alias
	}{
		InfoType: "chat",
		Alias: (*Alias)(&u),
	})
}