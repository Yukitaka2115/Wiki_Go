package dao

import (
	"fmt"
	"gorm.io/driver/mysql" // V2 需要独立的驱动包
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // V2 的日志包
)

var Db *gorm.DB

func Init() {
	dsn := "root:114514@tcp(127.0.0.1:3306)/wiki?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	// V2 使用 mysql.Open(dsn)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 开启 SQL 日志打印（替代原来的 LogMode）
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	fmt.Println("GORM V2 connect success")
}
