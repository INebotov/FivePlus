package db

import (
	"fmt"
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

func (c *Confirmation) SendPhone(phone string) error {
	fmt.Printf("Confirmation send on Phone: %s;\nAction: %s\nCode: %d\n", phone, c.ActionString, c.Code)
	return nil
}
func (c *Confirmation) SendEmail(email string) error {
	fmt.Printf("Confirmation send on Email: %s;\nAction: %s\nCode: %d\n", email, c.ActionString, c.Code)
	return nil
}

func (db *DB) CreateConfirmation(c *Confirmation, value string, expired time.Duration) error {
	c.GenerateConfirmation(expired)
	if c.Type == EmailConfirmationType {
		err := c.SendEmail(value)
		if err != nil {
			return err
		}
	} else if c.Type == PhoneConfirmationType {
		err := c.SendPhone(value)
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
