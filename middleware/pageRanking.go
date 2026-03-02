package middleware

import (
	"context"
	"log"
	"strconv"
	"time"
	"wiki/dao"

	"github.com/go-redis/redis/v8"
)

type PageRanking struct {
	Client *redis.Client   // Redis客户端
	ctx    context.Context // Redis操作的上下文
}

func NewPageRanking() *PageRanking {
	// 创建Redis客户端
	client := dao.InitRedis()
	log.Print(client)
	return &PageRanking{
		Client: client,
		ctx:    context.Background(),
	}
}

// IncreasePageVisit  添加日访问量排行榜
// IncreasePageVisit 词条访问量 +1
func (p *PageRanking) IncreasePageVisit(pageID int) error {
	// 建议超时时间设短一点，Redis 操作很快
	ctx, cancel := context.WithTimeout(p.ctx, time.Second*5)
	defer cancel()

	// ZIncrBy: 如果 daily_ranking 里没这个 pageID，会自动创建并设为 1；如果有了，就 +1
	err := p.Client.ZIncrBy(ctx, "daily_ranking", 1, strconv.Itoa(pageID)).Err()
	if err != nil {
		log.Printf("Redis 增加访问量失败: %v", err)
		return err
	}
	return nil
}

// GetDailyRanking 获取日访问量排行榜
func (p *PageRanking) GetDailyRanking(limit int64) ([]redis.Z, error) {
	// 设置超时时间
	ctx, cancel := context.WithTimeout(p.ctx, time.Hour*24)
	defer cancel()

	// 获取日访问量排行榜
	return p.Client.ZRevRangeWithScores(ctx, "daily_ranking", 0, limit).Result()
}
