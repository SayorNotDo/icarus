# 测试完成时间预估
当前时间 +（单次测试用例执行耗时（update）* 剩余执行次数（update））

# 目录结构
models  模型存放
repository  数据库操作结构体存放
service 业务逻辑代码存放
controller  控制器存放

#热重启iris项目
rizla main

#zookeeper docker 镜像
zookeeper:latest
wurmeister/kafka

# 启动 zookeeper
docker -d -p 2181:2181 --name icarus_zookeeper zookeeper

# 启动kafka镜像生成容器
# KAFKA_BROKER_ID：配置broker_id，在kafka集群中是唯一的
# KAFKA_ZOOKEEPER_CONNECT：配置连接 zookeeper 的 ip 和 port
# KAFKA_LISTENERS：配置 kafka 监听的端口
docker run -d --name kafka1 -p 9092:9092 -e KAFKA_BROKER_ID=0 -e KAFKA_ZOOKEEPER_CONNECT=172.17.0.2:2181 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://127.0.0.1:9092 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 wurstmeister/kafka

cd /opt/kafka/bin

# --zookeeper：指定 zookeeper 容器的 ip 和 端口
# --topic：指定主题的名字
./kafka-topics.sh --create --zookeeper 172.17.0.2:2181 --replication-factor 1 --partitions 1 --topic test

生产者客户端
# --broker-list：指定 kafka 的 ip 和 port
# --topic：指定消息发送到哪个主题
./kafka-console-producer.sh --broker-list 127.0.0.1:9092 --topic test

消费者客户端
# --bootstrap-server：指定 kafka 的 ip 和 port
# --topic：指定订阅的主题
./kafka-console-consumer.sh --bootstrap-server 127.0.0.1:9092 --topic test
