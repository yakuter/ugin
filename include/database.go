package include

import (
	"fmt"
	"ugin/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB
var err error

// InitDB opens a database and saves the reference to `Database` struct.
func InitDB() *gorm.DB {
	var db = DB

	config := config.InitConfig()

	driver := config.Database.Driver
	database := config.Database.Dbname
	username := config.Database.Username
	password := config.Database.Password
	host := config.Database.Host
	port := config.Database.Port

	if driver == "sqlite" {

		db, err = gorm.Open("sqlite3", "./ugin.db")

		if err != nil {
			fmt.Println("db err: ", err)
		}

	} else if driver == "postgres" {

		db, err = gorm.Open("postgres", "host="+host+" port="+port+" user="+username+" dbname="+database+"  sslmode=disable password="+password)
		if err != nil {
			fmt.Println("db err: ", err)
		}

	} else if driver == "mysql" {

		db, err = gorm.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			fmt.Println("db err: ", err)
		}

	}

	DB = db

	return DB
}

// GetDB helps you to get a connection
func GetDB() *gorm.DB {
	return DB
}
