package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wiki/service"
)

func Ranking(ctx *gin.Context) {
	ranking := service.NewPageRanking()
	dailyRanking, err := ranking.GetDailyRanking(-1)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"今日词条排行榜": dailyRanking})
}
