package middleware

import (
	"encoding/json"
	"log"
	"wiki/dao"
	"wiki/model"

	"github.com/streadway/amqp"
)

// MQURL 注意：如果没创建过名为 Cello 的 Virtual Host，请把最后的 /Cello 改成 /
const MQURL = "amqp://Cello:114514@127.0.0.1:5672/Cello"

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	MqUrl     string
}

type CurJsonToComment struct {
	Userid  int    `json:"Userid"`
	Pageid  int    `json:"Pageid"`
	Comment string `json:"Comment"`
}

// NewRabbitMQSimple 简单模式：内部处理所有连接逻辑
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		MqUrl:     MQURL,
	}
	var err error
	// 1. 建立 TCP 连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	if err != nil {
		log.Printf("【错误】无法连接 RabbitMQ: %v", err)
		return nil
	}

	// 2. 建立信道 (Channel)
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	if err != nil {
		log.Printf("【错误】无法打开 Channel: %v", err)
		return nil
	}
	return rabbitmq
}

// Destroy 只有在程序彻底退出时才调用
func (r *RabbitMQ) Destroy() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// PublishSimple 生产者：发送评论消息
func (r *RabbitMQ) PublishSimple(message string) {
	if r.channel == nil {
		log.Println("【错误】Channel 未初始化，无法发送")
		return
	}

	// 1. 声明队列（确保队列存在）
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否排他
		false, // 是否阻塞
		nil,   // 额外参数
	)
	if err != nil {
		log.Printf("【错误】队列声明失败: %v", err)
		return
	}

	// 2. 发送消息 注意：Publish 返回的是 error，不是队列对象！
	err = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	if err != nil {
		log.Printf("【错误】消息发送失败: %v", err)
	} else {
		log.Printf("【成功】已发送消息: %s", message)
	}
}

// ConsumeSimple 消费者：持续监听并写入数据库
func (r *RabbitMQ) ConsumeSimple() {
	// 1. 声明队列（确保队列存在）
	q, err := r.channel.QueueDeclare(
		r.QueueName, false, false, false, false, nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 获取消息流
	msgs, err := r.channel.Consume(
		q.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		log.Printf(" [*] 评论接收器已启动，正在等待消息...")

		for d := range msgs {
			// d.Body 就是你截图里看到的那个 JSON 字符串
			log.Printf("接收到原始消息: %s", d.Body)

			// 3.1 解析 JSON
			var clearJson CurJsonToComment
			err := json.Unmarshal(d.Body, &clearJson)
			if err != nil {
				log.Printf("JSON 解析失败: %v", err)
				continue // 解析失败跳过本条，继续看下一条
			}

			// 3.2 构造数据库模型
			// 假设你的 model.Comment 结构体包含这些字段
			newComment := model.Comment{
				UserID:  clearJson.Userid,
				PageID:  clearJson.Pageid,
				Content: clearJson.Comment,
			}

			// 3.3 写入数据库
			// 使用你项目中已有的 dao.Db
			err = dao.Db.Create(&newComment).Error
			if err != nil {
				log.Printf("评论存入数据库失败: %v", err)
				continue
			}

			log.Printf("【成功】评论已写入数据库: 用户ID %d, 页面ID %d", clearJson.Userid, clearJson.Pageid)
		}
		currentQ, err := r.channel.QueueDeclarePassive(
			r.QueueName, false, false, false, false, nil,
		)
		if err == nil {
			log.Printf("【状态报告】消息处理完成。当前队列剩余消息数: %d", currentQ.Messages)
		}
	}()

	log.Printf(" [*] 消费者已启动...")
	<-forever
}
