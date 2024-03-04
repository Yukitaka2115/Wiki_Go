package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wiki/config"
	"wiki/service"
)

func GetUsers(ctx *gin.Context) {
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
	fmt.Println("从声明中获取的角色信息：", claims["role"])
	role, ok := claims["role"].(float64)
	if !ok {
		fmt.Println("无法将角色信息转换为int类型")
		return
	}

	// 验证用户权限
	if int(role) != 1 {
		ctx.JSON(http.StatusForbidden, "暂无权限")
		return
	}

	// 获取用户信息
	var user []service.User

	users, err := service.GetUsers(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "获取用户信息失败")
		return
	}
	ctx.JSON(http.StatusOK, users)
} //获取全部用户信息

func DeleteUserByID(ctx *gin.Context) {
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
	role, ok := claims["role"].(float64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, "无效的用户角色信息")
		return
	}

	// 从 URL 参数中获取要删除的用户ID
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 验证用户权限
	if role != 0 && int(claims["id"].(float64)) != id {
		ctx.JSON(http.StatusForbidden, "没有权限删除其他用户信息")
		return
	}

	// 验证用户是否存在
	_, err = service.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, "用户不存在")
		return
	}

	// 执行删除操作
	err = service.DeleteUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "删除用户失败")
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("成功删除用户 ID: %d", id))
}

func UpdateUserByID(ctx *gin.Context) {
	var newUser service.User
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
	role, ok := claims["role"].(float64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, "无效的用户角色信息")
		return
	}

	// 从 URL 参数中获取要修改的用户ID
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 验证用户权限
	if role == 0 {
		// 如果角色为0，则具有修改任意用户的权限
		_, err := service.UpdateUserInfo(id, newUser)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, fmt.Sprintf("成功更新用户 ID 的信息: %d", id))
		return
	}

	// 如果角色不为0，则只能修改自身角色的用户
	userID, ok := claims["id"].(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, "无效的用户ID")
		return
	}

	if id != userID {
		ctx.JSON(http.StatusForbidden, "没有权限修改其他用户信息")
		return
	}
	_, err = service.UpdateUserInfo(id, newUser)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, fmt.Sprintf("成功修改用户 ID 的信息: %d", id))
}
