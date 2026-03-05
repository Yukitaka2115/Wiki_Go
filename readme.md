# Wiki_Go

一个基于 Go 语言构建的百科/词条管理系统后端。本项目旨在通过实际业务场景（词条浏览、异步评论、热度排行）来实践微服务架构中常见的中间件解耦、缓存优化及大数据处理技术。

---

## 🛠️ 技术栈

* **核心开发**: [Go (Golang)](https://golang.org/)
* **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
* **数据库 (ORM)**: [GORM](https://gorm.io/) (MySQL)
* **消息队列**: [RabbitMQ](https://www.rabbitmq.com/) (实现评论异步写入)
* **缓存/排行榜**: [Redis](https://redis.io/) (ZSet 实现实时热度排行)
* **数据清洗**: [Apache Spark](https://spark.apache.org/) (计划中：用于非法内容过滤)

---

## ✨ 已完成功能

### 1. 异步评论系统 (RabbitMQ)
* **解耦设计**：用户提交评论后，系统通过生产者将消息推送到 RabbitMQ 队列，避免数据库写入延迟影响用户响应。
* **后台消费**：采用独立协程作为消费者，稳定地将消息从队列同步至 MySQL，有效缓解高并发压力。

### 2. 词条热度排行榜 (Redis ZSet)
* **实时计数**：用户每查看一次词条，触发异步 `ZIncrBy` 操作，实时累加词条访问量。
* **高效排序**：利用 Redis Sorted Set 结构，支持 `O(log(N))` 复杂度的排行榜查询。

### 3. 数据映射与解析
* **模型解耦**：使用 GORM 的标签系统（如 `gorm:"-"`）处理复杂的业务字段，解决 JSON 序列化字段与数据库列的映射冲突。
* **动态解析**：支持词条内多维信息（主要角色、相关人物、团队）的动态 JSON 转结构体解析。

---

## 📅 开发计划 (Roadmap)

- [x] RabbitMQ 生产者/消费者闭环
- [x] Redis 实时热度排行榜实现
- [x] GORM 数据库自动映射优化
- [ ] **Spark 评论清洗** (Next Step)：接入 Spark 任务，对评论内容进行敏感词扫描与清洗。
- [x] 用户权限验证 (JWT 机制)

---

## 🚀 快速开始

### 1. 环境依赖
* **Redis**: 端口 `6379`
* **RabbitMQ**: 端口 `5672`, 推荐使用虚拟主机 (Virtual Host) `/Cello`
* **MySQL**: 确保已创建相关数据库及表结构

### 2. 运行项目
```bash
# 克隆项目
git clone [https://github.com/Yukitaka2115/Wiki_Go.git](https://github.com/Yukitaka2115/Wiki_Go.git)

# 进入目录
cd Wiki_Go

# 启动服务
go run main.go
