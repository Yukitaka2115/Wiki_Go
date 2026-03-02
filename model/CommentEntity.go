package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	// 注意：这里的 json 标签必须和你发送的测试用例 Key 大小写完全一致
	UserID  int    `json:"UserID"`
	PageID  int    `json:"PageID"`
	Content string `json:"Content"`
	//UserName string `json:"UserName"` // 如果不需要传这个，可以不写
}
