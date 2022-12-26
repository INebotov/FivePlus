package db

import (
	"gorm.io/gorm"
	"time"
)

type Subject struct {
	ID   uint `gorm:"primarykey"`
	Name string

	Photo       string
	Description string

	StaticID string
}

type Lesson struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	SubjectID uint

	TeacherID        uint
	TeacherBillingID string

	StudentID      uint
	PayerBillingID string

	TimeStarted int64
	TimeEnded   int64

	ChatID string `gorm:"uniqueIndex"`
}

type LessonRequest struct {
	ID uint `gorm:"primarykey"`

	UserID    uint
	CreatedAt time.Time

	SubjectID uint
	Question  string

	Satisfied bool
	LessonID  uint
}

func (db *DB) CreateLessonRequest(r *LessonRequest) error {
	r.Satisfied = false
	r.LessonID = 0

	return db.Engine.Create(r).Error
}
func (db *DB) GetLessonRequest(r *LessonRequest) error {
	return db.Engine.Where(r).Find(r).Error
}

func (db *DB) GetAllPendingLessonRequests(subjectids []uint, limit int) ([]LessonRequest, error) {
	var res []LessonRequest
	return res, db.Engine.Model(&LessonRequest{}).Where("subject_id IN ? AND satisfied = ?", subjectids, false).Limit(limit).Order("created_at desc").Find(&res).Error
}
func (db *DB) CreateLesson(l *Lesson, r *LessonRequest) error {
	l.SubjectID = r.SubjectID

	teacher := Teacher{
		ID: l.TeacherID,
	}
	err := db.GetTeacherSimple(&teacher)
	if err != nil {
		return err
	}

	userid, err := db.GetUsersBillingID(User{Model: gorm.Model{ID: l.StudentID}})
	if err != nil {
		return err
	}
	teacherid, err := db.GetUsersBillingID(User{Model: gorm.Model{ID: teacher.UserID}})
	if err != nil {
		return err
	}

	l.PayerBillingID = userid
	l.TeacherBillingID = teacherid

	err = db.Engine.Create(&l).Error
	if err != nil {
		return err
	}

	r.Satisfied = true
	r.LessonID = l.ID

	return db.Engine.Model(r).Updates(*r).Error
}

func (db *DB) StartLesson(l *Lesson) error {
	l.TimeStarted = time.Now().Unix()

	return db.Engine.Model(l).Updates(*l).Error
}

func (db *DB) EndLesson(l *Lesson) error {
	err := db.GetLesson(l)
	if err != nil {
		return err
	}

	l.TimeEnded = time.Now().Unix()

	err = db.Engine.Model(l).Updates(*l).Error
	if err != nil {
		return err
	}

	timeSpent := l.TimeEnded - l.TimeStarted

	t := Transaction{
		SenderBillingID:  l.PayerBillingID,
		ReciverBillingID: l.TeacherBillingID,
		Amount:           timeSpent / 60,
	}
	err = db.AddTransaction(&t)
	if err != nil {
		return err
	}
	return nil
}
