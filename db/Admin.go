package db

import (
	"Backend/Errors"
)

type SubjectToTeacher struct {
	SubjectID uint `gorm:"primarykey"`
	TeacherID uint `gorm:"primarykey"`
}
type Subject struct {
	ID   uint `gorm:"primarykey"`
	Name string

	Photo       string
	Description string

	StaticID string
}

type GradeToTeacher struct {
	GradeID   uint `gorm:"primarykey"`
	TeacherID uint `gorm:"primarykey"`
}
type Grade struct {
	ID          uint   `gorm:"primarykey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
}

func (db *DB) CreateSubject(s *Subject) error {
	var count int64
	db.Engine.Model(&Subject{}).Where(s).Count(&count)
	if count != 0 {
		return Errors.AlreadyExistsError
	}

	return db.Engine.Create(s).Error
}

func (db *DB) CreateGrade(s *Grade) error {
	var count int64
	db.Engine.Model(&Grade{}).Where(s).Count(&count)
	if count != 0 {
		return Errors.AlreadyExistsError
	}

	return db.Engine.Create(s).Error
}
func (db *DB) MakeUserTeacher(userid uint) error {
	var cou int64
	err := db.Engine.Model(&Teacher{}).Where("user_id = ?", userid).Count(&cou).Error
	if err != nil {
		return err
	}
	if cou != 0 {
		return Errors.AlreadyExistsError
	}

	user := User{}
	user.ID = userid
	user.AccessLevel = TeacherAccessLevel
	user.LevelString = LevelsOfAccess(TeacherAccessLevel).GetLevel()

	err = db.UpdateUser(user)
	if err != nil {
		return err
	}

	var teacher Teacher
	teacher.UserID = userid

	return db.Engine.Create(&teacher).Error
}
