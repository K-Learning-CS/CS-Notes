### NoSQL介绍
1.NoSQL 简介
    NoSQL(NoSQL = Not Only SQL )，意即"不仅仅是SQL"，指的是非关系型的数据库，是对不同于传统的关系型数据库的数据库管理系统的统称。
    
    在现代的计算系统上每天网络上都会产生庞大的数据量。
    这些数据有很大一部分是由关系数据库管理系统（RDMBSs）来处理，也有一部分使用非系型数据库处理
    
    对NoSQL最普遍的解释是"非关联型的"，强调Key-Value Stores和文档数据库的优点，而不是单纯的反对RDBMS。
    
    NoSQL用于超大规模数据的存储。（例如谷歌或Facebook每天为他们的用户收集万亿比特的数据）。这些类型的数据存储不需要固定的模式，无需多余操作就可以横向扩展。
2.为什么使用NoSQL
    关系型数据库对数据要求严格，而非关系型数据库没有那么严格，对于大量不同字段的数据，存储更加方便
### MongoDB简介
    Mongodb由C++语言编写的，是一个基于分布式文件存储的开源数据库系统。是专为可扩展性，高性能和高可用性而设计的数据库， 是非关系型数据库中功能最丰富，最像关系型数据库的，它支持的数据结构非常散，是类似 json 的 bjson 格式，因此可以存储比较复杂的数据类型。
    
    MongoDB的（来自于英文单词“Humongous”，中文含义为“庞大”）是可以应用于各种规模的企业，各个行业以及各类应用程序的开源数据库。作为一个适用于敏捷开发的数据库，MongoDB的的数据模式可以随着应用程序的发展而灵活地更新。
    
    MongoDB 以一种叫做 BSON（二进制 JSON）的存储形式将数据作为文档存储。具有相似结构的文档通常被整理成集合。可以把这些集合看成类似于关系数据库中的表： 文档和行相似， 字段和列相似。
    
    json格式：{key:value,key:value}
    bjson格式：{key:value,key:value}
    #区别在于：对于数据{id:1}，在JSON的存储上1只使用了一个字节，而如果用BJSON，那就是至少4个字节
1.MySQL与mongoDB对比
1）结构对比
mysql	MongoDB
库	库
表	集合
字段	键值
行	文档
1）数据库中数据（student库，user表）
uid	name	age
1	zhangyu	18
2	chencgheng	28
2）mongoDB中的数据（student库，user集合）
    1) {uid:1,name:zhangyu,age:18}
    2) {uid:2,name:chencgheng,age:28}
3）区别总结
    1.数据结构不同
    2.数据库添加不存在字段的数据时报错
    3.mongoDB可以添加不存在的字段的数据
    4.mongoDB不需要提前创建好库和表，创建数据直接会帮助我们创建好
2.MongoDB 特点
    1.高性能：
      Mongodb 提供高性能的数据持久性，索引支持更快的查询
    
    2.丰富的语言查询：
      Mongodb 支持丰富的查询语言来支持读写操作（CRUD）以及数据汇总
    
    3.高可用性：
      Mongodb 的复制工具，成为副本集，提供自动故障转移和数据冗余，
    
    4.水平可扩展性：
      Mongodb 提供了可扩展性，作为其核心功能的一部分，分片是将数据分在一组计算机上。
    
    5.支持多种存储引擎：
      WiredTiger存储引擎和、 MMAPv1存储引擎和 InMemory 存储引擎
       3.0以上版本            3.0以下版本
      新的引擎压缩比特别大，原来100个G，可能升级之后所有数据都在，只占用10个G
       
    6.强大的索引支持：
      地理位置索引可用于构建 各种 O2O 应用、文本索引解决搜索的需求、TTL索引解决历史数据自动过期的需求
3.MongoDB应用场景
    1.游戏场景： 使用 MongoDB 存储游戏用户信息，用户的装备、积分等直接以内嵌文档的形式存储，方便查询、更新 
    2.物流场景： 使用 MongoDB 存储订单信息，订单状态在运送过程中会不断更新，以 MongoDB 内嵌数组的形式来存储，一次查询就能将订单所有的变更读取出来。 
    3.社交场景： 使用 MongoDB 存储存储用户信息，以及用户发表的朋友圈信息，通过地理位置索引实现附近的人、地点等功能   将送快递骑手、快递商家的信息（包含位置信息）存储在 MongoDB，然后通过 MongoDB 的地理位置查询，这样很方便的实现了查找附近的商家、骑手等功能，使得快递骑手能就近接单   地图软件、打车软件、外卖软件，MongoDB强大的地理位置索引功能使其最佳选择 
    4.物联网场景： 使用 MongoDB 存储所有接入的智能设备信息，以及设备汇报的日志信息，并对这些信息进行多维度的分析 
    5.视频直播： 使用 MongoDB 存储用户信息、礼物信息等 
    6.电商场景： 上衣有胸围，裤子有腰围，如果用数据库需要分成两个库，如果使用MongoDB都可以存在一起
### MongoDB安装部署
0.安装依赖
[root@redis01 ~]# yum install -y libcurl openssl
1.上传或下载包
#下载地址：https://www.mongodb.com/download-center/community
[root@redis01 ~]# rz mongodb-linux-x86_64-3.6.13.tgz 
2.解压包
[root@redis01 ~]# tar xf mongodb-linux-x86_64-3.6.13.tgz -C /usr/local/
[root@redis01 ~]# ln -s /usr/local/mongodb-linux-x86_64-3.6.13 /usr/local/mongodb
3.配置
#创建目录
[root@redis01 ~]# mkdir /server/mongo_27017/{conf,logs,pid,data} -p

#配置
[root@redis01 ~]# vim /server/mongo_27017/conf/mongodb.conf
systemLog:
  destination: file   
  logAppend: true  
  path: /server/mongo_27017/logs/mongodb.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongo_27017/data
  directoryPerDB: true
  wiredTiger:
     engineConfig:
        cacheSizeGB: 1
        directoryForIndexes: true
     collectionConfig:
        blockCompressor: zlib
     indexConfig:
        prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongo_27017/pid/mongod.pid

net:
  port: 27017
  bindIp: 127.0.0.1,10.0.0.91
  
  
#配置详解
#日志相关
systemLog:
  #以文件格式存储
  destination: file
  #每次重启，不生成新文件，每次都追加到文件
  logAppend: true
  #指定文件路径
  path: /server/mongo_27017/logs/mongodb.log

#数据部分
storage:
  #数据回滚。类似于mysql的undolog
  journal:
    enabled: true
  #数据目录
  dbPath: /server/mongo_27017/data
  #默认 false，不适用 inmemory engine
  directoryPerDB: true
  #存储引擎
  wiredTiger:
     #存储引擎设置
     engineConfig:
        #想把数据存到缓存，缓存的大小
        cacheSizeGB: 1
        #设置一个库就是一个目录，关闭就全放到一个目录下，很乱
        directoryForIndexes: true
     #压缩相关
     collectionConfig:
        blockCompressor: zlib
     #索引压缩（与压缩一起使用）
     indexConfig:
        prefixCompression: true
#守护进程的模式
processManagement:
  fork: true
  #指定pid文件
  pidFilePath: /server/mongo_27017/pid/mongod.pid
#指定端口和监听地址
net:
  port: 27017
  bindIp: 127.0.0.1,10.0.0.91
4.启动
[root@redis01 ~]# /usr/local/mongodb/bin/mongod -f /server/mongo_27017/conf/mongodb.conf 
about to fork child process, waiting until server is ready for connections.
forked process: 11547
child process started successfully, parent exiting

#验证启动
[root@redis01 ~]# ps -ef | grep mongo
root      11547      1  6 08:48 ?        00:00:00 /usr/local/mongodb/bin/mongod -f /server/mongo_27017/conf/mongodb.conf
5.配置环境变量
[root@redis01 ~]# vim /etc/profile.d/mongo.sh
export PATH="/usr/local/mongodb/bin:$PATH"

[root@redis01 ~]# source /etc/profile

### mongo登录警告处理
```
1.警告一
#访问设置没有被允许
WARNING: Access control is not enabled for the database.

#解决方式：开启安全认证
[root@redis01 ~]# vim /server/mongo_27017/conf/mongodb.conf 
security:
  authorization: enabled
2.警告二
#以root用户运行了
WARNING: You are running this process as the root user, which is not recommended.

#解决方式：使用普通用户启动
1.先关闭mongodb
[root@redis01 ~]# mongod -f /server/mongo_27017/conf/mongodb.conf --shutdown
2.创建mongo用户
[root@redis01 ~]# useradd mongo
[root@redis01 ~]# echo '123456'|passwd --stdin mongo
Changing password for user mongo.
passwd: all authentication tokens updated successfully.
3.授权目录
[root@redis01 ~]# chown -R mongo.mongo /usr/local/mongodb/
[root@redis01 ~]# chown -R mongo.mongo /server/mongo_27017/
4.重新启动
[root@redis01 ~]# su - mongo
Last login: Wed May 27 09:07:51 CST 2020 on pts/1
[mongo@redis01 ~]$ mongod -f /server/mongo_27017/conf/mongodb.conf
3.警告三、四
#你使用的是透明大页，可能导致mongo延迟和内存使用问题。
WARNING: /sys/kernel/mm/transparent_hugepage/enabled is 'always'.
       We suggest setting it to 'never'
WARNING: /sys/kernel/mm/transparent_hugepage/defrag is 'always'.
       We suggest setting it to 'never'

#解决方法：执行 echo never > /sys/kernel/mm/transparent_hugepage/enabled 修复该问题
#         执行 echo never > /sys/kernel/mm/transparent_hugepage/defrag 修复该问题

#配置之后重启
[root@redis01 ~]# su - mongo
[mongo@redis01 ~]$ mongod -f /server/mongo_27017/conf/mongodb.conf --shutdown
[mongo@redis01 ~]$ mongod -f /server/mongo_27017/conf/mongodb.conf

#这样设置是临时的，我们要把他加到 rc.local，再授权
4.警告五
#rlimits太低，MongoDB的软件进程被限制了，MongoDB希望自己是最少rlimits 32767.5
WARNING: soft rlimits too low. rlimits set to 7837 processes, 65535 files. Number of processes should be at least 32767.5 : 0.5 times number of files.

#解决方法：
[root@redis01 ~]# vim /etc/profile
ulimit -f unlimited
ulimit -t unlimited
ulimit -v unlimited
ulimit -n 65535
ulimit -m unlimited
ulimit -u 65535

[root@redis01 ~]# vim /etc/security/limits.d/20-nproc.conf
*          soft    nproc     65535
root       soft    nproc     unlimited

[root@redis01 ~]# source /etc/profile
```

### 基本操作
```
1.操作说明
    CRUD操作是create(创建)， read(读取)， update(更新)和delete(删除) 文档。
    MongoDB不支持SQL但是支持自己的丰富的查询语言。
    
    在MongoDB中，存储在集合中的每个文档都需要一个唯一的 _id字段，作为主键。如果插入的文档省略了该_id字段，则MongoDB驱动程序将自动为该字段生成一个ObjectId_id。也用于通过更新操作插入的文档upsert: true.如果文档包含一个_id字段，该_id值在集合中必须是唯一的，以避免重复键错误。
    
    在MongoDB中，插入操作针对单个集合。 MongoDB中的所有写操作都是在单个文档的级别上进行的
2.基本操作
    show databases/show dbs             #查看库列表
    show tables/show collections        #查看所有的集合
    use admin                           #切换库（没有的库也可以进行切换，没有数据是看不到）
    db                                #查看当前库
    show users                          #打印当前数据库的用户列表
    
    test:登录时默认存在的库
    admin库:系统预留库,MongoDB系统管理库
    local库:本地预留库,存储关键日志
    config库:MongoDB配置信息库
3.插入数据
1）单条数据插入
db.test.insert({"name":"lhd","age":18,"sex":"男"})

db.test.insert({"name":"lhd","age":18,"sex":"男","address":"上海浦东新区"})

db.test.insertOne({"name":"lhd","age":18,"sex":"男","address":"上海浦东新区"})
2）多条数据插入
db.inventory.insertMany( [
    { "name": "lhd", "age": 18, "figure": { "h": 182, "w": 200 }, "size": "big" },
    { "name": "qiudao", "age": 88, "figure": { "h": 120, "w": 160 }, "size": "very bittle" },
    { "name": "zengdao", "age": 18, "figure": { "h": 180, "w": 160 }, "size": "nomel" },
 ]);
4.查询数据
1）查询所有数据
> db.test.find()
{ "_id" : ObjectId("5ecdcdac13a4155a65ecb332"), "name" : "lhd", "age" : 18, "sex" : "男" }
{ "_id" : ObjectId("5ecdcdc413a4155a65ecb333"), "name" : "lhd", "age" : 18, "sex" : "男", "address" : "上海浦东新区" }
{ "_id" : ObjectId("5ecdcdd213a4155a65ecb334"), "name" : "lhd", "age" : 18, "sex" : "男" }
{ "_id" : ObjectId("5ecdcdd813a4155a65ecb335"), "name" : "lhd", "age" : 18, "sex" : "男", "address" : "上海浦东新区" }
2）查询单条数据
> db.test.findOne()
{
    "_id" : ObjectId("5ecdcdac13a4155a65ecb332"),
    "name" : "lhd",
    "age" : 18,
    "sex" : "男"
}
3）按条件查询
#如果查询条件为数字，不需要加引号
> db.test.findOne({"name" : "lhd"})
{
    "_id" : ObjectId("5ecdcdac13a4155a65ecb332"),
    "name" : "lhd",
    "age" : 18,
    "sex" : "男"
}
4）查询多条件
#并且的多个条件
> db.inventory.find({"figure.h":120,"size":"very bittle"})

> db.inventory.find(
    {
        "figure.h":120,
        "size":"very bittle"
    }
)

#表示多条件或者
> db.inventory.find({$or:[{"figure.h":120},{"size":"big"}]})

> db.inventory.find(
    {
        $or [
            {"figure.h":120},
            {"size":"big"}
        ]
    }
)
5）条件加范围的查询
> db.inventory.find({$or:[{"figure.h":{$lt:130}},{"size":"big"}]})

> db.inventory.find(
    {
        $or [
            {"figure.h":{$lt:130}},
            {"size":"big"}
        ]
    }
)
5.修改数据
1）修改单个数据
> db.inventory.updateOne({"name":"qiudao"},{$set:{"figure.h":130}})

> db.inventory.updateOne(
    #条件
    {"name":"qiudao"},
    {
        $set:
            #修改的值
            {"figure.h":130}
    }
)
2）修改多条数据
> db.table.updateMany({name:"niulei"},{$set:{age:"18"}})
{ "acknowledged" : true, "matchedCount" : 4, "modifiedCount" : 3 }
6.索引
1）查看执行计划
> db.inventory.find().explain()
{   #查询计划
    "queryPlanner" : {
        #计划版本
        "plannerVersion" : 1,
        #被查询的库和集合
        "namespace" : "test2.inventory",
        #查询索引设置
        "indexFilterSet" : false,
        #查询条件
        "parsedQuery" : {
            
        },
        #成功的执行计划
        "winningPlan" : {
            #全表扫描
            "stage" : "COLLSCAN",
            #查询方向
            "direction" : "forward"
        },
        #拒绝的计划
        "rejectedPlans" : [ ]
    },
    #服务器信息
    "serverInfo" : {
        "host" : "redis01",
        "port" : 27017,
        "version" : "3.6.13",
        "gitVersion" : "db3c76679b7a3d9b443a0e1b3e45ed02b88c539f"
    },
    "ok" : 1
}

COLLSCAN 全表扫描
IXSCAN 索引扫描
2）创建索引
> db.inventory.createIndex({"age":1},{background:true})
{
    "createdCollectionAutomatically" : true,
    "numIndexesBefore" : 1,
    "numIndexesAfter" : 2,
    "ok" : 1
}

#添加索引
createIndex({索引的名称:1}) ：1表示正序，-1表示倒序

#创建方式
1.前台方式 
缺省情况下，当为一个集合创建索引时，这个操作将阻塞其他的所有操作。即该集合上的无法正常读写，直到索引创建完毕
任意基于所有数据库申请读或写锁都将等待直到前台完成索引创建操作
 
2.后台方式
将索引创建置于到后台，适用于那些需要长时间创建索引的情形
这样子在创建索引期间，MongoDB依旧可以正常的为提供读写操作服务
等同于关系型数据库在创建索引的时候指定online，而MongoDB则是指定background
其目的都是相同的，即在索引创建期间，尽可能的以一种占用较少的资源占用方式来实现，同时又可以提供读写服务
后台创建方式的代价：索引创建时间变长

#规范
1.如果要查询的内容都是最近的，那建立索引就用倒序，如果要通盘查询那就用正序。
2.比如说一个数据集合查询占的比较多就用索引，如果查询少而是插入数据比较多就不用建立索引。因为：当没有索引的时候，插入数据的时候MongoDB会在内存中分配出一块空间，用来存放数据。当有索引的时候在插入数据之后还会给自动添加一个索引，浪费了时间。
3.不是所有数据都要建立索引，要在恰当并且需要的地方建立才是最好的。
4.大数量的排序工作时可以考虑创建索引。
3）查看索引
> db.inventory.getIndexes()
[
    {
        "v" : 2,
        "key" : {
            "_id" : 1
        },
        "name" : "_id_",
        "ns" : "test2.test2"
    },
    {
        "v" : 2,
        "key" : {
            "age" : 1
        },
        "name" : "name_1",
        "ns" : "test2.test2",
        "background" : true
    }
]
4）再次查看执行计划
> db.inventory.find({"age":{$lt:40}}).explain()
        "winningPlan" : {
            "stage" : "FETCH",
            "inputStage" : {
                #走索引了
                "stage" : "IXSCAN",
6.删除
1）删除单条数据
> db.inventory.deleteOne({"name":"lhd"})
{ "acknowledged" : true, "deletedCount" : 1 }
2）删除多个数据
> db.inventory.deleteMany({"name":"lhd"})
3）删除索引
> db.test.dropIndex({ age: 1 })
{
    "ok" : 0,
    "errmsg" : "ns not found",
    "code" : 26,
    "codeName" : "NamespaceNotFound"
}
4）删除集合
#先确认自己在哪个库
> db
test2

#确认集合
> show tables;
inventory
test

#删除集合
> db.inventory.drop()
true
5）删除库
#先确认自己在哪个库
> db
test2

#删除库
> db.dropDatabase()
```
### mongo工具
```
    mongo               #登录命令
    mongodump           #备份导出，全备（数据时压缩过的）
    mongorestore        #恢复数据
    mongostat           #查看运行状态的
    mongod              #启动命令
    mongoexport         #备份，导出json格式
    mongoimport         #恢复数据
    mongos              #集群分片命令
    mongotop            #查看运行状态
1.mongostat命令
#不加任何参数时，每秒访问一次
[mongo@redis01 ~]$ mongostat
insert query update delete getmore command dirty used flushes vsize   res qrw arw net_in net_out conn                time
    *0    *0     *0     *0       0     2|0  0.0% 0.0%       0  972M 56.0M 0|0 1|0   158b   60.9k    1 May 27 11:23:08.248

insert      #每秒插入数据的数量
query       #每秒查询操作的数量
update      #每秒更新数据的数量
delete      #没面删除操作的数量
getmore     #每秒查询游标时的操作数
command     #每秒执行的命令数
dirty       #脏数据占缓存的多少
used        #使用中的缓存
flushes
            #在 wiredtiger引擎，表示轮询间隔
            #在MMapv1引擎，表示每秒写入磁盘次数
vsize       #虚拟内存使用量
res         #物理内存使用量
qrw         #客户端等待读数据的队列长度
arw         #客户端等待写入数据的队列长度
net_in      #网络进流量
net_out     #网络出流量
conn        #连接总数
time        #时间

#一般该命令搭配  mongotop 命令使用，可以显示每个集合的响应速度
```
### 用户授权认证
1.授权命令
用户管理界面
要添加用户， MongoDB提供了该db.createUser()方法。添加用户时，您可以为用户分配色以授予权限。
注意：
在数据库中创建的第一个用户应该是具有管理其他用户的权限的用户管理员。
您还可以更新现有用户，例如更改密码并授予或撤销角色。

db.auth() 将用户验证到数据库。
db.changeUserPassword() 更改现有用户的密码。
db.createUser() 创建一个新用户。
db.dropUser() 删除单个用户。
db.dropAllUsers() 删除与数据库关联的所有用户。
db.getUser() 返回有关指定用户的信息。
db.getUsers() 返回有关与数据库关联的所有用户的信息。
db.grantRolesToUser() 授予用户角色及其特权。
db.removeUser() 已过时。从数据库中删除用户。
db.revokeRolesFromUser() 从用户中删除角色。
db.updateUser() 更新用户数据。
2.创建用户和角色
[mongo@db01 ~]$ mongo
> use admin
> db.createUser({user: "admin",pwd: "123456",roles:[ { role: "root", db:"admin"}]})
Successfully added user: {
        "user" : "admin",
        "roles" : [
                {
                        "role" : "root",
                        "db" : "admin"
                }
        ]
}
3.查看用户
> db.getUsers()
[
    {
        "_id" : "test.admin",
        "userId" : UUID("b840b96c-3442-492e-a45f-6ca7dff907fd"),
        "user" : "admin",
        "db" : "test",
        "roles" : [
            {
                "role" : "root",
                "db" : "admin"
            }
        ]
    }
]
4.配置开启认证
[root@redis01 ~]# vim /server/mongo_27017/conf/mongodb.conf 
security:
  authorization: enabled
  
#重启
[mongo@redis01 ~]$ mongod -f /server/mongo_27017/conf/mongodb.conf --shutdown
[mongo@redis01 ~]$ mongod -f /server/mongo_27017/conf/mongodb.conf
5.配置认证以后查操作不了
> show databases;
2020-05-27T11:41:06.186+0800 E QUERY    [thread1] Error: listDatabases failed:{
    "ok" : 0,
    "errmsg" : "there are no users authenticated",
    "code" : 13,
    "codeName" : "Unauthorized"
} :
_getErrorWithCode@src/mongo/shell/utils.js:25:13
Mongo.prototype.getDBs@src/mongo/shell/mongo.js:67:1
shellHelper.show@src/mongo/shell/utils.js:860:19
shellHelper@src/mongo/shell/utils.js:750:15
@(shellhelp2):1:1
6.使用账号密码连接
[mongo@redis01 ~]$ mongo -uadmin -p --authenticationDatabase admin
MongoDB shell version v3.6.13
Enter password: 

> show databases;
admin   0.000GB
config  0.000GB
local   0.000GB
test    0.000GB
test2   0.000GB
7.创建普通用户
> use test
> db.createUser(
  {
    user: "test",
    pwd: "123456",
    roles: [ { role: "readWrite", db: "write" },
             { role: "read", db: "read" } ]
  }
)
8.创建测试数据
use write
db.write.insert({"name":"linhaoda","age":17,"ad":"上海浦东新区"})
db.write.insert({"name":"linhaoda","age":17,"ad":"上海浦东新区"})
db.write.insert({"name":"haoda","age":18,"ad":"上海浦东新区"})
db.write.insert({"name":"linda","age":18,"ad":"上海浦东新区"})
db.write.insert({"name":"linhao","age":18,"ad":"上海浦东新区","sex":"boy"})

use read
db.read.insert({"name":"linhaoda","age":17,"ad":"上海浦东新区"})
db.read.insert({"name":"linhaoda","age":17,"ad":"上海浦东新区"})
db.read.insert({"name":"haoda","age":18,"ad":"上海浦东新区"})
db.read.insert({"name":"linda","age":18,"ad":"上海浦东新区"})
db.read.insert({"name":"linhao","age":18,"ad":"上海浦东新区","sex":"boy"})
9.验证
mongo -utest -p --authenticationDatabase test
use write
db.write.find()
db.write.insert({"name":"linhaoda","age":17,"ad":"上海浦东新区"})

use read
db.read.find()
db.read.insert({"name":"linhaoda","age":17,"ad":"上海浦东新区"})
10.修改用户密码
#修改用户信息
db.updateUser("test",{pwd:"123"})

#修改密码
db.changeUserPassword("admin","123")
### 副本集的搭建
介绍副本集
    官网的图片 https://docs.mongodb.com/manual/replication/
1.创建多实例目录
[root@redis03 ~]# mkdir /server/mongodb/2801{7,8,9}/{conf,logs,pid,data} -p
2.编辑多实例配置文件
[root@redis03 ~]# vim /server/mongodb/28017/conf/mongo.conf
systemLog:
  destination: file
  logAppend: true
  path: /server/mongodb/28017/logs/mongodb.log
  #path: /server/mongodb/28018/logs/mongodb.log
  #path: /server/mongodb/28019/logs/mongodb.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/28017/data
  #dbPath: /server/mongodb/28018/data
  #dbPath: /server/mongodb/28019/data
  directoryPerDB: true
  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/28017/pid/mongod.pid
  #pidFilePath: /server/mongodb/28018/pid/mongod.pid
  #pidFilePath: /server/mongodb/28019/pid/mongod.pid

net:
  port: 28017
  #port: 28018
  #port: 28019
  bindIp: 127.0.0.1,10.0.0.93
  
replication:
  #类似于binlog，指定大小
  oplogSizeMB: 1024
  #副本记得名称，集群名称
  replSetName: dba
3.启动多实例
[root@redis03 ~]# chown -R mongo.mongo /server/mongodb/
[root@redis03 ~]# su - mongo

[mongo@redis03 ~]$ mongod -f /server/mongodb/28017/conf/mongo.conf
[mongo@redis03 ~]$ mongod -f /server/mongodb/28018/conf/mongo.conf
[mongo@redis03 ~]$ mongod -f /server/mongodb/28019/conf/mongo.conf

#验证
[mongo@redis03 ~]$ netstat -lntp       
tcp        0      0 10.0.0.93:28017         0.0.0.0:*               LISTEN      32893/mongod        
tcp        0      0 127.0.0.1:28017         0.0.0.0:*               LISTEN      32893/mongod        
tcp        0      0 10.0.0.93:28018         0.0.0.0:*               LISTEN      32938/mongod        
tcp        0      0 127.0.0.1:28018         0.0.0.0:*               LISTEN      32938/mongod        
tcp        0      0 10.0.0.93:28019         0.0.0.0:*               LISTEN      32981/mongod        
tcp        0      0 127.0.0.1:28019         0.0.0.0:*               LISTEN      32981/mongod
4.登录多实例
[mongo@redis03 ~]$ mongo 10.0.0.93:28017
[mongo@redis03 ~]$ mongo 10.0.0.93:28018
[mongo@redis03 ~]$ mongo 10.0.0.93:28019
5.初始化副本集
#配置副本集
config = {
  _id : "dba", 
  members : [
    {_id:0, host:"10.0.0.93:28017"},
    {_id:1, host:"10.0.0.93:28018"},
    {_id:2, host:"10.0.0.93:28019"},
  ]
}

#读取副本集
rs.initiate(config) 
6.查看副本集状态
dba:PRIMARY> rs.status()
            #健康状态 1表示正常 0表示故障
            "health" : 1,
            #表示状态 1是主库 2是从库 3表示恢复数据中 7表示投票者 8表示down机
            "state" : 1,
            #标注是主库还是从库
            "stateStr" : "PRIMARY",
            #集群启动时间
            "uptime" : 579,
            #另一种格式的时间
            "optime" : {
                "ts" : Timestamp(1590593779, 1),
                "t" : NumberLong(1)
            },
            #上一次心跳传过来数据的时间
            "optimeDate" : ISODate("2020-05-27T15:36:19Z"),
            #检测上一次心跳时间
            "lastHeartbeat" : ISODate("2020-05-27T15:36:25.815Z"),
            
#查看集群与主节点
dba:PRIMARY> rs.isMaster()

#oplog信息
dba:PRIMARY> rs.printReplicationInfo()
configured oplog size:   1024MB
log length start to end: 1543secs (0.43hrs)
oplog first event time:  Wed May 27 2020 23:26:46 GMT+0800 (CST)
oplog last event time:   Wed May 27 2020 23:52:29 GMT+0800 (CST)
now:                     Wed May 27 2020 23:52:38 GMT+0800 (CST)

#查看延时从库信息
dba:PRIMARY> rs.printSlaveReplicationInfo()
source: 10.0.0.93:28018
    syncedTo: Wed May 27 2020 23:54:19 GMT+0800 (CST)
    0 secs (0 hrs) behind the primary 
source: 10.0.0.93:28019
    syncedTo: Wed May 27 2020 23:54:19 GMT+0800 (CST)
    0 secs (0 hrs) behind the primary 
    
#打印副本集配置文件
dba:PRIMARY> rs.config()
7.主库创建数据，从库查看数据
#主库插入数据
db.table.insertMany([{"name":"gcc","age":10},{"name":"zzy","age":9},{"name":"hxh","age":11}])
#主库查看数据
dba:PRIMARY> show tables
table
dba:PRIMARY> db.table.find()

#从库查看数据
dba:SECONDARY> show databases
2020-05-27T23:43:40.020+0800 E QUERY    [thread1] Error: listDatabases failed:{
    "operationTime" : Timestamp(1590594219, 1),
    "ok" : 0,
    "errmsg" : "not master and slaveOk=false",
    "code" : 13435,
    "codeName" : "NotMasterNoSlaveOk",
    "$clusterTime" : {
        "clusterTime" : Timestamp(1590594219, 1),
        "signature" : {
            "hash" : BinData(0,"AAAAAAAAAAAAAAAAAAAAAAAAAAA="),
            "keyId" : NumberLong(0)
        }
    }
} 
#连查看库都会被拒绝，因为从库不提供读写

#执行命令（从库都要执行）
dba:SECONDARY> rs.slaveOk()
dba:SECONDARY> show databases
admin    0.000GB
cluster  0.000GB
config   0.000GB
local    0.000GB

#每次重新连接都要执行以上命令才能读取
#可以配置永久生效
[root@redis03 ~]# vim ~/.mongorc.js
rs.slaveOk()
### 副本集实现高可用
1.故障切换测试
#主库使用 localhost 连接，执行关闭数据库的操作（使用ip连接是不能执行的）
[root@db01 ~]# mongod -f /server/mongodb/28018/conf/mongo.conf --shutdown

#查看其他从库中会有一台从库，变成主库

#故障转移实现了，但是我的程序连接mongodb的配置还需要修改怎么办？？
2.程序怎么实现连接切换的
1.如果使用的是单节点，那么程序里面直接配置写死mongodb的ip和端口即可

2.如果是副本集集群的形式，在程序里面写的就是一个列表，列表里面写
    mongo_reip=[10.0.0.91:28017,10.0.0.92:28018,10.0.0.93:29019]
    程序会去使用命令询问谁是主节点，得到结果后在写入数据
3.恢复主库
#重新启动主库，他会自动判断谁是主库，自动成为新的从库

#注意：三台节点，只能坏一台，坏两台就有问题了
4.指定节点提升优先级
#原来的主库配置高，性能好，想恢复之后还让他是主库怎么办

#查看优先级
dba:PRIMARY> rs.conf()
            #权重值
            "priority" : 1,
            
#临时修改配置文件
dba:PRIMARY> config=rs.conf()
#修改配置文件中 id 为0 的priority值为10
dba:PRIMARY> config.members[0].priority=10
#配置文件生效
dba:PRIMARY> rs.reconfig(config)

#新版本调整完直接切换主库，旧版本需要主动降级
dba:PRIMARY> rs.stepDown()

#恢复权重
dba:PRIMARY> config=rs.conf()
dba:PRIMARY> config.members[0].priority=1
dba:PRIMARY> rs.reconfig(config)
扩容与删减节点
1.配置一台新的节点
#创建目录
[root@redis03 ~]# mkdir /server/mongodb/28016/{conf,logs,pid,data} -p

#配置新节点
[root@redis03 ~]# cp /server/mongodb/28017/conf/mongo.conf /server/mongodb/28016/conf/
[root@redis03 ~]# sed -i 's#28017#28016#g' /server/mongodb/28016/conf/mongo.conf 

#启动新节点
[root@redis03 ~]# chown -R mongo.mongo /server/mongodb/
[root@redis03 ~]# su - mongo
[mongo@redis03 ~]$ mongod -f /server/mongodb/28016/conf/mongo.conf
2.将新节点加入集群
#主库操作
dba:PRIMARY> rs.add("10.0.0.93:28016")
{
    "ok" : 1,
    "operationTime" : Timestamp(1590597530, 1),
    "$clusterTime" : {
        "clusterTime" : Timestamp(1590597530, 1),
        "signature" : {
            "hash" : BinData(0,"AAAAAAAAAAAAAAAAAAAAAAAAAAA="),
            "keyId" : NumberLong(0)
        }
    }
}

#查看集群状态
dba:PRIMARY> rs.status()

#注意：四个节点也不能坏两台机器
3.删除节点
#主库操作
dba:PRIMARY> rs.remove("10.0.0.93:28016")
{
    "ok" : 1,
    "operationTime" : Timestamp(1590597842, 1),
    "$clusterTime" : {
        "clusterTime" : Timestamp(1590597842, 1),
        "signature" : {
            "hash" : BinData(0,"AAAAAAAAAAAAAAAAAAAAAAAAAAA="),
            "keyId" : NumberLong(0)
        }
    }
}

#查看集群状态
dba:PRIMARY> rs.status()
4.添加仲裁节点
#创建目录
[root@redis03 ~]# mkdir /server/mongodb/28015/{conf,logs,pid,data} -p

#配置新节点
[root@redis03 ~]# cp /server/mongodb/28017/conf/mongo.conf /server/mongodb/28015/conf/
[root@redis03 ~]# sed -i 's#28017#28015#g' /server/mongodb/28015/conf/mongo.conf 

#启动新节点
[root@redis03 ~]# chown -R mongo.mongo /server/mongodb/
[root@redis03 ~]# su - mongo
[mongo@redis03 ~]$ mongod -f /server/mongodb/28015/conf/mongo.conf 

#主库操作加入仲裁节点
dba:PRIMARY> rs.addArb(("10.0.0.93:28015")

#查看该库是否有数据

#注意，五个节点时，可以坏两个节点

### 备份与恢复数据
1.备份恢复工具
    1.mongoexport/mongoimport       #数据分析时使用
    2.mongodump/mongorestore        #单纯备份时使用
2.导出工具mongoexport
#备份成json格式
[mongo@redis03 ~]$ mongoexport --port 28018 -d test -c testtable -o ~/table.json

[mongo@redis03 ~]$ mongoexport -uadmin -p123456 --port 27017 --authenticationDatabase admin -d database -c table -o ~/table.json

#备份成csv格式
[mongo@redis03 ~]$ mongoexport --port 27017 -d database -c table --type=csv -f name,age -o ~/table.csv

[mongo@redis03 ~]$ mongoexport -uadmin -p123456 --port 27017 --authenticationDatabase admin -d database -c table --type=csv -f name,age -o ~/table.csv


-h:指明数据库宿主机的IP
-u:指明数据库的用户名
-p:指明数据库的密码
-d:指明数据库的名字
-c:指明集合的名字
-f:指明要导出那些列
-o:指明到要导出的文件名
-q:指明导出数据的过滤条件
--type:指明导出数据的类型
3.恢复工具mongoimport
#删除集合
> use database
switched to db database
> show tables
table
> db.table.drop()
true

#恢复数据
[mongo@redis03 ~]$ mongoimport --port 27017 -d database -c table ~/table.json

[mongo@redis03 ~]$ mongoimport --port 27017 -d database -c table --type=csv --headerline --file ~/table.csv

-h:指明数据库宿主机的IP
-u:指明数据库的用户名
-p:指明数据库的密码
-d:指明数据库的名字
-c:指明集合的名字
-f:指明要导入那些列
4.生产案例：数据库数据迁移至mongodb
1）搭建数据库
2）导入数据
3）配置数据库
#开启安全路径
[root@redis04 ~]# vim /etc/my.cnf

[mysqld]
basedir=/usr/local/mysql
datadir=/usr/local/mysql/data
secure-file-priv=/tmp

#重启数据库
[root@redis04 ~]# systemctl restart mysql
4）将数据库导出成csv
#使用第三方工具导出csv表格
mysql> select * from world.city into outfile '/tmp/city1.csv' fields terminated by ',';
5）查看文件
[root@redis04 ~]# cat /tmp/city1.csv
6）手动处理文件
#将数据库字段加到文件的第一行
[root@redis04 ~]# vim /tmp/city.csv
ID,Name,CountryCode,District,Population
1,Kabul,AFG,Kabol,1780000
2,Qandahar,AFG,Qandahar,237500
7）将数据导入mongodb
[root@redis04 ~]# scp /tmp/city1.csv 172.16.1.93:/tmp/

[mongo@redis03 ~]$ mongoimport -uadmin -p123456 --port 27017 --authenticationDatabase admin -d world -c city --type=csv --headerline --file /tmp/city1.csv
8）查看数据
[mongo@redis03 ~]$ mongo -uadmin -p123456 --authenticationDatabase admin
> show dbs
admin   0.000GB
config  0.000GB
local   0.000GB
read    0.000GB
world   0.000GB
write   0.000GB
> use world
switched to db world
> show tables
city
> db.city.find()
......
> it
5.生产案例：数据误删除恢复
1）过程
每天凌晨1点进行全备
10点进行误操作，删除了数据
恢复数据
2）模拟全备数据
#连接副本集的主库（只有在副本集模式才能使用mongodump）
[mongo@redis03 ~]$ mongo localhost:28018
dba:PRIMARY> use backup
dba:PRIMARY> db.backuptable.insertMany([{id:1},{id:2},{id:3}])
{
    "acknowledged" : true,
    "insertedIds" : [
        ObjectId("5ecfe698e99e372e2e4fe1fd"),
        ObjectId("5ecfe698e99e372e2e4fe1fe"),
        ObjectId("5ecfe698e99e372e2e4fe1ff")
    ]
}
dba:PRIMARY> db.backuptable.find()
{ "_id" : ObjectId("5ecfe698e99e372e2e4fe1fd"), "id" : 1 }
{ "_id" : ObjectId("5ecfe698e99e372e2e4fe1fe"), "id" : 2 }
{ "_id" : ObjectId("5ecfe698e99e372e2e4fe1ff"), "id" : 3 }
3）执行全备
[mongo@redis03 ~]$ mongodump --port 28018 --oplog -o /data

[mongo@redis03 ~]$ ll /data/oplog.bson 
-rw-rw-r-- 1 mongo mongo 110 May 29 02:01 /data/oplog.bson
4）模拟增量数据
[mongo@redis03 ~]$ mongo 10.0.0.93:28018
dba:PRIMARY> use backup
switched to db backup
dba:PRIMARY> db.backuptable.insertMany([{id:4},{id:5},{id:6}])
{
    "acknowledged" : true,
    "insertedIds" : [
        ObjectId("5ecfe86f5c1085fcf692a3cb"),
        ObjectId("5ecfe86f5c1085fcf692a3cc"),
        ObjectId("5ecfe86f5c1085fcf692a3cd")
    ]
}
5）删除数据
dba:PRIMARY> use backup
switched to db backup
dba:PRIMARY> db.backuptable.drop()
true
dba:PRIMARY> show tables
7）oplog
    oplog是local库下的一个固定集合，从库就是通过查看主库的oplog这个集合来进行复制的。每个节点都有oplog，记录这从主节点复制过来的信息，这样每个成员都可以保证切换主库时的数据同步
6）查找删除动作的时间点
#连接mongodb
[mongo@redis03 ~]$ mongo 10.0.0.93:28018
#切换到local库
dba:PRIMARY> use local
#查看oplog信息
dba:PRIMARY> db.oplog.rs.find()
dba:PRIMARY> db.oplog.rs.find().pretty()
{   
    #同步的时间点，选举时会选择最新的时间戳提升为主库
    "ts" : Timestamp(1590640219, 1),
    "t" : NumberLong(1),
    "h" : NumberLong("-8962736529514397515"),
    "v" : 2,
    #操作类型 i代表insert u代表update d代表delete n代表没有操作只是保持连接发送消息
    "op" : "n",
    #当前数据库的库、表
    "ns" : "",
    "wall" : ISODate("2020-05-28T04:30:19.080Z"),
    #操作的内容
    "o" : {
        "msg" : "periodic noop"
    }
}

#oplog信息
dba:PRIMARY> rs.printReplicationInfo()
configured oplog size:   1024MB                                 #oplog文件大小
log length start to end: 1543secs (0.43hrs)                     #oplog日志的启用时间段
oplog first event time:  Wed May 27 2020 23:26:46 GMT+0800 (CST)    #第一个事务日志的产生时间
oplog last event time:   Wed May 27 2020 23:52:29 GMT+0800 (CST)    #最后一个事务日志的产生时间
now:                     Wed May 27 2020 23:52:38 GMT+0800 (CST)    #现在的时间

#查找到删除的时间点
dba:PRIMARY> db.oplog.rs.find({ns:"backup.$cmd"}).pretty()
{
    "ts" : Timestamp(1590683811, 1),
    "t" : NumberLong(2),
    "h" : NumberLong("3968458855036608631"),
    "v" : 2,
    "op" : "c",
    "ns" : "backup.$cmd",
    "ui" : UUID("bec471f5-cd2a-44fe-8056-4c5c2de5de03"),
    "wall" : ISODate("2020-05-28T16:36:51.227Z"),
    "o" : {
        "drop" : "backuptable"
    }
}

1590683811
7）备份最新的oplog
[mongo@redis03 ~]$ mongodump --port 28018 -d local -c oplog.rs -o /data/

[mongo@redis03 ~]$ ll /data/local/
total 140
-rw-rw-r-- 1 mongo mongo 138093 May 29 00:56 oplog.rs.bson
-rw-rw-r-- 1 mongo mongo    125 May 29 00:56 oplog.rs.metadata.json
8）把原来的全备备份
[mongo@redis03 data]$ mv oplog.bson oplog.bson.bak
[mongo@redis03 ~]$ mv /data/local/oplog.rs.bson /data/oplog.bson
9）恢复数据
#删掉新备份的库数据，否则会覆盖
[mongo@redis03 data]$ rm -rf /data/local
#恢复到指定时间点的数据
[mongo@redis03 data]$ mongorestore --port 28018 --oplogReplay --oplogLimit "1590690412:1" --drop /data/
10）查看数据
[mongo@redis03 ~]$ mongo localhost:28018
dba:PRIMARY> show databases
dba:PRIMARY> use backup
switched to db backup
dba:PRIMARY> show tables;
dba:PRIMARY> db.backuptable.find()
6.mongo升级
1.首先确保是副本集状态
2.先关闭1个副本节点
3.检测数据是否可以升级
4.升级副本节点的可执行文件
5.更新配置文件
6.启动升级后的副本节点
7.确保集群工作正常
8.滚动升级其他副本节点
9.最后主节点降级
10.确保集群 可用
11.关闭降级的老的主节点
12.升级老的主节点
13.重新加入集群

### mongodb的分片
1.分片的概念
    mongodb的副本集跟redis的高可用相同，只能读，分担不了主库的压力，只能在主库出现故障的时候接替主库的工作
    mongodb能够使用的内存，只是主库的内存和磁盘，当副本集中机器配置不一致时也会有问题
2.分片的介绍
    优点：
       1.提高机器资源的利用率
       2.减轻主库的压力
    缺点：
       1.机器需要的更多
       2.配置和管理更加的复杂和困难
       3.分片配置好之后想修改很困难
3.分片的原理
1）路由服务 mongos server
    类似于代理，跟数据库的atlas类似，可以将客户端的数据分配到后端的mongo服务器上
2）分片配置服务器信息 config server
    mongos server是不知道后端服务器mongo有几台，地址是什么，他只能连接到这个config server，而config server就是记录后端服务器地址和数据的一个服务
    作用：
       1.记录后端mongo节点的信息
       2.记录数据写入存到了哪个节点
       3.提供给mongos后端服务器的信息
3）片键
    config server只存储信息，而不会主动将数据写入节点，所以还有一个片键的概念，片键就是索引
    作用：
       1.将数据根据规则分配到不同的节点
       2.相当于建立索引，加快访问速度
    分类：
       1.区间片键（很有可能出现数据分配不均匀的情况）
          可以以时间区间分片，根据时间建立索引
          可以以地区区间分片，根据地区建立索引
       2.hash片键（足够平均，足够随机）
          根据id或者数据数量进行分配
4）分片
    存储数据的节点，这种方式就是分布式集群

### 分片的高可用
    做分片只是针对单节点，mongo服务相当于还是只有一个，所以我们还有对分片进行副本集的操作 
    跟ES一样，我们不能一台机器上部署多节点，自己做自己的副本，那当机器挂了时，还是有问题 
    所以我们要错开进行副本集的建立，而且一台机器上不能有相同的数据节点，否则选举又会出现问题
1.服务器规划
主机	ip	部署	端口
mongodb01	10.0.0.81	Shard1_Master Shard2_Slave Shard3_Arbiter Config Server Mongos Server	20010 28020 28030 40000 60000
mongodb02	10.0.0.82	Shard2_Master Shard3_Slave Shard1_Arbiter Config Server Mongos Server	20010 28020 28030 40000 60000
mongodb03	10.0.0.83	Shard3_Master Shard1_Slave Shard2_Arbiter Config Server Mongos Server	20010 28020 28030 40000 60000
2.目录规划
#服务目录
mkdir /server/mongodb/master/{conf,log,pid,data} -p
mkdir /server/mongodb/slave/{conf,log,pid,data} -p
mkdir /server/mongodb/arbiter/{conf,log,pid,data} -p
mkdir /server/mongodb/config/{conf,log,pid,data} -p
mkdir /server/mongodb/mongos/{conf,log,pid} -p
3.安装mongo
#安装依赖
yum install -y libcurl openssl
#上传或下载包
rz mongodb-linux-x86_64-3.6.13.tgz
#解压
tar xf mongodb-linux-x86_64-3.6.13.tgz -C /usr/local/
#做软连接
ln -s /usr/local/mongodb-linux-x86_64-3.6.13 /usr/local/mongodb
4.配置mongodb01
1）配置master
vim /server/mongodb/master/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/master/log/master.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/master/data/
  directoryPerDB: true

  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/master/pid/master.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 28010
  bindIp: 127.0.0.1,10.0.0.81

replication:
  oplogSizeMB: 1024 
  replSetName: shard1

sharding:
  clusterRole: shardsvr
2）配置slave
vim /server/mongodb/slave/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/slave/log/slave.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/slave/data/
  directoryPerDB: true

  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/slave/pid/slave.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 28020
  bindIp: 127.0.0.1,10.0.0.81

replication:
  oplogSizeMB: 1024
  replSetName: shard2

sharding:
  clusterRole: shardsvr
3）配置arbiter
vim /server/mongodb/arbiter/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/arbiter/log/arbiter.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/arbiter/data/
  directoryPerDB: true

  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/arbiter/pid/arbiter.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 28030
  bindIp: 127.0.0.1,10.0.0.81

replication:
  oplogSizeMB: 1024
  replSetName: shard3

sharding:
  clusterRole: shardsvr
5.配置mongodb02
1）配置master
vim /server/mongodb/master/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/master/log/master.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/master/data/
  directoryPerDB: true

  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/master/pid/master.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 28010
  bindIp: 127.0.0.1,10.0.0.82

replication:
  oplogSizeMB: 1024 
  replSetName: shard2

sharding:
  clusterRole: shardsvr
2）配置slave
vim /server/mongodb/slave/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/slave/log/slave.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/slave/data/
  directoryPerDB: true

  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/slave/pid/slave.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 28020
  bindIp: 127.0.0.1,10.0.0.82

replication:
  oplogSizeMB: 1024
  replSetName: shard3

sharding:
  clusterRole: shardsvr
3）配置arbiter
vim /server/mongodb/arbiter/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/arbiter/log/arbiter.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/arbiter/data/
  directoryPerDB: true

  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/arbiter/pid/arbiter.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 28030
  bindIp: 127.0.0.1,10.0.0.82

replication:
  oplogSizeMB: 1024
  replSetName: shard1

sharding:
  clusterRole: shardsvr
6.配置mongodb03
1）配置master
vim /server/mongodb/master/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/master/log/master.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/master/data/
  directoryPerDB: true

  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/master/pid/master.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 28010
  bindIp: 127.0.0.1,10.0.0.83

replication:
  oplogSizeMB: 1024 
  replSetName: shard3

sharding:
  clusterRole: shardsvr
2）配置slave
vim /server/mongodb/slave/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/slave/log/slave.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/slave/data/
  directoryPerDB: true

  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/slave/pid/slave.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 28020
  bindIp: 127.0.0.1,10.0.0.83

replication:
  oplogSizeMB: 1024
  replSetName: shard1

sharding:
  clusterRole: shardsvr
3）配置arbiter
vim /server/mongodb/arbiter/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/arbiter/log/arbiter.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/arbiter/data/
  directoryPerDB: true

  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/arbiter/pid/arbiter.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 28030
  bindIp: 127.0.0.1,10.0.0.83

replication:
  oplogSizeMB: 1024
  replSetName: shard2

sharding:
  clusterRole: shardsvr
7.配置环境变量
[root@redis01 ~]# vim /etc/profile.d/mongo.sh
export PATH="/usr/local/mongodb/bin:$PATH"

[root@redis01 ~]# source /etc/profile
8.优化警告
useradd mongo -s /sbin/nologin -M 
echo "never"  > /sys/kernel/mm/transparent_hugepage/enabled
echo "never"  > /sys/kernel/mm/transparent_hugepage/defrag
9.配置system管理
1）配置master管理
vim /usr/lib/systemd/system/mongod-master.service

[Unit]
Description=MongoDB Database Server
Documentation=https://docs.mongodb.org/manual
After=network.target

[Service]
User=mongo
Group=mongo
ExecStart=/usr/local/mongodb/bin/mongod -f /server/mongodb/master/conf/mongo.conf
ExecStartPre=/usr/bin/chown -R mongo:mongo /server/mongodb/master/
ExecStop=/usr/local/mongodb/bin/mongod -f /server/mongodb/master/conf/mongo.conf --shutdown
PermissionsStartOnly=true
PIDFile=/server/mongodb/master/pid/master.pid
Type=forking

[Install]
WantedBy=multi-user.target
2）配置管理salve
vim /usr/lib/systemd/system/mongod-slave.service

[Unit]
Description=MongoDB Database Server
Documentation=https://docs.mongodb.org/manual
After=network.target

[Service]
User=mongo
Group=mongo
ExecStart=/usr/local/mongodb/bin/mongod -f /server/mongodb/slave/conf/mongo.conf
ExecStartPre=/usr/bin/chown -R mongo:mongo /server/mongodb/slave/
ExecStop=/usr/local/mongodb/bin/mongod -f /server/mongodb/slave/conf/mongo.conf --shutdown
PermissionsStartOnly=true
PIDFile=/server/mongodb/slave/pid/slave.pid
Type=forking

[Install]
WantedBy=multi-user.target
3）配置管理arbiter
vim /usr/lib/systemd/system/mongod-arbiter.service

[Unit]
Description=MongoDB Database Server
Documentation=https://docs.mongodb.org/manual
After=network.target

[Service]
User=mongo
Group=mongo
ExecStart=/usr/local/mongodb/bin/mongod -f /server/mongodb/arbiter/conf/mongo.conf
ExecStartPre=/usr/bin/chown -R mongo:mongo /server/mongodb/arbiter/
ExecStop=/usr/local/mongodb/bin/mongod -f /server/mongodb/arbiter/conf/mongo.conf --shutdown
PermissionsStartOnly=true
PIDFile=/server/mongodb/arbiter/pid/arbiter.pid
Type=forking

[Install]
WantedBy=multi-user.target
4）刷新启动程序
    systemctl daemon-reload
10.启动mongodb所有节点
    systemctl start mongod-master.service
    systemctl start mongod-slave.service
    systemctl start mongod-arbiter.service
11.配置副本集
1）mongodb01初始化副本集
#连接主库
mongo --port 28010
rs.add("10.0.0.83:28020")
rs.addArb("10.0.0.82:28030")
2）mongodb02初始化副本集
#连接主库
mongo --port 28010
rs.add("10.0.0.81:28020")
rs.addArb("10.0.0.83:28030")
3）mongodb03初始化副本集
#连接主库
mongo --port 28010
rs.add("10.0.0.82:28020")
rs.addArb("10.0.0.81:28030")
4）检查所有节点副本集状态
#三台主节点
mongo --port 28010
rs.status()
rs.isMaster()
12.配置config server
1）创建目录
2）配置config server
vim /server/mongodb/config/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/config/log/mongodb.log

storage:
  journal:
    enabled: true
  dbPath: /server/mongodb/config/data/
  directoryPerDB: true

  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
      directoryForIndexes: true
    collectionConfig:
      blockCompressor: zlib
    indexConfig:
      prefixCompression: true

processManagement:
  fork: true
  pidFilePath: /server/mongodb/config/pid/mongod.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 40000
  bindIp: 127.0.0.1,10.0.0.81

replication:
  replSetName: configset

sharding:
  clusterRole: configsvr
3）启动
/usr/local/mongodb/bin/mongod -f /server/mongodb/config/conf/mongo.conf
4）mongodb01上初始化副本集
mongo --port 40000

rs.initiate({
  _id:"configset", 
  1configsvr: true,
  members:[
    {_id:0,host:"10.0.0.51:40000"},
    {_id:1,host:"10.0.0.52:40000"},
    {_id:2,host:"10.0.0.53:40000"},
  ]})
5）检查
    rs.status()
    rs.isMaster()
13.配置mongos
1）创建目录
2）配置mongos
vim /server/mongodb/mongos/conf/mongo.conf

systemLog:
  destination: file 
  logAppend: true 
  path: /server/mongodb/mongos/log/mongos.log

processManagement:
  fork: true
  pidFilePath: /server/mongodb/mongos/pid/mongos.pid
  timeZoneInfo: /usr/share/zoneinfo

net:
  port: 60000
  bindIp: 127.0.0.1,10.0.0.81

sharding:
  configDB: 
    configset/10.0.0.81:40000,10.0.0.82:40000,10.0.0.83:40000
3）启动
/usr/local/mongodb/bin/mongod -f /server/mongodb/mongos/conf/mongo.conf
4）添加分片成员
#登录mongos
mongo --port 60000

#添加成员
use admin
db.runCommand({addShard:'shard1/10.0.0.81:28100,10.0.0.83:28200,10.0.0.82:28300'})
db.runCommand({addShard:'shard2/10.0.0.82:28100,10.0.0.81:28200,10.0.0.83:28300'})
db.runCommand({addShard:'shard3/10.0.0.83:28100,10.0.0.82:28200,10.0.0.81:28300'})
5）查看分片信息
    db.runCommand( { listshards : 1 } )
14.配置区间分片
1）区间分片
#数据库开启分片
mongo --port 60000
use admin 

#指定库开启分片
db.runCommand( { enablesharding : "test" } )
2）创建集合索引
mongo --port 60000 
use test
db.range.ensureIndex( { id: 1 } )
3）对集合开启分片，片键是id
use admin
db.runCommand( { shardcollection : "test.range",key : {id: 1} } )
4）插入测试数据
use test
for(i=1;i<10000;i++){ db.range.insert({"id":i,"name":"shanghai","age":28,"date":new Date()}); }
db.range.stats()
db.range.count()
15.设置hash分片
#数据库开启分片
mongo --port 60000
use admin
db.runCommand( { enablesharding : "testhash" } )
1）集合创建索引
use testhash
db.hash.ensureIndex( { id: "hashed" } )
2）集合开启哈希分片
use admin
sh.shardCollection( "testhash.hash", { id: "hashed" } )
3）生成测试数据
use testhash
for(i=1;i<10000;i++){ db.hash.insert({"id":i,"name":"shanghai","age":70}); }
4）验证数据
分片验证
#mongodb01
mongo --port 28010
use testhash
db.hash.count()
33755

#mongodb01
mongo --port 28010
use testhash
db.hash.count()
33142

#mongodb01
mongo --port 28010
use testhash
db.hash.count()
33102
16.分片集群常用管理命令
1.列出分片所有详细信息
    db.printShardingStatus()
    sh.status()

2.列出所有分片成员信息
    use admin
    db.runCommand({ listshards : 1})

3.列出开启分片的数据库
    use config
    db.databases.find({"partitioned": true })

4.查看分片的片键
    use config
    db.collections.find().pretty()

### mongo配置密码做副本集
openssl rand -base64 123 > /server/mongodb/mongo.key
chown -R mongod.mongod /server/mongodb/mongo.key
chmod -R 600 /server/mongodb/mongo.key

scp -r /server/mongodb/mongo.key 192.168.1.82:/server/mongodb/
scp -r /server/mongodb/mongo.key 192.168.1.83:/server/mongodb/