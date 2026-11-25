package validation

import (
	"go_pakito/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"

)

var v *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func IsValidMsg(msg *models.ChatView, err error) (int, string){
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			return websocket.CloseNormalClosure, "Bye!"
		} else {
			return websocket.CloseInvalidFramePayloadData, err.Error()
		}
	}
	if v.Struct(msg) != nil {
		return websocket.ClosePolicyViolation, "Invalid JSON"
	}
	return 0, ""
}