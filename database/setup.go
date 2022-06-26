package database

import (
	"log"

	"github.com/OtchereDev/go-gorm-user-app/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	DB *gorm.DB
}

var AppDBInstance DBInstance

func ConnectDB() {

	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.Logger = logger.Default.LogMode(logger.Info)

	// AutoMigrate
	db.AutoMigrate(&models.UserModel{})

	AppDBInstance = DBInstance{DB: db}
}

func ConnectToModel(model interface{}) *gorm.DB {
	return AppDBInstance.DB.Model(model)
}
