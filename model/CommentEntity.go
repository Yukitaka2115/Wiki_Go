package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey"`
	UserID  int    `gorm:"column:user_id"` // 显式指定列名
	PageID  int    `gorm:"column:page_id"`
	Content string `gorm:"column:content"`
}
