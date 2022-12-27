package db

import (
	"time"
)

type ChatRoom struct {
	ID string `gorm:"primarykey"`

	CreatedAt time.Time
	ClosedAt  time.Time

	Messages []Message `gorm:"foreignKey:ID;many2many:room_message"`
}

type Message struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time

	SenderID uint
	Message  string
}

func (db *DB) AddMessage(mess *Message, chatId string) error {
	return db.Engine.Model(&ChatRoom{}).Where("id = ?", chatId).Association("Messages").Append(&mess)
}

func (db *DB) CreateRoom(room *ChatRoom) error {
	return db.Engine.Create(&room).Error
}

func (db *DB) CloseRoom(room *ChatRoom) error {
	room.ClosedAt = time.Now()
	return db.Engine.Updates(*room).Error
}

func (db *DB) GetAllUserChats(userid uint, limit int) ([]ChatRoom, error) {
	var lessons []Lesson

	err := db.Engine.Model(&Lesson{}).Where("teacher_id = ? OR StudentID = ? ", userid, userid).Limit(limit).
		Order("time_ended").Find(&lessons).Error
	if err != nil {
		return nil, err
	}

	chatids := make([]string, len(lessons))
	for i, e := range lessons {
		chatids[i] = e.ChatID
	}
	var res []ChatRoom
	return res, db.Engine.Model(&ChatRoom{}).Find(&res, chatids).Error
}

func (db *DB) GetChat(chat *ChatRoom) error {
	return db.Engine.Model(&ChatRoom{}).Where(chat).Preload("Messages").Find(chat).Error
}

func (db *DB) IsChatUsers(chatID string, userid uint) bool {
	var counter int64
	teachID, err := db.GetTeacherIDByUserId(userid)
	if err != nil {
		return false
	}
	if db.Engine.Model(&Lesson{}).Where("ChatID = ? AND ( teacher_id = ? OR student_id = ? )", chatID, teachID, userid).Count(&counter).Error != nil {
		return false
	}
	return counter != 0
}
