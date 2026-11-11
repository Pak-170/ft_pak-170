package models

import "encoding/json"

type User struct {
	ID				uint	`gorm:"primaryKey"`
	Username		string	`json:"username"`
	Alias			*string	`json:"alias"`
	Avatar			*string `json:"avatar"`
	PasswordHash	string	`json:"-"`
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