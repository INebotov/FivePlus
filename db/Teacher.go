package db

import (
	"gorm.io/gorm"
)

type Comment struct {
	ID        uint `gorm:"primarykey"`
	TeacherID uint

	UserID  uint
	Mark    float32
	Message string
}

type Teacher struct {
	ID     uint `gorm:"primarykey" json:"id"`
	UserID uint `gorm:"not null;uniqueIndex" json:"user_id"`

	Comments []Comment `json:"comments"`

	ResultMark float32 `json:"result_mark"`

	Active bool `json:"active"`

	// Ignoring by DB
	Subjects []Subject `gorm:"-" json:"subjects"`
	Grades   []Grade   `gorm:"-" json:"grades"`
}

func (db *DB) GetTeacher(t *Teacher) error {
	err := db.Engine.Model(&Teacher{}).Where(t).Preload("Comments").First(t).Error
	if err != nil {
		return err
	}

	var subsIds []uint
	err = db.Engine.Model(&SubjectToTeacher{}).Where("teacher_id = ?", t.ID).Select("subject_id").Find(&subsIds).Error
	if err != nil {
		return err
	}

	var Subjects []Subject
	err = db.Engine.Model(&Subject{}).Find(&Subjects, subsIds).Error
	if err != nil {
		return err
	}
	t.Subjects = Subjects
	return nil
}
func (db *DB) GetTeacherSimple(t *Teacher) error {
	return db.Engine.Where(t).First(t).Error
}

func (db *DB) GetTeacherIDByUserId(userid uint) uint {
	var id uint
	if db.Engine.Model(&Teacher{}).Where("user_id = ?", userid).Select("id").
		First(&id).Error != nil {
		return 0
	}
	return id
}
func (db *DB) UpdateTeacher(t Teacher) error {
	return db.Engine.Model(&t).Updates(t).Error
}

func (db *DB) GetTeachers(ids []uint) ([]Teacher, error) {
	res := make([]Teacher, len(ids))
	for i, e := range ids {
		res[i].ID = e
		err := db.GetTeacherSimple(&res[i])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (db *DB) GetSubjectIDFromName(name string) (uint, error) {
	var res uint
	return res, db.Engine.Model(&Subject{}).Where("name = ?", name).Select("id").First(&res).Error
}

func (db *DB) AddSubject(subID, teachID uint) error {
	var c int64
	db.Engine.Model(&Subject{}).Where("id = ?", subID).Count(&c)
	if c == 0 {
		return gorm.ErrRecordNotFound
	}

	return db.Engine.Model(&SubjectToTeacher{}).Create(&SubjectToTeacher{
		SubjectID: subID,
		TeacherID: teachID,
	}).Error
}
func (db *DB) GetSubjects(s *Subject) error {
	return db.Engine.Where(s).First(s).Error
}
func (db *DB) DeleteSubject(subID, teachID uint) error {
	return db.Engine.Model(&SubjectToTeacher{}).
		Where("subject_id = ? AND teacher_id = ?", subID, teachID).Delete(&SubjectToTeacher{}).Error
}

func (db *DB) AddGrade(grID, userID uint) error {
	var c int64
	db.Engine.Model(&Grade{}).Where("id = ?", grID).Count(&c)
	if c == 0 {
		return gorm.ErrRecordNotFound
	}

	return db.Engine.Model(&GradeToTeacher{}).Create(&GradeToTeacher{
		GradeID:   grID,
		TeacherID: userID,
	}).Error
}
func (db *DB) GetGrades(s *Grade) error {
	return db.Engine.Where(s).First(s).Error
}
func (db *DB) DeleteGrade(grID, teachID uint) error {
	return db.Engine.Model(&GradeToTeacher{}).
		Where("grade_id = ? AND teacher_id = ?", grID, teachID).Delete(&GradeToTeacher{}).Error
}
