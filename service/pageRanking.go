package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
	"wiki/dao"
)

type PageRanking struct {
	client *redis.Client   // Redis客户端
	ctx    context.Context // Redis操作的上下文
}

func NewPageRanking() *PageRanking {
	// 创建Redis客户端
	client := dao.InitRedis()
	log.Print(client)
	return &PageRanking{
		client: client,
		ctx:    context.Background(),
	}
}

// IncreasePageVisit  添加日访问量排行榜
func (p *PageRanking) IncreasePageVisit(pageID int) error {
	// 设置超时时间
	ctx, cancel := context.WithTimeout(p.ctx, time.Hour*24)
	defer cancel()

	// 增加页面访问量
	_, err := p.client.ZIncrBy(ctx, "daily_ranking", 1, strconv.Itoa(pageID)).Result()

	dailyRanking, err := p.GetDailyRanking(-1)
	if err != nil {
		log.Println("Failed to get daily ranking:", err)
	} else {
		log.Println("Daily ranking:", dailyRanking)
	}
	return nil
}

// GetDailyRanking 获取日访问量排行榜
func (p *PageRanking) GetDailyRanking(limit int64) ([]redis.Z, error) {
	// 设置超时时间
	ctx, cancel := context.WithTimeout(p.ctx, time.Hour*24)
	defer cancel()

	// 获取日访问量排行榜
	return p.client.ZRevRangeWithScores(ctx, "daily_ranking", 0, limit).Result()
}
