package service

import (
	"encoding/json"
	"fmt"
)

type CurComment struct {
	Userid  int
	Pageid  int
	Comment string
}

func AddComment(uid int, pid int, comment string) {
	var commentStr = CurComment{
		Userid:  uid,
		Pageid:  pid,
		Comment: comment,
	}
	commentStrJson, _ := json.Marshal(commentStr)
	if commentStrJson == nil {
		fmt.Println("未转换成功")
		return
	} else {
		fmt.Println(commentStrJson)
	}
	mq := NewRabbitMQSimple("comments")
	mq.PublishSimple(string(commentStrJson))
	go mq.ConsumeSimple()
}
