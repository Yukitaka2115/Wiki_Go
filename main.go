package main

import (
	"github.com/gin-gonic/gin"
	"wiki/handler"
)

func main() {
	router := gin.Default()
	page := router.Group("/page")
	{
		page.POST("/add", handler.AddPage)
		page.POST("/addComment", handler.AddComment)
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
