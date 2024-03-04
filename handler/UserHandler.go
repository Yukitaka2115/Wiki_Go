package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wiki/config"
	"wiki/service"
)

func Register(ctx *gin.Context) {
	var user service.User
	_ = ctx.ShouldBindJSON(&user)
	if user.Username == "" || user.Pwd == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "输入信息为空"})
	} //判空

	exists := service.IsUsernameExists(user.Username) //判断重复
	if exists != false {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "已存在用户名"})
	} else {
		user = service.User{
			Role:     user.Role,
			Username: user.Username,
			Pwd:      user.Pwd,
		}
		service.AddUser(user)

		//token, err := config.CreateToken(user.ID)
		//if err != nil {
		//	return
		//} //生成token
		ctx.JSON(http.StatusOK, gin.H{"user": user})
	}
} //注册不用生成token

func Login(ctx *gin.Context) {
	var user service.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Username == "" || user.Pwd == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "输入信息为空"})
		return
	}

	curUser := service.GetUserByUserNameAndPwd(user.Username, user.Pwd)
	token, err := config.CreateToken(curUser.ID, curUser.Role)
	if err != nil {
		return
	} //生成token
	ctx.JSON(http.StatusOK, gin.H{"currentUserToken": token})
}
