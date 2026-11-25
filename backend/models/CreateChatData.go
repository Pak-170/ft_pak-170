package models

type CreateChatData struct {
	InfoType	string	`json:"infoType"`
	Name		string	`json:"name"`
	ChatType	string	`json:"type"`
	UserIDs		[]uint	`json:"userIDs"`
}
