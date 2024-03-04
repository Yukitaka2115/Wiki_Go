package service

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"wiki/dao"
)

// MQURL 连接信息amqp://Cello:114514@127.0.0.1:5672/Cello://事固定参数后面两个是用户名密码ip地址端口号Virtual Host
const MQURL = "amqp://Cello:114514@127.0.0.1:5672/Cello"

// RabbitMQ rabbitMQ结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	MqUrl string
}

type CurJsonToComment struct {
	Userid  int    `json:"Userid"`
	Pageid  int    `json:"Pageid"`
	Comment string `json:"Comment"`
}

// NewRabbitMQ 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, MqUrl: MQURL}
}

// Destroy 断开channel 和 connection
func (r *RabbitMQ) Destroy() {
	err := r.channel.Close()
	if err != nil {
		return
	}
	err = r.conn.Close()
	if err != nil {
		return
	}
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// NewRabbitMQSimple 创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// PublishSimple 直接模式队列生产
func (r *RabbitMQ) PublishSimple(message string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, _ = r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if r.QueueName != "" {
		fmt.Println(r.QueueName)
	}
	//调用channel 发送消息到队列中
	msgQueue := r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routeKey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if msgQueue != nil {
		fmt.Println(msgQueue)
		fmt.Println(message)
	} else {
		fmt.Println("no msg in queue")
	}
}

// ConsumeSimple simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if r.QueueName != "" {
		fmt.Println(r.QueueName)
	}
	if err != nil {
		fmt.Println(err)
	}

	//接收消息
	msgs, err := r.channel.Consume(
		q.Name, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		true, // auto-ack
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Connection中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if msgs != nil {
		fmt.Println(msgs)
	} else {
		fmt.Println("no msgs received")
	}
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var clearJson CurJsonToComment
			var user User

			fmt.Println(string(d.Body))

			// 解析消息
			err := json.Unmarshal(d.Body, &clearJson)
			if err != nil {
				return
			}
			log.Println("Received message:", clearJson)

			// 查询用户信息
			err = dao.Db.Find(&user, clearJson.Userid).Error
			if err != nil {
				fmt.Println("Failed to find user:", err)
				continue
			}

			// 构建评论对象
			comment := Comment{
				Username: user.Username,
				Comment:  clearJson.Comment,
			}
			log.Println("comments", comment)

			// 将评论插入到数据库
			log.Println(clearJson.Pageid)
			commentJson, _ := json.Marshal(comment)
			log.Println(string(commentJson))
			tx := dao.Db.Begin()
			tx.Model(&Page{}).Where("id = ?", clearJson.Pageid).Update("comments", commentJson)
			tx.Commit()
			if err != nil {
				fmt.Println("Failed to insert comment to database:", err)
				continue
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
