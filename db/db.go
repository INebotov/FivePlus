package db

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	Name     string `yaml:"Name"`
	Type     int    `yaml:"Type"`
	Host     string `yaml:"Host"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	TimeZone string `yaml:"TimeZone"`
	Port     int    `yaml:"Port"`
	SslMode  bool   `yaml:"SslMode"`
}

type DB struct {
	ID string

	Params DBParams
	Engine *gorm.DB

	Log *zap.Logger

	CloserFunc func() error
}

func (db *DB) Connect() error {
	var err error
	db.ID = (uuid.New()).String()

	if db.Params.Type == TypePostgres {
		err = db.InitPostgres()
		if err != nil {
			return err
		}
	} else if db.Params.Type == TypeSQLLite {
		err = db.InitLite(db.Params.Name)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("wrong type of database")
	}
	DATA, err := db.Engine.DB()
	if err != nil {
		return err
	}
	db.CloserFunc = DATA.Close
	return db.Migration()
}

func (db *DB) Migration() error {
	return db.Engine.AutoMigrate(&Session{})
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
