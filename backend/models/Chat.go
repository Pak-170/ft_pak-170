package models

import "encoding/json"

type Chat struct {
	ID				uint	`gorm:"primaryKey"`
	Name			string	`gorm:"name"`
	Owner			uint	`gorm:"owner" json:"-"`
	Participants	[]User	`gorm:"-" json:"Participants"`
}

func (Chat) TableName() string {
	return "chats"
}

func (c Chat) MarshalJSON() ([]byte, error) {
	type Alias Chat
	var chatType	string
	if c.Name != "" {
		chatType = "group"
	} else {
		chatType = "individual"
	}

	return json.Marshal(&struct{
		InfoType		string	`json:"infoType"`
		ChatType		string	`json:"Type"`
		*Alias
	}{
		InfoType: "chat",
		ChatType: chatType,
		Alias: (*Alias)(&c),
	})
}

func NewChat(id uint, name string, owner uint, participants []User) *Chat {
	return &Chat{
		ID: id,
		Name: name,
		Owner: owner,
		Participants: participants,
	}
}
