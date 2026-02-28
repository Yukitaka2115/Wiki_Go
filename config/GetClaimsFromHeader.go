package config

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// GetClaimsFromHeader 统一从 Gin 的上下文中提取并解析 Token
func GetClaimsFromHeader(c *gin.Context) (map[string]interface{}, error) {
	// 1. 获取 Header 里的 Authorization 字段
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("请求头中缺失 Authorization")
	}

	// 2. 去掉 Bearer 前缀（如果有的话）
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	// 3. 调用你之前写好的 ExtractClaimsFromToken (解析 JWT)
	// 这里的 ExtractClaimsFromToken 应该也是在 config 包里定义的
	claims, err := ExtractClaimsFromToken(tokenString)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
