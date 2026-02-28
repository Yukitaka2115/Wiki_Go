package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 调用你刚摘出来的公共方法
		claims, err := GetClaimsFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或凭证无效"})
			c.Abort() // 必须调用 Abort，否则后面的 Handler 还会执行
			return
		}

		// 2. 将解析出的常用字段直接塞进上下文，方便后续直接 c.Get
		// 注意类型转换，JWT 解析出来的数字通常是 float64
		if uid, ok := claims["user_id"].(float64); ok {
			c.Set("userID", uint(uid))
		}
		if role, ok := claims["role"].(float64); ok {
			c.Set("userRole", int(role))
		}

		c.Next()
	}
}
