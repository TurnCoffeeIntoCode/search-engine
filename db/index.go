package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDB() {
	dburl := os.Getenv("DATABASE_URL")
	var err error
	DBConn, err = gorm.Open(postgres.Open(dburl))
	if err != nil {
		fmt.Println("Failed to connect to database")
		panic("Failed to connect to database")
	}

	// Enable uuid-ossp extension
	err = DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		fmt.Println("Failed to enable uuid-ossp extension")
		panic(err)
	}

	err = DBConn.AutoMigrate(&User{}, &SearchSettings{}, &CrawledUrl{}, &SearchIndex{})
	if err != nil {
		fmt.Println("Failed to migrate")
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return DBConn
}
