package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"wiki/dao"
	"wiki/middleware"
	"wiki/model"
)

//做个comment列的索引，直接拉索引crud

type DeleteCommentReq struct {
	CommentID    uint
	OperatorID   uint
	OperatorRole int
}

func AddComment(c model.Comment) {
	// 1. 类型转换：把数据库模型 c 转为 MQ 传输用的结构体
	// 假设你消费者那边解析的是 CurJsonToComment
	msgObj := middleware.CurJsonToComment{
		Userid:  c.UserID,  // 从 model.Comment 提取并转类型
		Pageid:  c.PageID,  // 从 model.Comment 提取并转类型
		Comment: c.Content, // 从 model.Comment 提取内容
	}

	// 2. 序列化为 JSON 字符串
	msgJson, err := json.Marshal(msgObj)
	if err != nil {
		log.Println("JSON 序列化失败:", err)
		return
	}

	// 3. 定义队列名并发送
	// 这里的 "comment_queue" 必须和你 ConsumeSimple 传入的参数一致
	const QueueName = "comment_queue"
	mq := middleware.NewRabbitMQSimple(QueueName)
	fmt.Println("mq start success")

	// 4. 调用你写好的发送方法
	mq.PublishSimple(string(msgJson))

	fmt.Printf("评论已入队: 用户%d 对 页面%d 的评论\n", msgObj.Userid, msgObj.Pageid)
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
