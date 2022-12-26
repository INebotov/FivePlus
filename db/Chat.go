package db

import (
	"gorm.io/gorm"
	"time"
)

type ChatRoom struct {
	ID string `gorm:"primarykey"`

	LessonID uint

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Message struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time

	SenderID uint
	RoomID   string
	Message  string
}

func (db *DB) CreateRoom(room *ChatRoom) error {
	return db.Engine.Create(&room).Error
}
func (db *DB) CloseRoom(room *ChatRoom) error {
	return db.Engine.Delete(&room).Error
}
