package Chat

import (
	"github.com/google/uuid"
)

type Room struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	users      []uint
}

// NewRoom creates a new Room
func NewRoom(name string, users []uint) Room {
	return Room{
		ID:         uuid.New().String(),
		Name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		users:      users,
	}
}

func (room *Room) RunRoom() {
	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.BroadcastMessage(message.encode(), message.Sender)
		}

	}
}

func (room *Room) registerClientInRoom(client *Client) {
	mess := Message{
		Action: JoinRoomAction,
		Sender: client,
	}
	room.BroadcastMessage(mess.encode(), client)
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) BroadcastMessage(message []byte, client *Client) {
	for c := range room.clients {
		if c.ID != client.ID {
			c.send <- message
		}
	}
}

func (room *Room) GetName() string {
	return room.Name
}
