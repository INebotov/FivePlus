package db

import (
	"Backend/Errors"
	"Backend/auth"
	"reflect"
	"time"
)

type Session struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Refresh   string
	ExpiresIn int64

	UserLevel int
	UserID    uint

	Device string
}

func (db *DB) CreateSession(session Session) error {
	if len(session.Refresh) >= 16 && session.ExpiresIn > time.Now().Unix() && session.UserID != 0 && session.UserLevel < 6 && session.UserLevel > 0 {
		return db.Engine.Model(&Session{}).Create(&session).Error
	}
	return Errors.SessionIsIncorrectError
}

func (db *DB) RefreshSession(old, new auth.Token) error {
	if reflect.DeepEqual(time.Time{}, new.Expires) {
		return Errors.TokenIsIvalidError
	}
	return db.Engine.Model(&Session{}).Where("refresh = ?", old.Key).Limit(1).Update("refresh", new.Key).Update("expires_in", new.Expires).Error
}

func (db *DB) DeleteSession(token auth.Token) error {
	return db.Engine.Model(&Session{}).Where("refresh = ?", token.Key).Limit(1).Delete(&Session{}).Error
}

func (db *DB) GetUserDataFromRefresh(token string) (uint, error) {
	var id uint
	return id, db.Engine.Model(&Session{}).Where("refresh = ? AND expires_in > ?", token, time.Now().Unix()).
		Select("user_id").First(&id).Error
}

func (db *DB) CheckSession(token auth.Token) error {
	now := time.Now().Unix()

	if !reflect.DeepEqual(time.Time{}, token.Expires) && token.Expires.Unix() < now {
		return Errors.TokenExpiredError
	}
	var count int64
	err := db.Engine.Model(&Session{}).Where("refresh = ? AND expires_in > ?", token.Key, now).Count(&count).Error

	if count == 0 {
		return Errors.TokenIsIvalidError
	}
	return err
}

type Device struct {
	Name string `json:"name"`
	Date int64  `json:"date"`
}

func (db *DB) DetAllUserDevices(userid uint) []Device {
	var dev []Device
	err := db.Engine.Model(&Session{}).Where("user_id = ? AND expires_in > ?", userid, time.Now().
		Unix()).Order("updated_at").Select([]string{"device", "updated_at"}).
		Select("device as name, updated_at as date").Find(&dev).Error
	if err != nil {
		return nil
	}
	return dev
}
