package Chat

//
//import (
//	"encoding/json"
//	"log"
//)
//
//const SendMessageAction = "send-message"
//const JoinRoomAction = "join-room"
//const LeaveRoomAction = "leave-room"
//
//type Message struct {
//	Action  string  `json:"action"`
//	Message string  `json:"message"`
//	Sender  *Client `json:"sender"`
//}
//
//func (message *Message) encode() []byte {
//	json, err := json.Marshal(message)
//	if err != nil {
//		log.Println(err)
//	}
//
//	return json
//}
