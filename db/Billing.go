package db

import (
	"Backend/Errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type BillingAccount struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserBillingID string `gorm:"not null" json:"-"`

	Balance float64 `json:"balance"`
}

func (b *BillingAccount) BeforeCreate(tx *gorm.DB) (err error) {
	b.Balance = 0
	return nil
}

type Transaction struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"satisfied_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	SenderBillingID  string `gorm:"not null" json:"-"`
	ReciverBillingID string `gorm:"not null" json:"-"`

	SenderID  uint `json:"sender_id" gorm:"-"`
	ReciverID uint `json:"Reciver_id" gorm:"-"`

	Amount  float64 `gorm:"not null" json:"amount"`
	Pending bool    `json:"pending"`
}

func (b *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if !b.Pending && b.Amount == 0 {
		return Errors.BadRequest
	}

	var counter int64
	err = tx.Model(&BillingAccount{}).Where("user_billing_id IN (?, ?)", b.SenderBillingID, b.ReciverBillingID).Count(&counter).Error
	if err != nil {
		return err
	}
	if counter != 2 {
		return Errors.BadRequest
	}

	// TODO Low balance send notification
	return nil
}

func (db *DB) CreateBillingAccount(billingId string) error {
	account := BillingAccount{
		UserBillingID: billingId,
		Balance:       0,
	}
	res := db.Engine.Create(&account)
	if res.RowsAffected == 0 {
		return Errors.NoRowsAffectedError
	}
	return res.Error
}
func (db *DB) UpdateBillingAccount(id string, newBalance float64) error {
	res := db.Engine.Model(&BillingAccount{}).Where("user_billing_id = ?", id).Update("balance", newBalance)
	if res.RowsAffected == 0 {
		return Errors.NoRowsAffectedError
	}
	return res.Error
}

func (db *DB) GetBillingAccount(account *BillingAccount) error {
	return db.Engine.Where(account).First(account).Error
}
func (db *DB) GetBalance(billingID string) (float64, error) {
	var currbalance float64
	err := db.Engine.Model(&BillingAccount{}).Where("user_billing_id = ?", billingID).Select("balance").First(&currbalance).Error
	if err != nil {
		return 0, err
	}
	return currbalance, nil
}

func (db *DB) AddTransaction(t *Transaction) error {
	var sender BillingAccount
	sender.UserBillingID = t.SenderBillingID
	err := db.GetBillingAccount(&sender)
	if err != nil {
		return err
	}

	var reciver BillingAccount
	reciver.UserBillingID = t.ReciverBillingID
	err = db.GetBillingAccount(&reciver)
	if err != nil {
		return err
	}

	t.SenderID = sender.ID
	t.ReciverID = reciver.ID

	t.Amount = 0
	t.Pending = true

	return db.Engine.Create(t).Error
}
func (db *DB) PreformTransaction(t *Transaction) error {
	if t.Amount == 0 {
		return Errors.BadRequest
	}

	q := db.Engine.Model(&Transaction{}).Where(t)

	var err error
	if err = q.Find(t).Error; err != nil {
		return err
	}

	sender := BillingAccount{UserBillingID: t.SenderBillingID}
	if err = db.GetBillingAccount(&sender); err != nil {
		return err
	}

	reciver := BillingAccount{UserBillingID: t.ReciverBillingID}
	if err = db.GetBillingAccount(&reciver); err != nil {
		return err
	}

	newBalance := sender.Balance - t.Amount

	if err = q.Update("pending", false).Error; err != nil {
		return err
	}

	err1 := db.UpdateBillingAccount(sender.UserBillingID, newBalance)
	err = db.UpdateBillingAccount(reciver.UserBillingID, reciver.Balance+t.Amount)
	if err != nil || err1 != nil {
		db.Log.Error("Danger! One of Billing accounts may be not updated!!!!")
		return fmt.Errorf("First: %s\t Second: %s", err, err1)
	}
	return nil
}
func (db *DB) GetTransaction(t *Transaction) error {
	return db.Engine.Where(t).First(t).Error
}

func (db *DB) GetAllTransactions(id string) ([]Transaction, error) {
	var res []Transaction
	return res, db.Engine.Model(&Transaction{}).Where("sender_billing_id = ? OR reciver_billing_id = ?", id, id).Find(&res).Error
}
