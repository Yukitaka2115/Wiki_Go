package handler

import (
	"net/http"
	"wiki/middleware"

	"github.com/gin-gonic/gin"
)

func Ranking(ctx *gin.Context) {
	ranking := middleware.NewPageRanking()
	dailyRanking, err := ranking.GetDailyRanking(-1)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"今日词条排行榜": dailyRanking})
}
