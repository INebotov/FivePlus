package db

import (
	"Backend/Errors"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

const (
	EmailConfirmationType = iota + 1
	PhoneConfirmationType
)

type TypesOfConfirmation int

func (t TypesOfConfirmation) String() string {
	return []string{"email", "phone"}[t-1]
}
func (t *TypesOfConfirmation) Parse(s string) {
	switch s {
	case "email":
		*t = TypesOfConfirmation(1)
	case "phone":
		*t = TypesOfConfirmation(2)
	}
}

type Confirmation struct {
	ID uint `gorm:"primarykey"`

	UserID uint `gorm:"not null"`

	Type int `gorm:"not null"`

	ActionType string `gorm:"not null"`

	Code int `gorm:"not null"`

	LiveUntill int64          `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (c *Confirmation) GenerateConfirmation(expired time.Duration) {
	c.LiveUntill = time.Now().Add(expired).Unix()
	c.Code = 1000 + rand.Intn(8999)
}

func (db *DB) CreateConfirmation(c *Confirmation, expired time.Duration) error {
	c.GenerateConfirmation(expired)
	return db.Engine.Create(c).Error
}
func (db *DB) ConfirmOperation(userid uint, code int, action string) error {
	var conf Confirmation
	qery := db.Engine.Where("user_id = ? AND code = ? AND action_type = ?", userid, code, action)

	if err := qery.First(&conf).Error; err != nil {
		return err
	}

	if time.Now().Unix() > conf.LiveUntill {
		return Errors.ConfirmationExpiredError
	}

	return qery.Delete(&Confirmation{}).Error
}
