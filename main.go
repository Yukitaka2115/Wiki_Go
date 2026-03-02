package main

import (
	"wiki/dao"
	"wiki/handler"
	"wiki/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	mq := middleware.NewRabbitMQSimple("comment_queue")
	go mq.ConsumeSimple() // ✅ 必须加 go，让它在后台跑

	router := gin.Default()
	dao.Init()
	page := router.Group("/page")
	{
		authGroup := page.Group("", middleware.AuthInterceptor())
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
