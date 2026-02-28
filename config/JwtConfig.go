package config

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func CreateToken(id int, role int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte("yukitaka2115")) //私钥
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractClaimsFromToken(tokenString string) (jwt.MapClaims, error) {
	// 解析 JWT Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("yukitaka2115"), nil
	})
	if err != nil {
		return nil, err
	}

	// 验证 JWT Token 是否有效
	if !token.Valid {
		return nil, errors.New("无效的身份验证信息")
	}

	// 从 JWT Token 中获取用户身份信息
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("无效的身份验证信息")
	}

	return claims, nil
}
