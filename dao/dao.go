package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDao(dsn string) {
	var err error

	dsnStr := dsn + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsnStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
