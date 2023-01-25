package db

import "Backend/Errors"

type RefreshToken struct {
	ID        uint `gorm:"primarykey"`
	Refresh   string
	ExpiresIn int64
	UserLevel int // TODO valid
	UserID    uint

	Client string
}

func (db *DB) SetToken(t *RefreshToken) error {
	result := db.Engine.Create(t)
	if result.RowsAffected == 0 {
		return Errors.NoRowsAffectedError
	}
	return result.Error
}

func (db *DB) DelToken(t *RefreshToken) error {
	result := db.Engine.Where(t).Delete(&RefreshToken{})
	if result.RowsAffected == 0 {
		return Errors.NoRowsAffectedError
	}
	return result.Error
}
func (db *DB) CheckAndGetTocken(t *RefreshToken) bool {
	res := db.Engine.Where(t).First(t)
	if res.Error != nil || res.RowsAffected == 0 {
		return false
	}
	return true
}
func (db *DB) GetAllUsersDevices(userid uint) ([]string, error) {
	var devices []string
	res := db.Engine.Model(&RefreshToken{}).Select("client").Where("user_id = ?", userid).Find(&devices)
	if res.RowsAffected == 0 {
		return nil, Errors.NoRowsAffectedError
	}
	return devices, res.Error
}
