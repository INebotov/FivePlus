package db

import (
	"BackendSimple/Sender"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	TypePostgres = iota
	TypeSQLLite
)

type TypeOfDB int

func (typ TypeOfDB) GetType() string {
	return [...]string{"PostgresSQL", "SQL Lite"}[typ]
}

type DBParams struct {
	Name     string
	Type     int
	Host     string
	User     string
	Password string
	TimeZone string
	Port     int
	SslMode  bool
}

type DB struct {
	ID string

	Params DBParams
	Engine *gorm.DB
	Sender Sender.EmailSender
}

func GetDB(Params DBParams, sender Sender.EmailSender) (DB, error) {
	var err error
	db := DB{
		Params: Params,
		ID:     (uuid.New()).String(),
		Sender: sender,
	}
	if Params.Type == TypePostgres {
		err = db.InitPostgres()
	} else if Params.Type == TypeSQLLite {
		err = db.InitLite(Params.Name)
	} else {
		return DB{}, fmt.Errorf("wrong type of database")
	}
	if err = db.Migration(); err != nil {
		return DB{}, err
	}
	return db, err
}

func (db *DB) Migration() error {
	return db.Engine.AutoMigrate(&User{}, &RefreshToken{}, &Confirmation{}, &UserToUser{}, &BillingAccount{},
		&Subject{}, &SubjectToTeacher{}, &Lesson{}, &LessonRequest{}, &ChatRoom{}, &Teacher{}, &Transaction{})

}

func (db *DB) InitPostgres() error {
	sslMode := "disable"
	if db.Params.SslMode {
		sslMode = "enable"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		db.Params.Host, db.Params.User, db.Params.Password, db.Params.Name, db.Params.Port, sslMode, db.Params.TimeZone)
	var err error
	db.Engine, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return err
}

func (db *DB) InitLite(name string) error {
	database, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	db.Engine = database

	return err
}
