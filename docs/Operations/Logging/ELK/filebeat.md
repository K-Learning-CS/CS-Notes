# filebeat

[官方文档](https://www.elastic.co/guide/en/beats/filebeat/current/index.html)


## 名词

- inputs       输入，获取日志的方式，日志发现配置
- harvesters   收割机，为每个日志输入启用一个收割机，日志读取跟踪组件
- output       输出，输出日志的方式
- filestream   文件流，input中的一种
- parsers      转换，这里使用它下面的multiline做多行日志处理
- processors   处理器，对经过的日志进行流式处理


## 原理 

Filebeat 由两个主要组件组成：

- 输入（inputs）：输入负责管理收集器并查找所有要读取的来源。每个输入都在其自己的 Go routine中运行。Filebeat 目前支持多种input类型。每种输入类型可以定义多次。log输入会检查每个文件，以查看是否需要启动收割机、收割机是否已在运行或是否可以忽略该文件（请参阅ignore_older）。仅当收割机关闭后文件大小发生变化时，才会拾取新行。
- 收割机（harvesters）：收割机负责读取单个文件的内容。收割机逐行读取每个文件，并将内容发送到输出。每个文件都会启动一个收割机。收割机负责打开和关闭文件，这意味着文件描述符在收割机运行时保持打开状态。如果在收割过程中删除或重命名文件，Filebeat 会继续读取该文件。这样做的副作用是，磁盘上的空间会被保留，直到收割机关闭。默认情况下，Filebeat 会保持文件打开，直到close_inactive达到。



## 安装

- 参考官方文档

## 配置

细节参考官方文档，有三种配置的方式：

- 数据收集模块——简化常见日志格式的收集、解析和可视化
- ECS 记录器 — 将应用程序日志结构化并格式化为与 ECS 兼容的 JSON
- 手动 Filebeat 配置

### 手动配置

```yaml
filebeat:
  inputs:
    - type: filestream # 文件流类型
      id: applog # 每个文件流都需要有唯一id
      paths:
        - /logs/app_log/*.log
      ignore_older: 1h #忽略在指定时间段之前修改的任何文件
      close: # 关闭收割机相关的选项
        on_state_change: # 在状态改变时
          renamed: true
      fields_under_root: true #  自定义字段将作为顶级字段存储在输出文档中，而不是分组到fields子词典下
      fields: # 自定义字段
        type: hospital_prd_applog
        format_tag: applog
#      scan_frequency: 2s
      parsers: # 日志行必须经过的转换列表
        - multiline: # 多行控制
            type: pattern
            pattern: '^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}'
            negate: true
            match: after

output:
    kafka:
      enabled: true # 启用输出
      # username: admin
      # password: admin
      hosts: ["kafka.logging.svc.cluster.local:9092"] # kafka地址
      topic: "%{[type]}" # 用type字段的值作为topic
      version: "2.0.0" # kafka版本
      max_retries: 3
      bulk_max_size: 2048 # 单个 Kafka 请求中批量处理的最大事件数
      timeout: 30s
      broker_timeout: 10s
      channel_buffer_size: 256
      keep_alive: 60
      compression: gzip
      max_message_bytes: 1000000
      required_acks: 1


logging: # filebeat的日志配置
  level: info
  to_files: true
  files:
    path: /tmp/
    rotateeverybytes: 10485760 # = 10MB
    name: filebeat
    keepfiles: 7
    permissions: 0644

http: # filebeat的http页面
  enabled: true
  host: 0.0.0.0
  port: 5678
```

*使用此配置，filebeat 将为`/logs/app_log/`目录中所有以`.log`结尾的文件启动一个收割机，为匹配的所有行都加上两个自定义字段。在输出部分 filebeat 将以`hospital_prd_applog`作为`topic`向 kafka 队列中输出数据，不存在则创建。*