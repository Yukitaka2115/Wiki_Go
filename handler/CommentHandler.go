package handler

import (
	"net/http"
	"strconv"
	"wiki/model"
	"wiki/service"

	"github.com/gin-gonic/gin"
)

type Comment struct {
	Comment string `json:"comment"`
}

func AddComment(ctx *gin.Context) {
	// 连获取 Token 的逻辑都不用写了，直接从 Set 好的地方取
	uid, _ := ctx.Get("userID")

	var comment model.Comment
	// ... 绑定 JSON ...
	comment.UserID = int(uid.(uint))
	service.AddComment(comment)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": comment, // 直接返回结构体，前端就能看到绑定的 pageid 了
	})
}

func DeleteCommentHandler(ctx *gin.Context) {
	// 1. 从路由获取评论 ID (例如 /api/comment/:id)
	idStr := ctx.Param("id")
	commentID, _ := strconv.Atoi(idStr)

	// 2. 直接从中间件 Set 好的地方取（注意类型断言）
	// 此时不需要再解析 Token，中间件没通过的话根本进不来这个函数
	uid, _ := ctx.Get("userID")
	role, _ := ctx.Get("userRole")

	// 3. 组装请求对象（还记得我们之前优化的 Request 结构体吗？）
	req := service.DeleteCommentReq{
		CommentID:    uint(commentID),
		OperatorID:   uid.(uint),
		OperatorRole: role.(int),
	}

	// 4. 调用 Service 执行带权限校验的删除
	if err := service.DeleteComment(req); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
