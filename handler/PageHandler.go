package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"wiki/service"
)

func AddPage(ctx *gin.Context) {
	var page service.Page
	if err := ctx.BindJSON(&page); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service.AddPage(page)
	ctx.JSON(http.StatusOK, gin.H{"message": "添加成功", "data": page})
} //1

func GetAllPage(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	sizeStr := ctx.Query("size")

	// 将页码和每页大小转换为整数类型
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	pages, err := service.GetPagesWithPagination(page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// 将结果转换为 JSON 格式并返回给客户端
	ctx.JSON(http.StatusOK, pages)
} //1

func GetPageByTitle(ctx *gin.Context) {
	title := ctx.Param("title")
	page := service.GetPageByTitle(title)
	pageRanking := service.NewPageRanking()
	log.Println(page.ID)
	err := pageRanking.IncreasePageVisit(page.ID)
	if err == nil {
		log.Print("Page visit increased successfully")
	} else {
		log.Print("Failed to increase page visit:", err)
	}
	ctx.JSON(http.StatusOK, page)
} //1

func UpdatePageById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "无效的页面 ID")
		return
	}
	var newPage service.Page
	if err := ctx.ShouldBindJSON(&newPage); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, _ := service.UpdatePage(id, newPage)
	ctx.JSON(http.StatusOK, page)
} //1

func DeletePageByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "无效的页面 ID")
		return
	}
	service.DeletePageByID(id)
	ctx.JSON(http.StatusOK, id)
} //1
