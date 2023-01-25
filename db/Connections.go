package db

//
//import (
//	"gorm.io/gorm"
//	"time"
//)
//
//const (
//	ParentToChildConnectionType = iota + 1
//	// Support To User
//)
//
//type TypesOfConnection int
//
//func (t TypesOfConnection) GetType() string {
//	return []string{"Parent To Child"}[t-1]
//}
//
//type UserToUser struct {
//	gorm.Model
//
//	Type       int
//	TypeString string
//
//	DominantID  uint
//	DependentID uint
//
//	Expires int64
//}
//
//func (u *UserToUser) BeforeCreate(tx *gorm.DB) error {
//	u.TypeString = TypesOfConnection(u.Type).GetType()
//	return nil
//}
//
//func (db *DB) CreateChild(parentid uint, u *User) error {
//	err := db.CreateUser(u)
//	if err != nil {
//		return err
//	}
//	res := db.Engine.Create(&UserToUser{
//		Type:        1,
//		DominantID:  parentid,
//		DependentID: u.ID,
//		Expires:     0,
//	})
//	return res.Error
//}
//func (db *DB) DeleteChild(parentid uint, c *User) error {
//	err := db.Engine.Where(c).Delete(c).Error
//	if err != nil {
//		return err
//	}
//	res := db.Engine.Where(&UserToUser{
//		Type:        1,
//		DominantID:  parentid,
//		DependentID: c.ID,
//	}).Delete(&UserToUser{})
//
//	if res.RowsAffected == 0 {
//		return NoRowsAffectedError
//	}
//	return res.Error
//}
//
//func (db *DB) GetUsersChild(userid uint) ([]User, error) {
//	var ids []uint
//	res := db.Engine.Model(&UserToUser{}).Where("dominant_id = ? AND (expires = 0 OR expires > ?) AND type = 1", userid, time.Now().Unix()).Select("dependent_id").Find(&ids)
//	if res.Error != nil {
//		return nil, res.Error
//	} else if res.RowsAffected == 0 {
//		return nil, nil
//	}
//
//	var users []User
//	res = db.Engine.Model(&User{}).Find(&users, ids)
//	if res.Error != nil {
//		return nil, res.Error
//	}
//	return users, nil
//}
//func (db *DB) GetUsersParents(userid uint) ([]User, error) {
//	var ids []uint
//	res := db.Engine.Model(&UserToUser{}).Where("dependent_id = ? AND (expires = 0 OR expires > ?) AND type = 1", userid, time.Now().Unix()).Select("dominant_id").Find(&ids)
//	if res.Error != nil {
//		return nil, res.Error
//	} else if res.RowsAffected == 0 {
//		return nil, nil
//	}
//
//	var users []User
//	res = db.Engine.Model(&User{}).Find(&users, ids)
//	if res.Error != nil {
//		return nil, res.Error
//	}
//	return users, nil
//}
