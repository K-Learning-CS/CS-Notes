# kafka

[官方文档](https://kafka.apache.org/documentation/#gettingStarted)

《深入理解Kafka：核心设计与实践原理》

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

- 参考官方文档

## 配置

