package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wiki/config"
	"wiki/service"
)

type Comment struct {
	Comment string `json:"comment"`
}

func AddComment(ctx *gin.Context) {
	var comment Comment
	// 从请求头中获取 JWT Token
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, "未提供身份验证信息")
		return
	}

	// 从 JWT Token 中获取用户身份信息
	claims, err := config.ExtractClaimsFromToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "身份验证失败")
		return
	}

	// 从声明中获取用户角色信息
	uidFloat64, ok := claims["user_id"].(float64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, "无效的用户角色信息")
		return
	}

	// 转换为int类型
	uid := int(uidFloat64)
	pid, err := strconv.Atoi(ctx.Query("id"))
	err = ctx.BindJSON(&comment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 将评论传递给service.AddComment函数
	service.AddComment(uid, pid, comment.Comment)
	ctx.JSON(http.StatusOK, gin.H{"添加评论成功": comment.Comment})
}
