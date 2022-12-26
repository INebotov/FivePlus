package db

import (
	"gorm.io/gorm"
)

type SubjectToTeacher struct {
	SubjectID uint `gorm:"primarykey"`
	TeacherID uint `gorm:"primarykey"`
}

type Comment struct {
	ID uint `gorm:"primarykey"`

	UserID  uint
	Mark    int
	Message string
}

type Teacher struct {
	ID     uint `gorm:"primarykey"`
	UserID uint `gorm:"not null;uniqueIndex"`

	Comments []Comment `gorm:"foreignKey:ID;many2many:teacher_comments"`

	ResultMark int

	Active bool

	// Ignoring by DB
	Subjects []Subject `gorm:"-"`
}

func (db *DB) MakeUserTeacher(userid uint) error {
	var cou int64
	err := db.Engine.Model(&Teacher{}).Where("user_id = ?", userid).Count(&cou).Error
	if err != nil {
		return err
	}
	if cou != 0 {
		return AlreadyExistsError
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

func (db *DB) GetTeacherIDByUserId(userid uint) (uint, error) {
	var id uint
	return id, db.Engine.Model(&Teacher{}).Where("user_id = ?", userid).Select("id").First(&id).Error
}

func (db *DB) UpdateTeacher(t Teacher) error {
	result := db.Engine.Model(&t).Updates(t)
	if result.RowsAffected == 0 {
		return NoRowsAffectedError
	}
	return result.Error
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
func (db *DB) GetSubject(s *Subject) error {
	return db.Engine.Where(s).First(s).Error
}
func (db *DB) GetTeacherBySubject(subid uint, limit int) ([]Teacher, error) {
	var tids []uint
	err := db.Engine.Model(&SubjectToTeacher{}).Where("subject_id = ?", subid).Select("teacher_id").Limit(limit).Find(&tids).Error
	if err != nil {
		return nil, err
	}

	res, err := db.GetTeachers(tids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db *DB) TeacherAddSubject(teacherid, subjectid uint) error {
	var c int64
	db.Engine.Model(&Subject{}).Where("id = ?", subjectid).Count(&c)
	if c == 0 {
		return gorm.ErrRecordNotFound
	}

	return db.Engine.Model(&SubjectToTeacher{}).Create(&SubjectToTeacher{
		SubjectID: subjectid,
		TeacherID: teacherid,
	}).Error
}
func (db *DB) TeacherDeleteSubject(teacherid, subjectid uint) error {
	return db.Engine.Model(&SubjectToTeacher{}).Delete(&SubjectToTeacher{
		SubjectID: subjectid,
		TeacherID: teacherid,
	}).Error
}

func (db *DB) GetLessons(userid uint, limit int, startfrom int) ([]Lesson, error) {
	var res []Lesson
	return res, db.Engine.Model(&Lesson{}).Where("(teacher_id = ? OR student_id = ?) AND time_ended > 0 AND id > ?", userid, userid, startfrom).
		Limit(limit).Order("id desc, time_ended").Find(&res).Error
}
func (db *DB) GetLesson(l *Lesson) error {
	return db.Engine.Where(l).First(l).Error
}

func (db *DB) MarkLesson(lessonid uint, c Comment) error {
	var l Lesson
	l.ID = lessonid
	err := db.GetLesson(&l)
	if err != nil {
		return err
	}

	err = db.Engine.Model(&Teacher{ID: l.TeacherID}).Association("comments").Append(&c)
	if err != nil {
		return err
	}

	t := Teacher{ID: l.TeacherID}
	err = db.GetTeacherSimple(&t)
	if err != nil {
		return err
	}

	t.ResultMark = (t.ResultMark + c.Mark) / 2
	return db.UpdateTeacher(t)
}

func (db *DB) ChangeActive(id uint, active bool) error {
	return db.Engine.Model(&Teacher{ID: id}).Update("active", active).Error
}
