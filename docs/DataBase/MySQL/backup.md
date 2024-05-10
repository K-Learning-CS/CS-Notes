## MySQL数据库备份方案


### 一、非生产环境

#### 服务器规划

|       ip        |   描述   | cpu | 内存 | 硬盘 |   硬盘    |    备注    |
|:---------------:|:------:|-----|:--:| :--: |:-------:|:--------:|
| 192.168.108.130 | slave  | 8C  | 8G | 100G |  200G   |  备份服务器   |
|  192.168.1.82   | slave  | 4C  | 4G | 100G | 80G SSD | 备用master |
|  192.168.1.89   | backup | 4C  | 4G | 100G | 80G SSD |  备用服务器   |


#### 备份策略

- 全量备份：
  1. 备份时间：每天凌晨一点
  2. 备份工具：Mydumper
  3. 备份类型：逻辑备份、热备
  4. 备份周期：每天凌晨两点
  5. 保存时间：7天 
  6. 备份机器：192.168.1.81 
  7. 保存地址：/export/mysql/backup/mydumper-time.tar.gz



- 增量备份：
  1. 备份时间：实时
  2. 备份工具：Mysqlbinlog
  3. 备份类型：裸文件备份、热备
  4. 备份周期：实时
  5. 保存时间：7天 
  6. 备份机器：192.168.108.130 
  7. 保存地址：/export/mysql/backup/192.168.1.81/binlogs/

#### 备份细节

1. 全量备份

* 在 192.168.1.81 服务器上的 /export/shell/ 目录中添加此脚本
* 在定时任务中添加 `00 02 * * * /export/shell/mydumper_backup.sh > /dev/null 2>&1`


```bash

#!/bin/bash
# ---------------------------------------------------------
# 脚本名称: mydumper_backup.sh
# 描述信息: 备份指定数据库的全量备份
# ---------------------------------------------------------
######### 设置变量 #########
# 连接MySQL数据库常用变量
MYSQLHOST='192.168.1.81'
MYSQLPORT='3306'
MYSQLUSER='root'
MYSQLPASS='66666666'
# 数据库备份存放目录
TIME=$(date "+%Y.%m.%d.%H.%M")
BACKUPDIR="/export/mysql/backup/${TIME}"
[ -d ${BACKUPDIR} ] || mkdir -p ${BACKUPDIR}
######### 开始备份全量备份 #########
# --set-names 默认为 binary，会导致 json 导入失败
# --exec-threads 指定启动的线程数
mydumper -u ${MYSQLUSER} -p ${MYSQLPASS}  -h ${MYSQLHOST} --regex '^(?!(mysql|information_schema|performance_schema|sys))' --set-names utf8mb4 -G -R -E  --exec-threads 4 -o ${BACKUPDIR}

# 备份配置文件
cd ${BACKUPDIR} 
cp /etc/my.cnf .
cd ..

# 压缩
tar zcf mydumper-${TIME}.tar.gz ${TIME}
rm -rf ${TIME}

# 发送
scp mydumper-${TIME}.tar.gz 192.168.1.82:/export/mysql/backup/
scp mydumper-${TIME}.tar.gz 192.168.108.130:/export/mysql/backup/

# 清理
find /export/mysql/backup/ -type f -mtime +7 | awk '{print "rm -f "$1}' | bash
```



2. 增量备份

* 在 192.168.108.130 服务器上的 /export/shell/ 目录中添加此脚本
* 使用`nohup /export/shell/binlog_backup.sh &`启动此备份脚本

```bash

#!/bin/bash
# ---------------------------------------------------------
# 脚本名称: binlog_backup.sh
# 描述信息: 备份指定数据库的binlog日志，需要用户有REPLICATION SLAVE权限
# ---------------------------------------------------------
######### 设置变量 #########
# 连接MySQL数据库常用变量
MYSQLHOST='192.168.1.81'
MYSQLPORT='3306'
MYSQLUSER='root'
MYSQLPASS='66666666'
CMDDIR="/usr/local/mysql/bin"
FIRST_BINLOG="mysql-bin.000001"
# 数据库备份存放目录
BACKUPDIR="/export/mysql/backup"
# Binlog日志存放目录
BINLOGDIR="${BACKUPDIR}/${MYSQLHOST}/binlogs"
[ -d ${BINLOGDIR} ] || mkdir -p ${BINLOGDIR}
# MYSQLBINLOG命令
MYSQLBINLOG="${CMDDIR}/mysqlbinlog"
# 备份日志
BACKUPLOG="/var/log/backup_binlog.log"

# 停止时间
SLEEP_SECONDS=10
######### 开始备份binglog日志 #########
# 运行while循环，连接断开后等待指定时间，重新连接
while true;do
    if [ $(ls -A ${BINLOGDIR} |wc -l) -eq 0 ];then
        LAST_FILE=${FIRST_BINLOG}
    else
        LAST_FILE=$(ls ${LOCAL_BACKUP_DIR}|tail -1 |awk '{print $NF}')
    fi
    cd ${BINLOGDIR} && \
    ${MYSQLBINLOG}  --raw --read-from-remote-server --stop-never --host=${MYSQLHOST} --port=${MYSQLPORT} --user=${MYSQLUSER} --password=${MYSQLPASS} ${LAST_FILE}
    echo "$(date +"%F %T") ${MYSQLHOST}:${MYSQLPORT}备份binlog日志停止，返回代码：$?" | tee -a ${BACKUPLOG}
    echo "$(date +"%F %T") ${MYSQLHOST}:${MYSQLPORT} ${SLEEP_SECONDS}秒后再次连接并继续备份" | tee -a ${BACKUPLOG}
    sleep ${SLEEP_SECONDS}
done
```


#### 恢复策略

- 宕机恢复：
  1. 优先使用 slave，根据 slave 节点的 point，查找 binlog 中未执行的语句，将语句导出后导入 slave，并将 slave 提升为主库
  2. 启用备用 mysql 服务器，将全备导入，再将差异从 binlog 中找出并导出为 sql，而后将 sql 导入新库

- 误删恢复：
  1. 启用备用 mysql 服务器，将全备导入，再从 binlog 中找到删除语句，从前一条导出 sql，而后将 sql 导入新库，再将删除的库/表导出，导入至 master


#### 恢复细节


1. slave 提升

```bash
1. 登陆 slave 找到当前使用的 master binlog 文件，以及point
# 192.168.1.82
mysql -uroot -p
show slave status\G
#              Master_Log_File: bin-log.002066
#          Read_Master_Log_Pos: 842254119

2. 从 master 节点或者 backup 节点找到 binlog 文件，并加其需要的部分导出为 sql
# 192.168.108.130
cd /export/mysql/backup/192.168.1.81/binlogs/
# 找到结束位置
mysqlbinlog bin-log.002066 > 1.txt
tail 1.txt
# at 961496292

# 导出 sql
mysqlbinlog -d del --start-position=842254119 --stop-position=961496292 bin-log.002066 > new.sql

3. 将 sql 发送到 slave，并由 slave 执行
scp new.sql root@192.168.1.82:/export/mysql/
mysql -uroot -p
source /export/mysql/new.sql;

4. 将 slave 提升为主库，并开放读写
reset slave;

```

2. 启用新库

```bash
# 192.168.1.89

1. 停止服务，删除数据目录
systemclt stop mysqld
rm -rf /export/mysql/data/*

2. 初始化数据目录
/usr/local/mysql/bin/mysqld --initialize --user=mysql --basedir=/usr/local/mysql --datadir=/export/mysql/data
# 记录生成的随机密码

3. 启动服务并修改密码和连接
systemctl start mysqld
mysql --socket=/usr/local/mysql/mysql.sock -uroot -p
ALTER USER 'root'@'localhost' IDENTIFIED BY '<new_password>';
grant all on *.* to root@'192.168.%' identified by '<password>';
flush privileges;

4. 解压最近全备压缩包，导入全备数据
cd /export/mysql/backup
tar xf mydumper-time.tar.gz
myloader -u root -p <password> --regex '^(?!(mysql|test|information_schema|performance_schema|sys))'  -o -d backup/

5. 从 backup 中找到起始 binlog 和 point
head <time>/metadata
#File = bin-log.002063
#Position = 111029088

6. 从 bin-log.002063 到结束的 binlog 中导出 sql 数据
mysqlbinlog -d del --start-position=111029088 bin-log.002063 > 002063.sql
mysqlbinlog  bin-log.002064 > 002064.sql
mysqlbinlog  bin-log.002065 > 002065.sql
mysqlbinlog  bin-log.002066 > 002066.sql
# 如果设及到删库删表恢复，请找出删库删表语句前的位置点，并指定结束位置

7. 导入 sql
mysql -uroot -p
source /export/mysql/002063.sql;
source /export/mysql/002064.sql;
source /export/mysql/002065.sql;
source /export/mysql/002066.sql;

# 误删恢复
# 恢复完成后导出库或表的 sql 数据，将其导入到原库即可
```

### 二、生产环境

#### 服务器规划

|       ip        |   描述   | cpu | 内存 | 硬盘 |  硬盘  |    备注    |
|:---------------:|:------:|-----|:--:| :--: |:----:|:--------:|
| 192.168.108.130 | master | 8C  | 8G | 100G | 200G | 内网生产数据库  |
|       10.       |  ecs   | 2C  | 2G | 100G |      | 云环境备份ecs |

#### 备份策略

- 同非生产

#### 备份细节

- 同非生产

#### 恢复策略

- 每天将前一天全备同步至非生产环境

#### 恢复细节

- 同非生产