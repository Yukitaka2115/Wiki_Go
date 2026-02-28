package main

import (
	"fmt"
	"wiki/config"
	"wiki/dao"
	"wiki/handler"
	"wiki/model"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	dao.Init()

	//fmt.Println("正在尝试连接数据库...")

	// 1. 强制打印当前连接的数据库名，看看是不是连错地方了
	var dbName string
	dao.Db.Raw("SELECT DATABASE()").Scan(&dbName)
	fmt.Println("当前真正连接的数据库是:", dbName)

	// 2. 强制建表并捕获错误
	err1 := dao.Db.AutoMigrate(model.Comment{})
	if err1 != nil {
		fmt.Printf("建表失败，报错内容: %v\n", err1)
	} else {
		fmt.Println("GORM 声称建表成功了！")
	}
	page := router.Group("/page")
	{
		page.POST("/add", handler.AddPage)
		authGroup := page.Group("/api/v1", config.AuthInterceptor())
		{
			authGroup.POST("/comment", handler.AddComment)
			authGroup.DELETE("/comment/:id", handler.DeleteCommentHandler)
		}
		page.GET("/all", handler.GetAllPage)
		page.GET("/:title", handler.GetPageByTitle)
		page.PUT("/update", handler.UpdatePageById)
		page.DELETE("/delete", handler.DeletePageByID)
		page.GET("/index", handler.Ranking)
	}
	user := router.Group("/user")
	{
		user.GET("/", handler.GetUsers)
		user.PUT("/update", handler.UpdateUserByID)
		user.DELETE("/delete", handler.DeleteUserByID)
		user.POST("/register", handler.Register)
		user.POST("/login", handler.Login)
	}

	err := router.Run(":8088")
	if err != nil {
		return
	}
}
