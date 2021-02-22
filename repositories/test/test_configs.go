package test

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func getDB() (*gorm.DB, error) {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		panic("NOT SET DSN")
	}
	dsn += "?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println(dsn)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
