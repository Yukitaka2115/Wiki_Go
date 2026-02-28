package service

import (
	"errors"
	"fmt"
	"wiki/dao"
	"wiki/model"
)

//做个comment列的索引，直接拉索引crud

type DeleteCommentReq struct {
	CommentID    uint
	OperatorID   uint
	OperatorRole int
}

func AddComment(c model.Comment) {
	// 1. 直接把结构体 c 存入数据库
	// 这里的 c 必须是字段大写开头的结构体
	err := dao.Db.Table("comments").Create(&c).Error

	if err != nil {
		fmt.Println("写入失败详情:", err)
		return
	}
	fmt.Println("写入数据库成功！ID 为:", c.ID)

	// 2. 如果你还想发 MQ（可选）
	//commentStrJson, _ := json.Marshal(c)
	// ... 发送 MQ 的逻辑
}

// service/comment_service.go

// 1. 确保参数类型是 DeleteCommentReq
func DeleteComment(req DeleteCommentReq) error {
	var comment model.Comment

	// 2. 内部全部改用 req 里的字段，不要再出现 c.Param 或 c.JSON
	if err := dao.Db.First(&comment, req.CommentID).Error; err != nil {
		return errors.New("评论不存在")
	}

	isAdmin := req.OperatorRole == 1
	isOwner := uint(comment.UserID) == req.OperatorID

	if !isAdmin && !isOwner {
		return errors.New("权限不足")
	}

	return dao.Db.Delete(&comment).Error
}
