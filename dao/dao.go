package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(mysql.Open(getDbDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func getDbDSN() string {
	dsn := os.Getenv("DSN")
	dsn += "?charset=utf8mb4&parseTime=True&loc=Local"
	return dsn
}
