package Chat

import (
	"BackendSimple/auth"
	"BackendSimple/db"
	"encoding/json"
	"log"
	"net/http"
)

type Chat struct {
	Auth auth.Auth
	DB   db.DB

	Rooms map[string]*Room
}

func Drop400Error(w http.ResponseWriter) {
	marshal, err := json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{400, "Bad request"})
	if err != nil {
		return
	}
	w.WriteHeader(400)
	_, err = w.Write(marshal)
	if err != nil {
		return
	}
}
func Drop500Error(w http.ResponseWriter) {
	marshal, err := json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{500, "Internal server error"})
	if err != nil {
		return
	}
	w.WriteHeader(500)
	_, err = w.Write(marshal)
	if err != nil {
		return
	}
}
func Drop401Error(w http.ResponseWriter) {
	marshal, err := json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{401, "Unauthorized"})
	if err != nil {
		return
	}
	w.WriteHeader(401)
	_, err = w.Write(marshal)
	if err != nil {
		return
	}
}

func (c *Chat) ServeWs(w http.ResponseWriter, r *http.Request) {
	token, ok := r.URL.Query()["token"]

	if !ok || len(token) == 0 {
		Drop400Error(w)
		return
	}

	roomID, ok := r.URL.Query()["room"]

	if !ok || len(roomID) == 0 {
		Drop400Error(w)
		return
	}

	claims, err := c.Auth.GetTokenClaims(token[0])
	if err != nil {
		Drop401Error(w)
		log.Println(err)
		return
	}

	idString, ok := claims["sub"]
	if !ok {
		Drop400Error(w)
		return
	}

	id := uint(idString.(float64))
	var user db.User
	user.ID = id

	err = c.DB.GetUser(&user)
	if err != nil {
		Drop500Error(w)
		return
	}

	room, ok := c.Rooms[roomID[0]]
	if !ok {
		Drop500Error(w)
		return
	}

	l := db.Lesson{ChatID: roomID[0]}
	err = c.DB.GetLesson(&l)
	if err != nil {
		Drop401Error(w)
		log.Println(err)
		return
	}
	if l.TimeStarted == 0 {
		err = c.DB.StartLesson(&l)
		if err != nil {
			Drop500Error(w)
			log.Println(err)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Drop500Error(w)
		log.Println(err)
		return
	}

	client := newClient(conn, user, room, l)
	client.ExitFunc = c.GetExitFunc(roomID[0])
	room.register <- client
}

func (c *Chat) GetExitFunc(roomid string) func() error {
	return func() error {
		var l db.Lesson
		l.ChatID = roomid
		err := c.DB.EndLesson(&l)
		if err != nil {
			return err
		}
		return nil
	}
}

func NewChat(db db.DB, auth auth.Auth) Chat {
	return Chat{
		Auth: auth,
		DB:   db,

		Rooms: make(map[string]*Room),
	}
}

func (c *Chat) CreateRoom(name string, users ...uint) *Room {
	room := NewRoom(name, users)
	go room.RunRoom()

	c.Rooms[room.ID] = &room

	return &room
}
