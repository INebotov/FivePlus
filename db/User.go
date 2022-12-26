package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"regexp"
)

const (
	ChildAccessLevel = iota + 1
	GeneralAccessLevel
	TeacherAccessLevel
	SupporterAccessLevel
	AdminAccessLevel
)

type LevelsOfAccess int

func (l LevelsOfAccess) GetLevel() string {
	return []string{"Child", "General", "Teacher", "Supporter", "Admin"}[l-1]
}

type User struct {
	gorm.Model

	UserName string `validate:"gte=3 & lte=25" gorm:"not null"`

	Name  string `validate:"gte=3 & lte=60" gorm:"not null"`
	Email string `validate:"format=email" gorm:"not null"`

	EmailConfermed bool

	Photo string

	Password    []byte
	AccessLevel int `validate:"gte=1 & lte=5" gorm:"not null"`
	LevelString string

	BillingID string
}

func (u *User) Validate() error {
	if !regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString(u.UserName) {
		return WrongUsernameError
	}
	if len(u.Password) != 64 {
		return WrongPasswordHashError
	}
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.LevelString = LevelsOfAccess(u.AccessLevel).GetLevel()

	if u.AccessLevel != ChildAccessLevel {
		var id string
		var c int64 = 1
		i := 0
		for c != 0 || i < 3 {
			id = uuid.New().String()
			tx.Model(&User{}).Where("user_billing_id = ?", id).Count(&c)
			i++
		}
		if c != 0 {
			return CantCreateBillingIDError
		}
		u.BillingID = id
	}
	return nil
} // But On update need to do it by hands

func (db *DB) CreateUser(user *User) error {
	result := db.Engine.Create(user)
	return result.Error
}
func (db *DB) DeleteUser(user *User) error {
	return db.Engine.Delete(user).Error
}
func (db *DB) GetUserById(userid uint) (User, error) {
	var userRes User
	return userRes, db.Engine.First(&userRes, userid).Error
}
func (db *DB) GetUser(user *User) error {
	result := db.Engine.Where(user).First(user)
	if result.RowsAffected == 0 {
		return NoRowsAffectedError
	}
	return result.Error
}
func (db *DB) UpdateUser(user User) error {
	result := db.Engine.Model(&user).Updates(user)
	if result.RowsAffected == 0 {
		return NoRowsAffectedError
	}
	return result.Error
}

func (db *DB) CheckUserPersistance(u User) bool {
	var user User
	res := db.Engine.Where("email = ?", u.Email).First(&user)
	if res.Error != nil || res.RowsAffected == 0 {
		return false
	}
	return true
}
func (db *DB) GetUserName(userid uint) string {
	var res string // Error isnt rasing!!
	db.Engine.Model(&User{}).Select("name").First(&res, userid)
	return res
}
func (db *DB) ChangeUserBool(userid uint, name string, status bool) error {
	result := db.Engine.Model(&User{}).Where("id = ?", userid).Update(name, status)
	if result.RowsAffected == 0 {
		return NoRowsAffectedError
	}
	return result.Error
}

func (db *DB) GetUsersBillingID(u User) (string, error) {
	var res string
	return res, db.Engine.Model(&User{}).Where(&u).Select("billing_id").First(&res).Error
}

func (db *DB) GetIDByUsername(username string) (uint, error) {
	var res uint
	result := db.Engine.Model(&User{}).Where("user_name = ?", username).Select("id").First(&res)
	if result.RowsAffected == 0 {
		return 0, NoRowsAffectedError
	}
	return res, result.Error
}
func (db *DB) GetUsernameByID(id uint) (string, error) {
	var res string
	result := db.Engine.Model(&User{}).Where("id = ?", id).Select("user_name").First(&res)
	if result.RowsAffected == 0 {
		return "", NoRowsAffectedError
	}
	return res, result.Error
}
