# RocketMQ

## 一.环境准备
*详细的部署方式及其优劣请参考官方文档*

|      ip       |     描述     | cpu | 内存 |  硬盘  | 硬盘 |   备注    |
|:-------------:|:----------:|-----|:--:|:----:|:--:|:-------:|
| 192.168.1.136 | rocketmq-0 | 4C  | 8G | 100G |    | slave-1 |
| 192.168.1.137 | rocketmq-1 | 4C  | 8G | 100G |    | slave-2 |
| 192.168.1.138 | rocketmq-2 | 4C  | 8G | 100G |    | slave-0 |

*从部署的角度需要考虑：*

* `local`模式，`local`模式与`cluster`模式的区别在于`proxy`组件是否独立部署，
这里为了简单内嵌了。
* 三节点三副本-同步双写集群，与异步复制区别在于`broker`的配置项`brokerRole`的值。
* 主备故障自动切换，5.0版本以后`nameserver`节点内嵌`controller`，需要修改配置并指定
启动命令后一并启动，也可以单独部署。由于其采用`Raft`协议所以至少需要三个节点。
* `acl`权限控制，只需要`broker`配置`aclEnable=true`便可以启用。权限控制配置
在`/usr/local/rocketmq/conf/plain_acl.yml`文件中，`/usr/local/rocketmq`为
`rocketmq`家目录。该配置在`broker`启用`acl`后为热加载模式，修改保存即生效。
* 图形界面，`rocketmq-dashboard`使用`docker`部署，其为可选项。

### 依赖安装
```bash
yum install -y java-1.8.0-openjdk-1.8.0.392.b08-2.el7_9.x86_64 unzip
```

### 二进制包下载
```bash
wget https://dist.apache.org/repos/dist/release/rocketmq/5.1.4/rocketmq-all-5.1.4-bin-release.zip
```

### 安装
```bash
1.解压并移动至 /usr/local/
unzip rocketmq-all-5.1.4-bin-release.zip

cp -r rocketmq-all-5.1.4-bin-release /usr/local/rocketmq

2.调整jvm内存配置，具体配置根据节点的内存大小来
vi /usr/local/rocketmq/bin/runserver.sh
#JAVA_OPT="${JAVA_OPT} -server -Xms1g -Xmx1g -Xmn256m

vi /usr/local/rocketmq/bin/runbroker.sh
#JAVA_OPT="${JAVA_OPT} -Xmn256m
#JAVA_OPT="${JAVA_OPT} -server -Xms2g -Xmx2g"
```

***


## 二、NameServer

* `nameserver`为无状态的注册中心，默认情况下无需任何配置，直接启动即可。这里的配置是为`controller`
准备的。
* `rocketmq`的日志目录为`/root/logs/rocketmqlogs/`，所有组件的日志都在此处，排查错误请在此
查看日志。
* 从现在开始需要在具体节点做具体配置，请注意其差异。

### rocketmq-0

1. 创建数据目录
```bash
mkdir -p /export/rocketmq/controller/DledgerController
```
2. `controller`配置

```bash
cat > /export/rocketmq/controller/namesrv-n0.conf <<EOF
#Namesrv config
listenPort = 9876
enableControllerInNamesrv = true

#controller config
controllerDLegerGroup = group1
controllerDLegerPeers = n0-192.168.1.136:9878;n1-192.168.1.137:9878;n2-192.168.1.138:9878
controllerDLegerSelfId = n0
controllerStorePath = /export/rocketmq/controller/DledgerController
EOF
```
3. 启动`nameserver`
```bash
nohup sh /usr/local/rocketmq/bin/mqnamesrv -c /export/rocketmq/controller/namesrv-n0.conf &
```


### rocketmq-1
1. 创建数据目录
```bash
mkdir -p /export/rocketmq/controller/DledgerController
```
2. `controller`配置

```bash
cat > /export/rocketmq/controller/namesrv-n1.conf <<EOF
#Namesrv config
listenPort = 9876
enableControllerInNamesrv = true

#controller config
controllerDLegerGroup = group1
controllerDLegerPeers = n0-192.168.1.136:9878;n1-192.168.1.137:9878;n2-192.168.1.138:9878
controllerDLegerSelfId = n1
controllerStorePath = /export/rocketmq/controller/DledgerController
EOF
```
3. 启动`nameserver`
```bash
nohup sh /usr/local/rocketmq/bin/mqnamesrv -c /export/rocketmq/controller/namesrv-n1.conf &
```

### rocketmq-2
1. 创建数据目录
```bash
mkdir -p /export/rocketmq/controller/DledgerController
```
2. `controller`配置

```bash
cat > /export/rocketmq/controller/namesrv-n2.conf <<EOF
#Namesrv config
listenPort = 9876
enableControllerInNamesrv = true

#controller config
controllerDLegerGroup = group1
controllerDLegerPeers = n0-192.168.1.136:9878;n1-192.168.1.137:9878;n2-192.168.1.138:9878
controllerDLegerSelfId = n2
controllerStorePath = /export/rocketmq/controller/DledgerController
EOF
```
3. 启动`nameserver`
```bash
nohup sh /usr/local/rocketmq/bin/mqnamesrv -c /export/rocketmq/controller/namesrv-n2.conf &
```

*参数解释：*

* enableControllerInNamesrv：Nameserver 中是否开启 controller，默认 false。
* controllerDLegerGroup：DLedger Raft Group 的名字，同一个 DLedger Raft Group 保持一致即可。
* controllerDLegerPeers：DLedger Group 内各节点的端口信息，同一个 Group 内的各个节点配置必须要保证一致。
* controllerDLegerSelfId：节点 id，必须属于 controllerDLegerPeers 中的一个；同 Group 内各个节点要唯一。
* controllerStorePath：controller 日志存储位置。controller 是有状态的，controller 重启或宕机需要依靠日志来恢复数据，该目录非常重要，不可以轻易删除。
* enableElectUncleanMaster：是否可以从 SyncStateSet 以外选举 Master，若为 true，可能会选取数据落后的副本作为 Master 而丢失消息，默认为 false。
* notifyBrokerRoleChanged：当 Broker 副本组上角色发生变化时是否主动通知，默认为 true。

***

## 三、Broker

*下列为各服务器上`broker`的配置和启动，从第一行注视开始为可选配置，
包括性能优化，数据盘定义，`controller`与`acl`的开启。*

*需要注意的是，必须先启动主库再启动从库，可以先配置完，再逐一启动。*

### rocketmq-0

1. 创建数据目录
```bash
mkdir -p /export/rocketmq/broker-0/{commitlog,consumequeue,index,conf,store}
mkdir -p /export/rocketmq/broker-1-s/{commitlog,consumequeue,index,conf,store}
```
2. `master`配置

```bash
cat > /export/rocketmq/broker-0/conf/broker-0.properties <<EOF
brokerClusterName=DefaultCluster
brokerName=broker-0
brokerId=0
deleteWhen=04
fileReservedTime=48
brokerRole=SYNC_MASTER
flushDiskType=ASYNC_FLUSH
namesrvAddr=192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876
brokerIP1=192.168.1.136
listenPort=10800
#commitLog每个文件的大小默认1G
mapedFileSizeCommitLog=1073741824
#ConsumeQueue每个文件默认存30W条，根据业务情况调整
mapedFileSizeConsumeQueue=300000
#destroyMapedFileIntervalForcibly=120000
#redeleteHangedFileInterval=120000
#检测物理文件磁盘空间
diskMaxUsedSpaceRatio=88
#存储路径
storePathRootDir=/export/rocketmq/broker-0/
#commitLog 存储路径
storePathCommitLog=/export/rocketmq/broker-0/commitlog
#消费队列存储路径存储路径
storePathConsumeQueue=/export/rocketmq/broker-0/consumequeue
#消息索引存储路径
storePathIndex=/export/rocketmq/broker-0/index
#checkpoint 文件存储路径
storeCheckpoint=/export/rocketmq/broker-0/checkpoint
#abort 文件存储路径
abortFile=/export/rocketmq/broker-0/abort
#限制的消息大小
maxMessageSize=65536
#发消息线程池数量
#sendMessageThreadPoolNums=128
#拉消息线程池数量
#pullMessaeThreadPoolNums=128

#发送消息是否使用可重入锁
useReentrantLockWhenPutMessage=true
waitTimeMillsInSendQueue=300  #或者更大

# controller
enableControllerMode = true
controllerAddr = 192.168.1.136:9878;192.168.1.137:9878;192.168.1.138:9878
storePathEpochFile = /export/rocketmq/broker-0/store

# acl
aclEnable=true
EOF
```

3. 启动`master`
```bash
nohup sh /usr/local/rocketmq/bin/mqbroker -n "192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876" -c /export/rocketmq/broker-0/conf/broker-0.properties --enable-proxy &
```


4. `slave`配置
```bash
cat > /export/rocketmq/broker-1-s/conf/broker-1-s.properties <<EOF
brokerClusterName=DefaultCluster
brokerName=broker-1
brokerId=1
deleteWhen=04
fileReservedTime=48
brokerRole=SLAVE
flushDiskType=ASYNC_FLUSH
namesrvAddr=192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876
listenPort=10900
#commitLog每个文件的大小默认1G
mapedFileSizeCommitLog=1073741824
#ConsumeQueue每个文件默认存30W条，根据业务情况调整
mapedFileSizeConsumeQueue=300000
#destroyMapedFileIntervalForcibly=120000
#redeleteHangedFileInterval=120000
#检测物理文件磁盘空间
diskMaxUsedSpaceRatio=88
#存储路径
storePathRootDir=/export/rocketmq/broker-1-s/
#commitLog 存储路径
storePathCommitLog=/export/rocketmq/broker-1-s/commitlog
#消费队列存储路径存储路径
storePathConsumeQueue=/export/rocketmq/broker-1-s/consumequeue
#消息索引存储路径
storePathIndex=/export/rocketmq/broker-1-s/index
#checkpoint 文件存储路径
storeCheckpoint=/export/rocketmq/broker-1-s/checkpoint
#abort 文件存储路径
abortFile=/export/rocketmq/broker-1-s/abort
#限制的消息大小
maxMessageSize=65536
#发消息线程池数量
#sendMessageThreadPoolNums=128
#拉消息线程池数量
#pullMessaeThreadPoolNums=128

#发送消息是否使用可重入锁
useReentrantLockWhenPutMessage=true
waitTimeMillsInSendQueue=300  #或者更大

# controller
enableControllerMode = true
controllerAddr = 192.168.1.136:9878;192.168.1.137:9878;192.168.1.138:9878
storePathEpochFile = /export/rocketmq/broker-1-s/store

# acl
aclEnable=true
EOF
```

5. `proxy`配置
```bash
cat > /export/rocketmq/proxyConfig.json <<EOF
{
  "rocketMQClusterName": "DefaultCluster",
  "grpcServerPort": 8091,
  "remotingListenPort": 8090
}
EOF
```

6. 启动`slave`
```bash
nohup sh /usr/local/rocketmq/bin/mqbroker -n "192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876" -pc /export/rocketmq/proxyConfig.json -c /export/rocketmq/broker-1-s/conf/broker-1-s.properties --enable-proxy &
```






### rocketmq-1

1. 创建数据目录
```bash
mkdir -p /export/rocketmq/broker-1/{commitlog,consumequeue,index,conf,store}
mkdir -p /export/rocketmq/broker-2-s/{commitlog,consumequeue,index,conf,store}
```
2. `master`配置

```bash
cat > /export/rocketmq/broker-1/conf/broker-1.properties <<EOF
brokerClusterName=DefaultCluster
brokerName=broker-1
brokerId=0
deleteWhen=04
fileReservedTime=48
brokerRole=SYNC_MASTER
flushDiskType=ASYNC_FLUSH
namesrvAddr=192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876
listenPort=10800
#commitLog每个文件的大小默认1G
mapedFileSizeCommitLog=1073741824
#ConsumeQueue每个文件默认存30W条，根据业务情况调整
mapedFileSizeConsumeQueue=300000
#destroyMapedFileIntervalForcibly=120000
#redeleteHangedFileInterval=120000
#检测物理文件磁盘空间
diskMaxUsedSpaceRatio=88
#存储路径
storePathRootDir=/export/rocketmq/broker-1/
#commitLog 存储路径
storePathCommitLog=/export/rocketmq/broker-1/commitlog
#消费队列存储路径存储路径
storePathConsumeQueue=/export/rocketmq/broker-1/consumequeue
#消息索引存储路径
storePathIndex=/export/rocketmq/broker-1/index
#checkpoint 文件存储路径
storeCheckpoint=/export/rocketmq/broker-1/checkpoint
#abort 文件存储路径
abortFile=/export/rocketmq/broker-1/abort
#限制的消息大小
maxMessageSize=65536
#发消息线程池数量
#sendMessageThreadPoolNums=128
#拉消息线程池数量
#pullMessaeThreadPoolNums=128

#发送消息是否使用可重入锁
useReentrantLockWhenPutMessage=true
waitTimeMillsInSendQueue=300  #或者更大

# controller
enableControllerMode = true
controllerAddr = 192.168.1.136:9878;192.168.1.137:9878;192.168.1.138:9878
storePathEpochFile = /export/rocketmq/broker-1/store

# acl
aclEnable=true
EOF
```

3. 启动`master`
```bash
nohup sh /usr/local/rocketmq/bin/mqbroker -n "192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876" -c /export/rocketmq/broker-2/conf/broker-2.properties --enable-proxy &
```


4. `slave`配置
```bash
cat > /export/rocketmq/broker-2-s/conf/broker-2-s.properties <<EOF
brokerClusterName=DefaultCluster
brokerName=broker-2
brokerId=1
deleteWhen=04
fileReservedTime=48
brokerRole=SLAVE
flushDiskType=ASYNC_FLUSH
namesrvAddr=192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876
listenPort=10900
#commitLog每个文件的大小默认1G
mapedFileSizeCommitLog=1073741824
#ConsumeQueue每个文件默认存30W条，根据业务情况调整
mapedFileSizeConsumeQueue=300000
#destroyMapedFileIntervalForcibly=120000
#redeleteHangedFileInterval=120000
#检测物理文件磁盘空间
diskMaxUsedSpaceRatio=88
#存储路径
storePathRootDir=/export/rocketmq/broker-2-s/
#commitLog 存储路径
storePathCommitLog=/export/rocketmq/broker-2-s/commitlog
#消费队列存储路径存储路径
storePathConsumeQueue=/export/rocketmq/broker-2-s/consumequeue
#消息索引存储路径
storePathIndex=/export/rocketmq/broker-2-s/index
#checkpoint 文件存储路径
storeCheckpoint=/export/rocketmq/broker-2-s/checkpoint
#abort 文件存储路径
abortFile=/export/rocketmq/broker-2-s/abort
#限制的消息大小
maxMessageSize=65536
#发消息线程池数量
#sendMessageThreadPoolNums=128
#拉消息线程池数量
#pullMessaeThreadPoolNums=128

#发送消息是否使用可重入锁
useReentrantLockWhenPutMessage=true
waitTimeMillsInSendQueue=300  #或者更大

# controller
enableControllerMode = true
controllerAddr = 192.168.1.136:9878;192.168.1.137:9878;192.168.1.138:9878
storePathEpochFile = /export/rocketmq/broker-2-s/store

# acl
aclEnable=true
EOF
```

5. `proxy`配置
```bash
cat > /export/rocketmq/proxyConfig.json <<EOF
{
  "rocketMQClusterName": "DefaultCluster",
  "grpcServerPort": 8091,
  "remotingListenPort": 8090
}
EOF
```

6. 启动`slave`
```bash
nohup sh /usr/local/rocketmq/bin/mqbroker -n "192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876" -pc /export/rocketmq/proxyConfig.json -c /export/rocketmq/broker-0-s/conf/broker-0-s.properties --enable-proxy &
```



### rocketmq-2

1. 创建数据目录
```bash
mkdir -p /export/rocketmq/broker-2/{commitlog,consumequeue,index,conf,store}
mkdir -p /export/rocketmq/broker-0-s/{commitlog,consumequeue,index,conf,store}
```
2. `master`配置

```bash
cat > /export/rocketmq/broker-2/conf/broker-2.properties <<EOF
brokerClusterName=DefaultCluster
brokerName=broker-2
brokerId=0
deleteWhen=04
fileReservedTime=48
brokerRole=SYNC_MASTER
flushDiskType=ASYNC_FLUSH
namesrvAddr=192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876
listenPort=10800
#commitLog每个文件的大小默认1G
mapedFileSizeCommitLog=1073741824
#ConsumeQueue每个文件默认存30W条，根据业务情况调整
mapedFileSizeConsumeQueue=300000
#destroyMapedFileIntervalForcibly=120000
#redeleteHangedFileInterval=120000
#检测物理文件磁盘空间
diskMaxUsedSpaceRatio=88
#存储路径
storePathRootDir=/export/rocketmq/broker-2/
#commitLog 存储路径
storePathCommitLog=/export/rocketmq/broker-2/commitlog
#消费队列存储路径存储路径
storePathConsumeQueue=/export/rocketmq/broker-2/consumequeue
#消息索引存储路径
storePathIndex=/export/rocketmq/broker-2/index
#checkpoint 文件存储路径
storeCheckpoint=/export/rocketmq/broker-2/checkpoint
#abort 文件存储路径
abortFile=/export/rocketmq/broker-2/abort
#限制的消息大小
maxMessageSize=65536
#发消息线程池数量
#sendMessageThreadPoolNums=128
#拉消息线程池数量
#pullMessaeThreadPoolNums=128

#发送消息是否使用可重入锁
useReentrantLockWhenPutMessage=true
waitTimeMillsInSendQueue=300  #或者更大

# controller
enableControllerMode = true
controllerAddr = 192.168.1.136:9878;192.168.1.137:9878;192.168.1.138:9878
storePathEpochFile = /export/rocketmq/broker-2/store

# acl
aclEnable=true
EOF
```

3. 启动`master`
```bash
nohup sh /usr/local/rocketmq/bin/mqbroker -n "192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876" -c /export/rocketmq/broker-1/conf/broker-1.properties --enable-proxy &
```

4. `slave`配置
```bash
cat > /export/rocketmq/broker-0-s/conf/broker-0-s.properties <<EOF
brokerClusterName=DefaultCluster
brokerName=broker-0
brokerId=1
deleteWhen=04
fileReservedTime=48
brokerRole=SLAVE
flushDiskType=ASYNC_FLUSH
namesrvAddr=192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876
listenPort=10900
#commitLog每个文件的大小默认1G
mapedFileSizeCommitLog=1073741824
#ConsumeQueue每个文件默认存30W条，根据业务情况调整
mapedFileSizeConsumeQueue=300000
#destroyMapedFileIntervalForcibly=120000
#redeleteHangedFileInterval=120000
#检测物理文件磁盘空间
diskMaxUsedSpaceRatio=88
#存储路径
storePathRootDir=/export/rocketmq/broker-0-s/
#commitLog 存储路径
storePathCommitLog=/export/rocketmq/broker-0-s/commitlog
#消费队列存储路径存储路径
storePathConsumeQueue=/export/rocketmq/broker-0-s/consumequeue
#消息索引存储路径
storePathIndex=/export/rocketmq/broker-0-s/index
#checkpoint 文件存储路径
storeCheckpoint=/export/rocketmq/broker-0-s/checkpoint
#abort 文件存储路径
abortFile=/export/rocketmq/broker-0-s/abort
#限制的消息大小
maxMessageSize=65536
#发消息线程池数量
#sendMessageThreadPoolNums=128
#拉消息线程池数量
#pullMessaeThreadPoolNums=128

#发送消息是否使用可重入锁
useReentrantLockWhenPutMessage=true
waitTimeMillsInSendQueue=300  #或者更大

# controller
enableControllerMode = true
controllerAddr = 192.168.1.136:9878;192.168.1.137:9878;192.168.1.138:9878
storePathEpochFile = /export/rocketmq/broker-0-s/store

# acl
aclEnable=true
EOF
```

5. `proxy`配置
```bash
cat > /export/rocketmq/proxyConfig.json <<EOF
{
  "rocketMQClusterName": "DefaultCluster",
  "grpcServerPort": 8091,
  "remotingListenPort": 8090
}
EOF
```

6. 启动`slave`
```bash
nohup sh /usr/local/rocketmq/bin/mqbroker -n "192.168.1.136:9876;192.168.1.137:9876;192.168.1.138:9876" -pc /export/rocketmq/proxyConfig.json -c /export/rocketmq/broker-2-s/conf/broker-2-s.properties --enable-proxy &
```

***

## 四、权限控制

* 上面的`broker`配置中已经启动了`acl`权限控制，这里修改会立刻生效，热加载。
* `console`用户为添加的新用户，给`dashboard`使用。同时注销了`IP`白名单。

```bash
cat > /usr/local/rocketmq/conf/plain_acl.yml <<EOF
globalWhiteRemoteAddresses:
#  - 10.10.103.*
#  - 192.168.*.*

accounts:
  - accessKey: RocketMQ
    secretKey: 12345678
    whiteRemoteAddress:
    admin: false
    defaultTopicPerm: DENY
    defaultGroupPerm: SUB
    topicPerms:
      - topicA=DENY
      - topicB=PUB|SUB
      - topicC=SUB
    groupPerms:
      # the group should convert to retry topic
      - groupA=DENY
      - groupB=PUB|SUB
      - groupC=SUB

  - accessKey: rocketmq2
    secretKey: 12345678
    whiteRemoteAddress:
    # if it is admin, it could access all resources
    admin: true

  - accessKey: console
    secretKey: 12345678
    whiteRemoteAddress:
    admin: true
EOF
```

***

## 五、图形化界面

*使用`<ip>:9090`访问`dashboard`，账号`admin`密码`admin`*

1. 创建数据目录
```bash
mkdir -p /export/rocketmq/console/data

chmod -R 777 /export/rocketmq/console/data
```

2. 添加`dashboard`访问用户
```bash
cat > /export/rocketmq/console/data/users.properties <<EOF
# This file supports hot change, any change will be auto-reloaded without Dashboard restarting.
# Format: a user per line, username=password[,N] #N is optional, 0 (Normal User); 1 (Admin)

# Define Admin
admin=admin,1

# Define Users
#user1=user1
#user2=user2
EOF
```

3. 启动`dashboard`
```bash
docker run -d --name rocketmq-dashboard --restart=always \
   -v /export/rocketmq/console/data:/tmp/rocketmq-console/data \
   -e "JAVA_OPTS=-Drocketmq.namesrv.addr=192.168.1.136:9876 -Dcom.rocketmq.sendMessageWithVIPChannel=false -Drocketmq.config.loginRequired=true -Drocketmq.config.accessKey=console -Drocketmq.config.secretKey=12345678" \
   -p 9090:8080 \
   -t apacherocketmq/rocketmq-dashboard:latest
```

## 六、参考链接

[RocketMQ的Docker镜像部署](https://blog.csdn.net/alionsss/article/details/135139421)

[部署方式](https://rocketmq.apache.org/zh/docs/deploymentOperations/01deploy)

[主备自动切换模式部署](https://rocketmq.apache.org/zh/docs/deploymentOperations/03autofailover)

[权限控制](https://rocketmq.apache.org/zh/docs/bestPractice/03access)

[图形化界面](https://rocketmq.apache.org/zh/docs/deploymentOperations/04Dashboard)

[proxy端口冲突](https://blog.csdn.net/Cooder_SXK/article/details/132482516)