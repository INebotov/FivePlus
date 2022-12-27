package db

import (
	"BackendSimple/Sender"
	"context"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

const (
	EmailConfirmationType = iota + 1
	PhoneConfirmationType
)

type TypesOfConfirmation int

func (t TypesOfConfirmation) GetType() string {
	return []string{"Email", "Phone Confirmation"}[t-1]
}

const (
	RegistrationConfirmationAction = iota + 1
	ProfileSecretsChangeConfirmationAction
	DeteteChildAction
)

type ConfirmationsActions int

func (t ConfirmationsActions) GetType() string {
	return []string{"Registration", "User Profile Secrets Change", "Delete Child"}[t-1]
}

type Confirmation struct {
	ID uint `gorm:"primarykey"`

	UserID uint `gorm:"not null"`

	Type       int    `gorm:"not null"`
	TypeString string `gorm:"not null"`

	Action       int    `gorm:"not null"`
	ActionString string `gorm:"not null"`

	Code int `gorm:"not null"`

	LiveUntill int64          `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (c *Confirmation) GenerateConfirmation(expired time.Duration) {
	if c.ActionString == "" {
		c.ActionString = ConfirmationsActions(c.Action).GetType()
	}
	if c.TypeString == "" {
		c.TypeString = TypesOfConfirmation(c.Type).GetType()
	}
	c.LiveUntill = time.Now().Add(expired).Unix()
	c.GenerateCode()
}

func (c *Confirmation) GenerateCode() {
	c.Code = 1000 + rand.Intn(8999)
}
func (db *DB) CreateConfirmation(c *Confirmation, u User, expired time.Duration) error {
	c.GenerateConfirmation(expired)
	if c.Type == EmailConfirmationType {
		err := db.Sender.SendEmail(context.Background(), Sender.EmailMessage{
			UserName:  u.Name,
			Type:      c.Action,
			Code:      c.Code,
			Recipient: u.Email,
		})
		if err != nil {
			return err
		}
	} else {
		return WrongConfirmationTypeError
	}

	result := db.Engine.Create(c)
	if result.RowsAffected == 0 {
		return NoRowsAffectedError
	}
	return result.Error
}
func (db *DB) ConfirmOperation(userid uint, code int, action int) error {
	var conf Confirmation // TODO: In handlers from request -> to func Iject possible danger!
	qery := db.Engine.Where("user_id = ? AND code = ? AND action = ?", userid, code, action)
	res := qery.First(&conf)
	if res.RowsAffected == 0 {
		return NoRowsAffectedError
	} else if res.Error != nil {
		return res.Error
	}

	if time.Now().Unix() > conf.LiveUntill {
		return ConfirmationExpiredError
	}

	res = qery.Delete(&Confirmation{})
	if res.RowsAffected == 0 {
		return NoRowsAffectedError
	} else if res.Error != nil {
		return res.Error
	}
	return nil
}
