package dao

import (
	"github.com/kabukky/journey/repositories/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDao() {
	var err error
	dsn, err := setting.GetGlobal("dsn")
	if err != nil {
		panic("Get DB DSN error: " + err.Error())
	}
	dsnStr := dsn.GetString() + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsnStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
