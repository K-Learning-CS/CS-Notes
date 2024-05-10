# MHA

### MHA介绍
简介
MHA能够在较短的时间内实现自动故障检测和故障转移，通常在10-30秒以内;在复制框架中，MHA能够很好地解决复制过程中的数据一致性问题，由于不需要在现有的replication中添加额外的服务器，仅需要一个manager节点，而一个Manager能管理多套复制，所以能大大地节约服务器的数量;另外，安装简单，无性能损耗，以及不需要修改现有的复制部署也是它的优势之处。
 
MHA还提供在线主库切换的功能，能够安全地切换当前运行的主库到一个新的主库中(通过将从库提升为主库),大概0.5-2秒内即可完成。
 
MHA由两部分组成：MHA Manager（管理节点）和MHA Node（数据节点）。MHA Manager可以独立部署在一台独立的机器上管理多个Master-Slave集群，也可以部署在一台Slave上。当Master出现故障时，它可以自动将最新数据的Slave提升为新的Master,然后将所有其他的Slave重新指向新的Master。整个故障转移过程对应用程序是完全透明的。
#在切换过程我们可以查看日志
原理
1.把宕机的master二进制日志保存下来。
2.找到binlog位置点最新的slave。
3.在binlog位置点最新的slave上用relay-log（差异日志）修复其它slave。
4.将宕机的master上保存下来的二进制日志恢复到含有最新位置点的slave上。
5.将含有最新位置点binlog所在的slave提升为master。
6.将其它slave重新指向新提升的master，并开启主从复制。
架构
1.MHA属于C/S结构
2.一个manager节点可以管理多套集群
3.集群中所有的机器都要部署node节点
4.node节点才是管理集群机器的
5.manager节点通过ssh连接node节点，管理
6.manager可以部署在集群中除了主库以外的任意一台机器上
工具
1）manager节点的工具
#解压tar包，查看
[root@db01 ~]# ll mha4mysql-manager-0.56/bin/
#检查主从状态
masterha_check_repl
#检查ssh连接（配置免密）
masterha_check_ssh
#检查MHA状态
masterha_check_status
#删除死掉机器的配置
masterha_conf_host
    [server2]
    hostname=10.0.0.52
    port=3306
 
    [server3]
    hostname=10.0.0.53
    port=3306
#启动程序
masterha_manager
#检测master是否宕机
masterha_master_monitor
#手动故障转移
masterha_master_switch
#建立TCP连接从远程服务器
masterha_secondary_check
#关闭进程的程序
masterha_stop
2）node节点工具
#解压node安装包
 
[root@db01 ~]# ll mha4mysql-node-0.56/bin/
#对比relay-log
apply_diff_relay_logs
#防止回滚事件
filter_mysqlbinlog
#删除relay-log
purge_relay_logs
#保存binlog
save_binary_logs
### MHA优点
1）Masterfailover and slave promotion can be done very quickly
自动故障转移快
 
2）Mastercrash does not result in data inconsistency
主库崩溃不存在数据一致性问题
 
3）Noneed to modify current MySQL settings (MHA works with regular MySQL)
不需要对当前mysql环境做重大修改
 
4）Noneed to increase lots of servers
不需要添加额外的服务器(仅一台manager就可管理上百个replication)
 
5）Noperformance penalty
性能优秀，可工作在半同步复制和异步复制，当监控mysql状态时，仅需要每隔N秒向master发送ping包(默认3秒)，所以对性能无影响。你可以理解为MHA的性能和简单的主从复制框架性能一样。
 
6）Works with any storage engine
只要replication支持的存储引擎，MHA都支持，不会局限于innodb
### 搭建MHA
保证主从的状态
#主库
mysql> show master status;
#从库
mysql> show slave status\G
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
#从库
mysql> show slave status\G
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
部署MHA之前配置
1.关闭从库删除relay-log的功能
relay_log_purge=0
 
2.配置从库只读
read_only=1
 
3.从库保存binlog
log_slave_updates
#禁用自动删除relay log 功能
mysql> set global relay_log_purge = 0;
#设置只读
mysql> set global read_only=1;
#编辑配置文件
[root@mysql-db02 ~]# vim /etc/my.cnf
#在mysqld标签下添加
[mysqld]
#禁用自动删除relay log 永久生效
relay_log_purge = 0
配置
1）主库配置
[mysqld]
server_id=1
binlog_format=row
basedir=/usr/local/mysql
datadir=/usr/local/mysql/data
socket=/var/lib/mysql/mysql.sock
log_bin=/usr/local/mysql/data/mysql-bin
log_err=/usr/local/mysql/data/mysql.err
 
gtid_mode=on
enforce_gtid_consistency
log-slave-updates
relay_log_purge=0
read_only=1
skip-name-resolve
 
[mysql]
socket=/var/lib/mysql/mysql.sock
2）从库01配置
[mysqld]
server_id=2
binlog_format=row
basedir=/usr/local/mysql
datadir=/usr/local/mysql/data
socket=/var/lib/mysql/mysql.sock
log_bin=/usr/local/mysql/data/mysql-bin
log_err=/usr/local/mysql/data/mysql.err
 
gtid_mode=on
enforce_gtid_consistency
log-slave-updates
relay_log_purge=0
read_only=1
skip-name-resolve
 
[mysql]
socket=/var/lib/mysql/mysql.sock
3）从库02配置
[mysqld]
server_id=3
binlog_format=row
basedir=/usr/local/mysql
datadir=/usr/local/mysql/data
socket=/var/lib/mysql/mysql.sock
log_bin=/usr/local/mysql/data/mysql-bin
log_err=/usr/local/mysql/data/mysql.err
 
gtid_mode=on
enforce_gtid_consistency
log-slave-updates
relay_log_purge=0
read_only=1
skip-name-resolve
 
[mysql]
socket=/var/lib/mysql/mysql.sock
### 部署MHA
~~~bash
1）安装依赖（所有机器）
yum install perl-DBD-MySQL -y
2）安装manager依赖（manager机器）
yum install -y perl-Config-Tiny epel-release perl-Log-Dispatch perl-Parallel-ForkManager perl-Time-HiRes
4）部署manager节点
rz mha4mysql-manager-0.56-0.el6.noarch.rpm
yum localinstall -y mha4mysql-manager-0.56-0.el6.noarch.rpm 
5）配置MHA
#创建MHA配置目录
[root@db03 ~]# mkdir -p /service/mha
 
#配置MHA
[root@db03 ~]# vim /service/mha/app1.cnf
[server default]
#指定日志存放路径
manager_log=/service/mha/manager
#指定工作目录
manager_workdir=/service/mha/app1
#binlog存放目录
master_binlog_dir=/usr/local/mysql/data
#MHA管理用户
user=mha
#MHA管理用户的密码
password=mha
#检测时间
ping_interval=2
#主从用户
repl_user=rep
#主从用户的密码
repl_password=123
#ssh免密用户
ssh_user=root
 
[server1]
hostname=172.16.1.51
port=3306
 
[server2]
#candidate_master=1
#check_repl_delay=0
hostname=172.16.1.52
port=3306
 
[server3]
hostname=172.16.1.53
port=3306
 
#设置为候选master，如果设置该参数以后，发生主从切换以后将会将此从库提升为主库，即使这个主库不是集群中事件最新的slave。
candidate_master=1
#默认情况下如果一个slave落后master 100M的relay logs的话，MHA将不会选择该slave作为一个新的master，因为对于这个slave的恢复需要花费很长时间，通过设置check_repl_delay=0,MHA触发切换在选择一个新的master的时候将会忽略复制延时，这个参数对于设置了candidate_master=1的主机非常有用，因为这个候选主在切换的过程中一定是新的master
check_repl_delay=0
 
 
# 配置文件
[root@db02 ~]# tree /apps
/apps
└── mha
    ├── app1
    │   ├── app1.failover.complete
    │   └── manager
    ├── app1.conf
    └── manager.log
vim /apps/mha/app1.conf 
 
[server default]
manager_log=/apps/mha/app1/manager
manager_workdir=/apps/mha/app1
master_binlog_dir=/usr/local/mysql/data/
password=mha
ping_interval=2
repl_password=123
repl_user=rep
ssh_user=root
user=mha
 
[server1]
hostname=172.16.1.51
port=3306
 
 
[server2]
candidate_master=1
check_repl_delay=0
hostname=172.16.1.52
port=3306
 
[server3]
hostname=172.16.1.53
port=3306
6）创建MHA管理用户
#主库执行即可
mysql> grant all on *.* to mha@'172.16.1.%' identified by 'mha';
Query OK, 0 rows affected (0.03 sec)
7）ssh免密（三台机器每一台都操作一下内容
#创建秘钥对
ssh-keygen -t dsa -P '' -f ~/.ssh/id_dsa >/dev/null 2>&1
#发送公钥，包括自己
ssh-copy-id -i /root/.ssh/id_dsa.pub root@172.16.1.51
ssh-copy-id -i /root/.ssh/id_dsa.pub root@172.16.1.52
ssh-copy-id -i /root/.ssh/id_dsa.pub root@172.16.1.53
8）检测MHA状态
#检测主从
masterha_check_repl --conf=/apps/mha/app1.conf
 
MySQL Replication Health is OK.
 
#检测ssh
masterha_check_ssh --conf=/apps/mha/app1.conf
 
Mon Jul 27 11:40:06 2020 - [info] All SSH connection tests passed successfully.
9）启动MHA
#启动
[root@db03 ~]# nohup masterha_manager --conf=/apps/mha/app1.conf --remove_dead_master_conf --ignore_last_failover < /dev/null > /apps/mha/manager.log 2>&1 &
 
nohup ... &   					#后台启动
masterha_manager 				#启动命令
--conf=/apps/mha/app1.conf		 #指定配置文件
--remove_dead_master_conf 		 #移除挂掉的主库配置
--ignore_last_failover 			 #忽略最后一次切换
< /dev/null > /apps/mha/manager.log 2>&1
 
#MHA保护机制：
	1.MHA主库切换后，8小时内禁止再次切换
	2.切换后会生成一个锁文件，下一次启动MHA需要检测该文件是否存在
~~~
### 恢复MHA
修复数据库
systemctl start mysqld.service
恢复主从
#将恢复的数据库当成新的从库加入集群
#找到binlog位置点
[root@db03 ~]# grep 'CHANGE MASTER' /service/mha/manager | awk -F: 'NR==1 {print $4}'
 CHANGE MASTER TO MASTER_HOST='172.16.1.52', MASTER_PORT=3306, MASTER_AUTO_POSITION=1, MASTER_USER='rep', MASTER_PASSWORD='xxx';
#恢复的数据库执行change master to
mysql> CHANGE MASTER TO MASTER_HOST='172.16.1.52', MASTER_PORT=3306, MASTER_AUTO_POSITION=1, MASTER_USER='rep', MASTER_PASSWORD='123';
Query OK, 0 rows affected, 2 warnings (0.20 sec)
mysql> start slave;
Query OK, 0 rows affected (0.05 sec)
恢复MHA
#将恢复的数据库配置到MHA配置文件
[root@db03 ~]# vim /service/mha/app1.cnf 
......
[server1]
hostname=172.16.1.51
port=3306
 
[server2]
hostname=172.16.1.52
port=3306
 
[server3]
hostname=172.16.1.53
port=3306
......
 
#启动MHA
[root@db03 ~]# nohup masterha_manager --conf=/apps/mha/app1.conf --remove_dead_master_conf --ignore_last_failover < /dev/null > /apps/mha/manager.log 2>&1 &
### MHA主库切换
MHA主库切换机制
1.读取配置中的指定优先级
	candidate_master=1
	check_repl_delay=0
2.如果数据量不同，数据量多的为主库
3.如果数据量相同，按照主机标签，值越小优先级越高
主机标签优先级测试
#配置MHA
[root@db03 ~]# vim /service/mha/app1.cnf
......
[serverc]
hostname=172.16.1.53
port=3306
 
[serverb]
hostname=172.16.1.52
port=3306
 
[serverd]
hostname=172.16.1.51
port=3306
......
 
#停掉主库，查看切换
mysql> show slave status\G
                  Master_Host: 172.16.1.52
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
指定优先级测试
#配置优先级
[root@db03 ~]# vim /service/mha/app1.cnf 
......
[server3]
candidate_master=1
check_repl_delay=0
hostname=172.16.1.53
port=3306
......
 
#停掉主库，查看切换
mysql> show slave status\G
                  Master_Host: 172.16.1.53
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
数据量不同时切换的优先级
1）主库建库
mysql> use test
Database changed
mysql> create table linux9(id int not null primary key auto_increment,name varchar(10));
Query OK, 0 rows affected (0.31 sec)
2）脚本循环插入数据
[root@db03 ~]# vim insert.sh 
#/bin/bash
while true;do
    mysql -e "use test;insert into linux9(name) values('lhd')"
    sleep 1
done
 
[root@db03 ~]# sh insert.sh 
3）停掉一台机器的IO线程
#此时db03是主库，停掉db01的IO线程
mysql> stop slave io_thread;
4）停掉db03的主库
[root@db03 ~]# systemctl stop mysqld
5）查看主库切换
mysql> show slave status\G
*************************** 1. row ***************************
               Slave_IO_State: Waiting for master to send event
                  Master_Host: 172.16.1.52
                  Master_User: rep
                  Master_Port: 3306
                Connect_Retry: 60
              Master_Log_File: mysql-bin.000008
          Read_Master_Log_Pos: 67437
               Relay_Log_File: db01-relay-bin.000002
                Relay_Log_Pos: 14904
        Relay_Master_Log_File: mysql-bin.000008
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
### 如果主库宕机，binlog如何保存
配置MHA实时备份binlog
[root@db03 ~]# vim /service/mha/app1.cnf
[root@db03 ~]# vim /service/mha/app1.cnf
......
[binlog1]
no_master=1
hostname=172.16.1.53
#不能跟当前机器数据库的binlog存放目录一样
master_binlog_dir=/root/binlog/
创建binlog存放目录
mkdir binlog
手动执行备份binlog
#进入该目录
[root@db03 ~]# cd /root/binlog/
#备份binlog
[root@db03 binlog]# mysqlbinlog  -R --host=172.16.1.52 --user=mha --password=mha --raw --stop-never mysql-bin.000001 &
重启MHA
[root@db03 binlog]# masterha_stop --conf=/service/mha/app1.cnf
Stopped app1 successfully.
 
[root@db03 binlog]# nohup masterha_manager --conf=/service/mha/app1.cnf --remove_dead_master_conf --ignore_last_failover < /dev/null > /service/mha/manager.log 2>&1 &
主库添加数据查看binlog
#主库
mysql> create database mha;
Query OK, 1 row affected (0.01 sec)
 
[root@db02 ~]# ll /usr/local/mysql/data/mysql-bin.000008 
-rw-rw---- 1 mysql mysql 67576 Jul 28 10:33 /usr/local/mysql/data/mysql-bin.000008
 
#MHA机器查看binlog
[root@db03 binlog]# ll
total 96
-rw-rw---- 1 root root   852 Jul 28 10:30 mysql-bin.000001
-rw-rw---- 1 root root   214 Jul 28 10:30 mysql-bin.000002
-rw-rw---- 1 root root   214 Jul 28 10:30 mysql-bin.000003
-rw-rw---- 1 root root   214 Jul 28 10:30 mysql-bin.000004
-rw-rw---- 1 root root   465 Jul 28 10:30 mysql-bin.000005
-rw-rw---- 1 root root   214 Jul 28 10:30 mysql-bin.000006
-rw-rw---- 1 root root   214 Jul 28 10:30 mysql-bin.000007
-rw-rw---- 1 root root 67576 Jul 28 10:33 mysql-bin.000008
### VIP漂移
VIP漂移的两种方式
1.keeplaived的方式
2.MHA自带的脚本进行VIP漂移
配置MHA读取VIP漂移脚本
#编辑配置文件
[root@db03 ~]# vim /service/mha/app1.cnf
#在[server default]标签下添加
[server default]
#使用MHA自带脚本
master_ip_failover_script=/service/mha/master_ip_failover
编写脚本
#默认脚本存放在
[root@db01 ~]# ll mha4mysql-manager-0.56/samples/scripts/
total 32
-rwxr-xr-x 1 4984 users  3648 Apr  1  2014 master_ip_failover
 
#上传现成的脚本
 
#编辑脚本
[root@db03 mha]# vim master_ip_failover
.......
my $vip = '172.16.1.50/24';
my $key = '1';
my $ssh_start_vip = "/sbin/ifconfig eth1:$key $vip";
my $ssh_stop_vip = "/sbin/ifconfig eth1:$key down";
......
手动绑定VIP
[root@db01 ~]# ifconfig eth1:1 172.16.1.50/24
[root@db01 ~]# ip a
3: eth1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:0c:29:2c:76:88 brd ff:ff:ff:ff:ff:ff
    inet 172.16.1.51/24 brd 172.16.1.255 scope global noprefixroute eth1
       valid_lft forever preferred_lft forever
    inet 172.16.1.50/24 brd 172.16.1.255 scope global secondary eth1:1
       valid_lft forever preferred_lft forever
    inet6 fe80::de8c:34e1:563e:9f2b/64 scope link noprefixroute 
       valid_lft forever preferred_lft forever
 
#解绑VIP
[root@db01 ~]# ifconfig eth1:1 [172.16.1.50] down
重启MHA
#启动MHA
[root@db03 mha]# nohup masterha_manager --conf=/service/mha/app1.cnf --remove_dead_master_conf --ignore_last_failover < /dev/null > /service/mha/manager.log 2>&1 &
 
#启动失败：
	1.检查配置文件语法是否正确
	2.授权是否正确
		[root@db03 mha]# chmod 755 master_ip_failover
	3.脚本格式要正确
		[root@db03 mha]# dos2unix master_ip_failover 
		dos2unix: converting file master_ip_failover to Unix format ...
        4.去掉VIP飘逸的脚本后，MHA切换时仍然报错说要读取脚本
	#报错：
	master_ip_failover_script or purge_relay_log 在同一台机器运行
	#解决：
	重新写一个MHA配置文件
测试VIP漂移
#停止主库
[root@db01 ~]# systemctl stop mysqld.service
 
#查看切换成主库的ip地址
[root@db02 ~]# ip a
3: eth1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:0c:29:3e:56:1f brd ff:ff:ff:ff:ff:ff
    inet 172.16.1.52/24 brd 172.16.1.255 scope global noprefixroute eth1
       valid_lft forever preferred_lft forever
    inet 172.16.1.50/24 brd 172.16.1.255 scope global secondary eth1:1
       valid_lft forever preferred_lft forever
    inet6 fe80::de8c:34e1:563e:9f2b/64 scope link tentative noprefixroute dadfailed 
       valid_lft forever preferred_lft forever
    inet6 fe80::a3a8:499:ce26:b3be/64 scope link noprefixroute 
       valid_lft forever preferred_lft forever
使用脚本恢复MHA
vim /root/ansible/roles/mysql/files/mha_backup.sh 
```bash
#!/bin/bash
# 指定manager IP
maneger='172.16.1.52'
# 指定 备份 binlog 数据库 IP
no_master='172.16.1.53'
# 指定mha日志
mha_log='/apps/mha/app1/manager'
# 指定mha配置文件
mha_conf='/apps/mha/app1.conf'
# 指定mha备份配置文件
mha_conf_buckup='/apps/mha/app1.bak'
 
# 获取slave设置
slave_set=`ssh $maneger "grep 'CHANGE MASTER TO' $mha_log | tail -1 |  sed 's#xxx#123#g' " | awk -F: '{print $4}'`
# 获取宕机主库 IP
died_master_ip=`ssh $maneger "grep 'is down\!' $mha_log " | tail -1 | awk -F'[ (]' '{print $2}'`
# 获取新主库 IP
new_master_ip=`ssh $maneger "grep 'Master failover to' $mha_log " | tail -1 | awk -F'[ (]' '{print $4}'`
# 获取当前所在主机 IP
now_ip=`ip a s eth1 | awk -F '[ /]+' 'NR==3{print $3}'`
 
# 如果当前主机是宕机主库
if [ $now_ip = $died_master_ip ];then
 
    systemctl restart mysql
    # 等待数据库重启  不等待会报错
    sleep 3
    # 重新建立主从复制
    mysql -e "$slave_set;start slave" || mysql -e "stop slave;reset slave all;$slave_set;start slave"
    # 恢复mha配置文件
    ssh $maneger "\cp $mha_conf_buckup  $mha_conf"
    # 重新获取binlog并备份
    ssh $no_master "cd /root/binlog/ ;  mysqlbinlog  -R --host=${new_master_ip} --user=mha --password=mha --raw --stop-never mysql-bin.000001 &>/dev/null  &"
    # 重启mha
    ssh $maneger "nohup masterha_manager --conf=$mha_conf --remove_dead_master_conf --ignore_last_failover < /dev/null > /apps/mha/manager.log 2>&1 &"
 
# 如果不是
else
    # 将脚本发送至宕机主库
    scp  -o stricthostkeychecking=no /root/mha_backup.sh ${died_master_ip}:/root  &>/dev/null
    # 执行脚本
    ssh  ${died_master_ip}   'sh /root/mha_backup.sh'   &>/dev/null
 
fi
```