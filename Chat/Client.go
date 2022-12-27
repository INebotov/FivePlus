package Chat

import (
	"BackendSimple/db"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Client struct {
	conn     *websocket.Conn `json:"-"`
	send     chan []byte     `json:"-"`
	ID       uint            `json:"id"`
	LessonID uint            `json:"lesson_id"`
	Name     string          `json:"name"`
	room     *Room           `json:"-"`
	ExitFunc func() error    `json:"-"`

	ClientParams
}

type ClientParams struct {
	WriteWait      time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
	MaxMessageSize int64
}

func newClient(conn *websocket.Conn, u db.User, room *Room, Lesson db.Lesson, params ClientParams) *Client {
	params.PingPeriod = (params.PongWait * 9) / 10
	client := Client{
		ID:       u.ID,
		Name:     u.Name,
		conn:     conn,
		send:     make(chan []byte, 256),
		room:     room,
		LessonID: Lesson.ID,

		ClientParams: params,
	}
	go client.writePump()
	go client.readPump()

	return &client
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(client.MaxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(client.PongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(client.PongWait)); return nil })

	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(jsonMessage)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(client.PingPeriod)
	defer func() {
		for c := range client.room.clients {
			c.room.unregister <- c
		}
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(client.WriteWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(client.WriteWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.room.unregister <- client
	close(client.send)
	client.conn.Close()
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	message.Sender = client

	switch message.Action {
	case SendMessageAction:
		client.room.broadcast <- &message
	case LeaveRoomAction:
		client.room.broadcast <- &message
		client.handleLeaveRoomMessage()
	}

}
func (client *Client) handleLeaveRoomMessage() {
	for c := range client.room.clients {
		c.room.unregister <- c
	}

	client.ExitFunc()
}
