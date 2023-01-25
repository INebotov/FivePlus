package db

import (
	"Backend/Errors"
	"gorm.io/gorm"
	"time"
)

type LessonType struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Unit  int     `json:"unit"`
	Price float64 `json:"price"`
}

func (l *LessonType) BeforeCreate(tx *gorm.DB) (err error) {
	if len(l.Name) < 6 || l.Description == "" || l.Unit == 0 || l.Price < 5 {
		return Errors.BadRequest
	}

	var count int64
	if err := tx.Model(&LessonType{}).Where("name = ?", l.Name).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return Errors.AlreadyExistsError
	}

	return nil
}

func (db *DB) GetLessonType(r *LessonType) error {
	return db.Engine.Where(r).Find(r).Error
}
func (db *DB) CreateLessonType(r *LessonType) error {
	return db.Engine.Create(r).Error
}

type Lesson struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Name      string `json:"name"`
	SubjectID uint   `json:"subject_id"`

	TypeID uint `json:"type_id" gorm:"not null"`

	TeacherID uint `json:"teacher_id"`
	StudentID uint `json:"student_id"`

	TimeRequested int64 `json:"time_requested" gorm:"not null"`

	TimeStarted int64 `json:"time_started"`
	TimeEnded   int64 `json:"time_ended"`

	TimeErased int64 `json:"time_erased"`
	LastErase  int64 `json:"-"`
	Paused     bool

	SuspendMessage string `json:"message"`

	TransactionID uint `json:"-"`

	ChatID string `json:"chat_id"`

	TeacherConfirmed bool `json:"-"`
	StudentConfirmed bool `json:"-"`
}

func (l *Lesson) BeforeCreate(tx *gorm.DB) (err error) {
	if len(l.Name) < 6 || l.TypeID == 0 || l.SubjectID == 0 || l.StudentID == 0 {
		return Errors.BadRequest
	}
	var count int64
	err = tx.Model(&LessonType{}).Where("id = ?", l.TypeID).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return Errors.BadRequest
	}

	l.TimeRequested = time.Now().Unix()
	l.ChatID = ""
	l.TimeStarted = 0
	l.TimeEnded = 0
	l.TeacherID = 0
	l.SuspendMessage = ""
	l.Paused = false
	l.LastErase = 0
	l.TimeErased = 0

	return nil
}

func (db *DB) GetLessons(userid uint, limit int, startfrom int) ([]Lesson, error) {
	var res []Lesson
	return res, db.Engine.Model(&Lesson{}).Where("(teacher_id = ? OR student_id = ?) AND time_ended > 0 AND id < ?", userid, userid, startfrom).
		Limit(limit).Order("id desc").Find(&res).Error
}

func (db *DB) CreateLesson(r *Lesson) error {
	var count int64
	if err := db.Engine.Model(&Lesson{}).Where("student_id = ?", r.StudentID).Count(&count).Error; err != nil {
		return err
	}

	if count != 0 {
		return Errors.AlreadyExistsError
	}

	return db.Engine.Create(r).Error
}
func (db *DB) GetLesson(r *Lesson) error {
	return db.Engine.Where(r).Find(r).Error
}

func (db *DB) SatisfyLessonRequest(teacherID uint, l *Lesson) error {
	var err error
	if err = db.GetLesson(l); err != nil {
		return err
	}

	teacher := Teacher{
		ID: teacherID,
	}
	if err = db.GetTeacherSimple(&teacher); err != nil {
		return err
	}

	var count int64
	if err = db.Engine.Model(&Lesson{}).Where("teacher_id = ? AND time_ended <= 0", teacherID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return Errors.AlreadyExistsError
	}

	PayerBillingID, err := db.GetUsersBillingID(User{ID: l.StudentID})
	if err != nil {
		return err
	}
	TeacherBillingID, err := db.GetUsersBillingID(User{ID: teacher.UserID})
	if err != nil {
		return err
	}

	t := Transaction{SenderBillingID: PayerBillingID, ReciverBillingID: TeacherBillingID}
	err = db.AddTransaction(&t)
	if err != nil {
		return err
	}

	l.TransactionID = t.ID
	l.TeacherConfirmed = true

	return db.Engine.Model(&Lesson{}).Updates(*l).Error
}
func (db *DB) StartLesson(l *Lesson) error {
	l.TimeStarted = time.Now().Unix()
	l.StudentConfirmed = true

	return db.Engine.Model(l).Updates(*l).Error
}
func (db *DB) EndLesson(l *Lesson) error {
	err := db.GetLesson(l)
	if err != nil {
		return err
	}

	l.TimeEnded = time.Now().Unix()
	timeSpent := float64((l.TimeEnded - l.TimeStarted) - l.TimeErased)

	lessonType := LessonType{ID: l.TypeID}
	err = db.GetLessonType(&lessonType)
	if err != nil {
		return err
	}

	t := Transaction{
		ID: l.TransactionID,
	}

	err = db.GetTransaction(&t)
	if err != nil {
		return err
	}

	t.Amount = (timeSpent / float64(lessonType.Unit)) * lessonType.Price

	err = db.Engine.Model(&Lesson{}).Where(l).Update("paused", false).Update("last_erase", 0).Error
	if err != nil {
		return err
	}

	err = db.PreformTransaction(&t)
	if err != nil {
		return err
	}

	return db.Engine.Model(l).Updates(*l).Error
}
func (db *DB) PauseLesson(l *Lesson) error {
	err := db.GetLesson(l)
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	if l.Paused {
		l.TimeErased += now - l.LastErase
		return db.Engine.Model(l).Update("paused", false).Update("last_erase", 0).
			Update("time_erased", l.TimeErased).Update("suspend_message", "").Error
	}

	l.Paused = true
	l.LastErase = now
	return db.Engine.Model(l).Updates(*l).Error
}

func (db *DB) MarkLesson(lessonid uint, c Comment) error {
	var l Lesson
	l.ID = lessonid
	err := db.GetLesson(&l)
	if err != nil {
		return err
	}

	if l.TimeEnded == 0 {
		return Errors.BadRequest
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
func (db *DB) GetAllPendingLessonRequests(subjectids []uint, limit int) ([]Lesson, error) {
	var res []Lesson
	return res, db.Engine.Model(&Lesson{}).Where("subject_id IN ? AND teacher_confirmed = ?", subjectids, false).
		Limit(limit).Order("created_at desc").Find(&res).Error
}
