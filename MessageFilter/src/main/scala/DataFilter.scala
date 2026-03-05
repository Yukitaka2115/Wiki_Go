import com.rabbitmq.client.{ConnectionFactory, DefaultConsumer, Envelope => RabbitEnvelope, AMQP}
import org.json4s._
import org.json4s.native.JsonMethods._

// 对应你 Go 里的 CurJsonToComment
case class CommentMsg(Userid: Int, Pageid: Int, Comment: String)

object SparkCommentCleaner {
  implicit val formats: Formats = DefaultFormats

  def main(args: Array[String]): Unit = {
    val factory = new ConnectionFactory()
    factory.setHost("host.docker.internal") // 对应你 Go 里的地址
    factory.setPort(5672)
    factory.setUsername("Cello")
    factory.setPassword("114514")
    factory.setVirtualHost("/Cello") // 注意 Java client 的 vhost 以 / 开头

    val conn = factory.newConnection()
    val rawChannel = conn.createChannel()
    val cleanChannel = conn.createChannel()

    // 声明两个队列：生产者写入 mq_raw，清洗后写入 mq_clean，后端消费者从 mq_clean 读取
    rawChannel.queueDeclare("mq_raw", false, false, false, null)
    cleanChannel.queueDeclare("mq_clean", false, false, false, null)

    println(" [*] Spark 过滤器已启动，正在监听 mq_raw...")

    val consumer = new DefaultConsumer(rawChannel) {
      override def handleDelivery(consumerTag: String, envelope: RabbitEnvelope, properties: AMQP.BasicProperties, body: Array[Byte]): Unit = {
        val message = new String(body, "UTF-8")

        try {
          // 1. 解析 JSON
          val json = parse(message)
          val comment = json.extract[CommentMsg]

          // 2. 清洗逻辑
          if (isSafe(comment.Comment)) {
            // 3. 如果合法，转发到 mq_clean
            cleanChannel.basicPublish("", "mq_clean", null, message.getBytes("UTF-8"))
            println(s"[通过] 用户 ${comment.Userid} 的评论已清洗通过")
          } else {
            // 如果非法，可以选择落盘、告警或写入专门的拒绝队列
            println(s"[拦截] 发现非法评论: ${comment.Comment}")
          }
        } catch {
          case e: Exception => println(s"解析失败: ${e.getMessage}")
        }
      }
    }

    rawChannel.basicConsume("mq_raw", true, consumer)

    // 阻塞主线程以保持连接（简单实现）
    while (true) {
      Thread.sleep(10000)
    }
  }

  // 简单的非法词过滤逻辑
  def isSafe(content: String): Boolean = {
    val blackList = List("非法", "广告", "刷屏")
    !blackList.exists(content.contains)
  }
}