### scp
```
## 远程copy，scp
scp 源文件 目标
 
# 推文件
scp /tmp/yum.log root@10.0.0.41:/root
# 推目录
scp -r /etc root@10.0.0.41:/root
 
# 拉文件
scp root@10.0.0.41:/root /tmp/yum.log
# 拉目录
scp -r root@10.0.0.41:/root /etc
 
ssh:22
ftp:21
rsync:873
 
C/S 架构：
Client/Server
 
B/S
Browser/Server
 
端口的范围：1-65535
```

### rsync
```
rsync   同步工具，把一台机器上的文件传输到另一台。
# 特性  rsync可以实现删除文件和目录的功能，相当于rm命令+scp命令+cp命令，但是优于他们的每一个命令相加。
 
# 选项：
-a           #归档模式传输, 等于-tropgDl
-v           #详细模式输出, 打印速率, 文件数量等
-z           #传输时进行压缩以提高效率
--delete     #让目标目录和源目录数据保持一致
--password-file=xxx #使用密码文件
 
------------------- -a 包含 ------------------
-r           #递归传输目录及子目录，即目录下得所有目录都同样传输。
-t           #保持文件时间信息
-o           #保持文件属主信息
-p           #保持文件权限
-g           #保持文件属组信息
-l           #保留软连接
-D           #保持设备文件信息
----------------------------------------------
 
-L           #保留软连接指向的目标文件
-e           #使用的信道协议,指定替代rsh的shell程序
--exclude=PATTERN   #指定排除不需要传输的文件模式
--exclude-from=file #文件名所在的目录文件
 
 
rsync备份类型
 
- 全量备份（支持）
# 全部备份一遍
- 增量备份（支持）
# 相对上次备份增加的部分
- 差异备份（不支持差异备份）
# 相对上次全量备份增加的部分
 
 
Rsync的传输模式
 
- 本地传输模式（命令cp）
# 语法
Local:  rsync [OPTION...] SRC... [DEST]
 
- 远程传输模式（命令scp）
# 这种模式是借助SSH的通道进行传输（ssh端口）
# 拉文件：rsync [选项...] 用户名@主机IP:路径 本地文件或目录
# 推文件：rsync [选项...] 本地文件或目录 用户名@主机IP:路径
 
- 守护进程模式（服务）
# 以服务的形式使用rsync，配置配置文件
# 拉文件：rsync [-avz] zls_bak@10.0.0.41::[模块] 源文件 目标
# 推文件：rsync [-avz] 源文件 zls_bak@10.0.0.41::[模块] 目标
```

### 守护进程模式详解

主机名	wanIP	lanIP	角色
web01	10.0.0.7	172.16.1.7	客户端
backup	10.0.0.41	172.16.1.41	服务端

```
######### 安装服务端
# 服务端：把备份的文件放在谁的磁盘上，谁就是服务端
 
1. 安装rsync
yum install -y rsync
2.修改配置文件（一般来说是以.conf 或 .cnf 或 .cfg结尾）
vim /etc/rsyncd.conf
## 指定进程启动uid
uid = rsync
## 指定进程启动gid
gid = rsync
## rsync服务的端口
port = 873
## 无需让rsync以root身份运行，允许接收文件的完整属性
fake super = yes
## 禁锢指定的目录
use chroot = no
## 最大连接数
max connections = 200
## 超时时间
timeout = 600
## 忽略错误
ignore errors
## 不只读（可读可写）
read only = false
## 不允许别人查看模块名
list = false
## 传输文件的用户
auth users = rsync_backup
## 传输文件的用户和密码文件
secrets file = /etc/rsync.passwd
## 日志文件
log file = /var/log/rsyncd.log
#####################################
## 模块名
[zls]
## 注释，没啥用
comment = 123
## 备份的目录
path = /backup
 
 
3.根据配置文件内容，创建出来需要的用户，目录，密码文件...
# 3.1 创建用户
useradd rsync -s /sbin/nologin -M
# 3.2 创建备份目录
mkdir /backup
# 3.3 修改属组和属主
chown -R rsync.rsync /backup/
# 3.4 创建用户名和密码存放的文件
vim /etc/rsync.passwd
zls_bak:123
或
echo 'zls_bak:123' > /etc/rsync.passwd
# 3.5 修改密码文件的权限为600
chmod 600 /etc/rsync.passwd
 
4.启动服务并且加入开机自启
systemctl start rsyncd
systemctl enable rsyncd
 
5.检测端口
netstat -lntup | grep 873
      
6.检测进程
ps -ef | grep [r]sync
 
 
######## 安装客户端
# rsync 客户端不用修改配置文件
 
1.安装rsync
yum install -y rsync
2.客户端需要创建一个密码文件
vim /etc/rsync.pass
123
3.修改密码文件的权限为600
chmod 600 /etc/rsync.pass 
## 二三步在写脚本时可以省略  换为给变量赋值
# export RSYNC_PASSWORD=密码
4.从客户端往服务端推送重要备份文件
rsync [-avz] 源文件 zls_bak@10.0.0.41::[模块]
rsync -avz /etc/shadow zls_bak@10.0.0.41::zls --password-file=/etc/rsync.pass
 
 
 
# rsync 无差异同步
--delete
 
- 拉取方式：
1.远端有什么，本地就有什么，远端没有的本地有也要删处。 
2.客户端目录数据可能丢失
 
- 推送方式：
1.本地有什么，远程就有什么，本地没有的远端有也要删除。 
2.服务器端的目录数据可能丢失。
# rsync 限速
--bwlimit=
sending incremental file list
# 限制了磁盘的I/O 吞吐量
```

### inotify

```bash
yum install -y inotify-tools
# inotifywait  监控命令
# 选项：
-m 持续监控
-r 递归
-q 静默，仅打印时间信息
--timefmt 指定输出时间格式
--format 指定事件输出格式
%Xe 事件
%w 目录
%f 文件
-e 指定监控的事件
access 访问
modify 内容修改
attrib 属性修改
close_write 修改真实文件内容
open 打开
create 创建
delete 删除
umount 卸载
 
# 例： 
inotifywait  -mrq  --format '%Xe  %w  %f' -e create,modify,delete,attrib,close_write  + 文件或目录

```

### rsync 守护进程模式工作流程

```bash
当执行`rsync -avz /etc/passwd zls_bak@10.0.0.41::zls`的一瞬间
 
# 网络是否可以通讯，客户端是否能连接服务端的IP和端口（通过）
1. 网络
2. 端口
3. 防火墙
4. selinux
 
# 验证用户名和密码（通过）
1. 检查配置文件的用户名
2. 检查密码文件里的内容
3. 检查密码文件的权限是不是600
4. 检查模块名：配置文件中指定的[模块名]
 
# 检查机器上是否有配置文件中的uid指定的用户
 
# 检查对应模块名下面指定的path目录，是否有权限（uid指定用户的权限）
```

### sersync配置文件

```bash
<?xml version="1.0" encoding="ISO-8859-1"?>
<head version="2.5">
    <host hostip="localhost" port="8008"></host>
    <debug start="false"/>
    <fileSystem xfs="false"/>
    <filter start="false">
	<exclude expression="(.*)\.svn"></exclude>
	<exclude expression="(.*)\.gz"></exclude>
	<exclude expression="^info/*"></exclude>
	<exclude expression="^static/*"></exclude>
    </filter>
    <inotify>
	<delete start="true"/>
	<createFolder start="true"/>
	<createFile start="false"/>
	<closeWrite start="true"/>
	<moveFrom start="true"/>
	<moveTo start="true"/>
	<attrib start="true"/>
	<modify start="true"/>
    </inotify>
 
    <sersync>
	<!-- 客户端需要监控的目录 -->
	<localpath watch="/user_upload">
 
	    <!-- rsync服务端的IP 和 name:模块 -->
	    <remote ip="10.0.0.41" name="nfs"/>
	    <!--<remote ip="192.168.8.39" name="tongbu"/>-->
	    <!--<remote ip="192.168.8.40" name="tongbu"/>-->
	</localpath>
	<rsync>
	    <!-- rsync命令执行的参数 -->
	    <commonParams params="-az"/>
            <!-- rsync认证start="true" users="rsync指定的匿名用户" passwordfile="指定一个密码文件的位置权限必须600" -->
	    <auth start="true" users="nfs_bak" passwordfile="/etc/rsync.pas"/>
	    <userDefinedPort start="false" port="874"/><!-- port=874 -->
	    <timeout start="false" time="100"/><!-- timeout=100 -->
	    <ssh start="false"/>
	</rsync>
	<failLog path="/tmp/rsync_fail_log.sh" timeToExecute="60"/><!--default every 60mins execute once-->
	<crontab start="false" schedule="600"><!--600mins-->
	    <crontabfilter start="false">
		<exclude expression="*.php"></exclude>
		<exclude expression="info/*"></exclude>
	    </crontabfilter>
	</crontab>
	<plugin start="false" name="command"/>
    </sersync>
 
    <plugin name="command">
	<param prefix="/bin/sh" suffix="" ignoreError="true"/>	<!--prefix /opt/tongbu/mmm.sh suffix-->
	<filter start="false">
	    <include expression="(.*)\.php"/>
	    <include expression="(.*)\.sh"/>
	</filter>
    </plugin>
 
    <plugin name="socket">
	<localpath watch="/opt/tongbu">
	    <deshost ip="192.168.138.20" port="8009"/>
	</localpath>
    </plugin>
    <plugin name="refreshCDN">
	<localpath watch="/data0/htdocs/cms.xoyo.com/site/">
	    <cdninfo domainname="ccms.chinacache.com" port="80" username="xxxx" passwd="xxxx"/>
	    <sendurl base="http://pic.xoyo.com/cms"/>
	    <regexurl regex="false" match="cms.xoyo.com/site([/a-zA-Z0-9]*).xoyo.com/images"/>
	</localpath>
    </plugin>
</head>
```