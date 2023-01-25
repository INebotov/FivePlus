package db

import (
	"Backend/Errors"
	"github.com/google/uuid"
	"github.com/gookit/validate"
	"gorm.io/gorm"
	"regexp"
	"time"
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
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserName string `validate:"gte=3 & lte=25" gorm:"not null" json:"user_name"`

	Name           string `validate:"gte=3 & lte=60" gorm:"not null" json:"name"`
	Email          string `validate:"format=email" gorm:"not null" json:"email"`
	EmailConfermed bool   `json:"email_confermed"`

	Phone          string `json:"phone"`
	PhoneConfirmed bool   `json:"phone_confirmed"`

	Photo string `json:"photo"`
	About string `json:"about"`

	Password    string `json:"-"`
	AccessLevel int    `validate:"gte=1 & lte=5" gorm:"not null" json:"-"`
	LevelString string `json:"user_type"`

	BillingID string `json:"-"`
}

func (u *User) Validate() error {
	if !regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString(u.UserName) {
		return Errors.WrongUsernameError
	}
	if !regexp.MustCompile(`^(\+7|7|8)?[\s\-]?\(?[489][0-9]{2}\)?[\s\-]?[0-9]{3}[\s\-]?[0-9]{2}[\s\-]?[0-9]{2}$`).MatchString(u.Phone) {
		return Errors.PhoneIncorrectError
	}

	return nil
}
func (u *User) BeforeCreate(tx *gorm.DB) error {
	v := validate.Struct(u)
	if !v.Validate() {
		return Errors.BadRequest
	}

	u.LevelString = LevelsOfAccess(u.AccessLevel).GetLevel()

	var c int64
	tx.Where("email = ? OR user_name = ? OR phone = ?", u.Email, u.UserName, u.Phone).Count(&c)
	if c != 0 {
		return Errors.AlreadyExistsError
	}

	if u.AccessLevel != ChildAccessLevel {
		var id string
		var c int64 = 1
		i := 0
		for c != 0 || i < 3 {
			id = uuid.New().String()
			tx.Model(&User{}).Where("billing_id = ?", id).Count(&c)
			i++
		}
		if c != 0 {
			return Errors.CantCreateBillingIDError
		}
		u.BillingID = id
	}
	return nil
}

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
	return db.Engine.Where(user).First(user).Error
}
func (db *DB) UpdateUser(user User) error {
	return db.Engine.Model(&user).Updates(user).Error
}
func (db *DB) CheckUserPersistance(u User) bool {
	var c int64
	return c != 0 && db.Engine.Where(&u).Count(&c).Error == nil
}
func (db *DB) GetUserName(userid uint) string {
	var res string // Error isnt rasing!!
	db.Engine.Model(&User{}).Select("name").First(&res, userid)
	return res
}
func (db *DB) ChangeUserBool(userid uint, name string, status bool) error {
	return db.Engine.Model(&User{}).Where("id = ?", userid).Update(name, status).Error
}

func (db *DB) GetUsersBillingID(u User) (string, error) {
	var res string
	return res, db.Engine.Model(&User{}).Where(&u).Select("billing_id").First(&res).Error
}

func (db *DB) GetIDByUsername(username string) (uint, error) {
	var res uint
	return res, db.Engine.Model(&User{}).Where("user_name = ?", username).Select("id").First(&res).Error
}
func (db *DB) GetUsernameByID(id uint) (string, error) {
	var res string
	return res, db.Engine.Model(&User{}).Where("id = ?", id).Select("user_name").First(&res).Error
}
func (db *DB) GetUserLevel(id uint) int {
	var res int = 0
	db.Engine.Model(&User{}).Where("id = ?", id).Select("access_level").First(&res)
	return res
}
