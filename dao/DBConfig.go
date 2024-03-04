package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func init() {
	dsn := "root:114514@tcp(127.0.0.1:3306)/wiki?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db, err = gorm.Open("mysql", dsn)
	Db.LogMode(true)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
}
