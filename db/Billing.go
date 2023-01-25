package db

//
//import (
//	"Backend/Errors"
//	"gorm.io/gorm"
//)
//
//type BillingAccount struct {
//	gorm.Model
//
//	UserBillingID string `gorm:"not null"`
//
//	Balance int64
//}
//
//type Transaction struct {
//	gorm.Model
//
//	SenderBillingID  string `gorm:"not null"`
//	ReciverBillingID string `gorm:"not null"`
//
//	Amount    int64 `gorm:"not null"`
//	Confirmed bool
//}
//
//func (db *DB) CreateBillingAccount(billingId string) error {
//	account := BillingAccount{
//		UserBillingID: billingId,
//		Balance:       0,
//	}
//	res := db.Engine.Create(&account)
//	if res.RowsAffected == 0 {
//		return Errors.NoRowsAffectedError
//	}
//	return res.Error
//}
//func (db *DB) UpdateBillingAccount(id string, newBalance int64) error {
//	res := db.Engine.Model(&BillingAccount{}).Where("user_billing_id = ?", id).Update("balance", newBalance)
//	if res.RowsAffected == 0 {
//		return Errors.NoRowsAffectedError
//	}
//	return res.Error
//}
//
//func (db *DB) GetBillingAccount(account *BillingAccount) error {
//	res := db.Engine.Where(account).First(account)
//	if res.RowsAffected == 0 {
//		return Errors.NoRowsAffectedError
//	}
//	return res.Error
//}
//func (db *DB) CheckBalance(billingID string, amount float64) bool {
//	var currbalance float64
//	err := db.Engine.Model(&BillingAccount{}).Where("user_billing_id = ?", billingID).Select("balance").First(&currbalance).Error
//	return err == nil && currbalance >= amount
//}
//
//func (db *DB) AddTransaction(t *Transaction) error {
//	var sender BillingAccount
//	sender.UserBillingID = t.SenderBillingID
//	err := db.GetBillingAccount(&sender)
//	if err != nil {
//		return err
//	}
//
//	var reciver BillingAccount
//	reciver.UserBillingID = t.ReciverBillingID
//	err = db.GetBillingAccount(&reciver)
//	if err != nil {
//		return err
//	}
//
//	newBalance := sender.Balance - t.Amount
//
//	err = db.UpdateBillingAccount(sender.UserBillingID, newBalance)
//	if err != nil {
//		return err
//	}
//	err = db.UpdateBillingAccount(reciver.UserBillingID, reciver.Balance+t.Amount)
//	if err != nil {
//		return err
//	}
//	return db.Engine.Create(t).Error
//}
//
//func (db *DB) GetAllTransactions(id string) ([]Transaction, error) {
//	var res []Transaction
//	return res, db.Engine.Model(&Transaction{}).Where("sender_billing_id = ? OR reciver_billing_id = ?", id, id).Find(&res).Error
//}
