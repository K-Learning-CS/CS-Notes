# kafka

version: 3.8

reference:

- [官方文档](https://kafka.apache.org/documentation/#gettingStarted)
- 《深入理解Kafka：核心设计与实践原理》

## 名词

- producer  生产者，生产和发送消息方
- broker    代理，接受生产者的消息，等待消费者消费，通常说的kafka指的就是broker
- consumer  消费者，broker的消费方
- topic     主题，所有的消息都有对应的主题，生产者与消费者通过主题进行消息的传递
- partition 分区，topic是逻辑上的概念，真正的消息队列是分区，一个topic下面会有多个分区
- replica   副本，kafka的容灾方式与mysql类似，一主多从，主副本读写，从副本仅同步，副本的粒度是分区


## 原理

### 机制

- leader   主副本
- follower 从副本
- AR（Assigned Replicas）   所有副本
- ISR（In-Sync Replicas）   在同步副本
- OSR（Out-of-SyncReplicas）未同步副本
- HW（High Watermark）      高水位
- LEO（Log End Offset）     最后一条日志的偏移量

*主副本负责维护在同步和未同步副本集合，只有在ISR中的副本才有资格被选为主副本。HW是一个虚拟指针，决定了用户能够访问消息的偏移量，它指向了所有副本都拥有消息的偏移量，既延迟最大的副本的偏移量，当所有副本都与主副本同步时么，HW=LEO。*

## 安装

*对于本地化部署建议直接使用服务器，而不是k8s集群，理由同mysql一样，图方便可以使用docker部署。*

- 容器部署参考官方文档：https://github.com/apache/kafka/blob/trunk/docker/examples/docker-compose-files/cluster/combined/plaintext/docker-compose.yml


## 配置

*对于容器化部署来说，有两种自定义配置的方式，挂载文件/使用环境变量，下面是使用环境变量的方式。*

- 配置部分的官方文档：https://github.com/apache/kafka/blob/trunk/docker/examples/README.md

```shell
# zookeeper地址
KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181
# 集群节点id，必须不同
KAFKA_NODE_ID=1
# 节点角色，这种方式是联合部署，可以将两种角色分开部署
KAFKA_PROCESS_ROLES='broker,controller'
# 监听协议列表，目前使用未开启安全验证的协议PLAINTEXT
KAFKA_LISTENER_SECURITY_PROTOCOL_MAP='CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT'
# 投票节点
KAFKA_CONTROLLER_QUORUM_VOTERS='1@kafka-1:9093,2@kafka-2:9093,3@kafka-3:9093'
# 监听地址，供内部访问
KAFKA_LISTENERS='PLAINTEXT://:19092,CONTROLLER://:9093,PLAINTEXT_HOST://:9092'
KAFKA_INTER_BROKER_LISTENER_NAME='PLAINTEXT'
# 监听地址，供外部访问
KAFKA_ADVERTISED_LISTENERS='PLAINTEXT://kafka.logging.svc.cluster.local:9092'
KAFKA_CONTROLLER_LISTENER_NAMES='CONTROLLER'
# 集群id，同一集群所有节点必须相同
CLUSTER_ID='4L6g3nShT-eMCtK--X86sw'
# 分区副本数量，副本数量需要根据节点数调整，如果副本数>节点数，那么未分配的副本将导致主题创建失败
KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1
# 新消费者组创建时，组协调器等待更多消费者加入该组的时间。在执行第一次再平衡之前，协调器会延迟指定的时间。
KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0
# 事务日志主题在写入时，必须有至少该数量的副本确认写入成功，才会被视为成功。
KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1
# Kafka 需要将事务相关的数据复制到多少个节点，以确保数据的可用性和容错能力。
KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1
# 持久化日志目录
KAFKA_LOG_DIRS='/tmp/kraft-combined-logs'
```

