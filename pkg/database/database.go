package database

import (
	"fmt"

	"github.com/yakuter/ugin/model"
	"github.com/yakuter/ugin/pkg/config"
	"github.com/yakuter/ugin/pkg/logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	DB    *gorm.DB
	err   error
	DBErr error
)

type Database struct {
	*gorm.DB
}

// Setup opens a database and saves the reference to `Database` struct.
func Setup() error {
	var db = DB

	config := config.GetConfig()

	driver := config.Database.Driver
	database := config.Database.Dbname
	username := config.Database.Username
	password := config.Database.Password
	host := config.Database.Host
	port := config.Database.Port

	switch driver {
	case "sqlite":
		db, err = gorm.Open("sqlite3", fmt.Sprintf("%s", database))
	case "mysql":
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", username, password, host, port, database))
	case "postgres":
		db, err = gorm.Open("postgres", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", username, password, host, port, database))
	default:
		logger.Errorf("Database driver not found")
		return err
	}

	if err != nil {
		DBErr = err
		logger.Errorf("Failed to load database error: %v", err)
		return err
	}

	// Change this to true if you want to see SQL queries
	db.LogMode(config.Database.LogMode)

	// Auto migrate project models
	db.AutoMigrate(&model.Post{}, &model.Tag{})
	DB = db

	return nil
}

// GetDB helps you to get a connection
func GetDB() *gorm.DB {
	return DB
}

// GetDBErr helps you to get a connection
func GetDBErr() error {
	return DBErr
}
