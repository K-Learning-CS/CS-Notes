### 服务区分

~~~bash
静态服务器
## 静态页： 不调用数据库的页面
# nginx
# apache
# IIS
# lighttpd
# tengine
# openresty-nginx
 
动态服务器
## 动态业：调用数据库的页面
# tomcat
# resin
# php
# weblogic
# jboss
~~~
### nginx安装

~~~bash

------------------------- 安装nginx ---------------------------
 
1.修改官方源
vim /etc/yum.repos.d/nginx.repo
[nginx-stable]
name=nginx stable repo
baseurl=http://nginx.org/packages/centos/$releasever/$basearch/
gpgcheck=1
enabled=1
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true
2.安装
yum install -y nginx
3.启动，并添加开机自启
systemctl start nginx
systemctl enable nginx
 
----------------- 检测nginx是否启动，安装成功 --------------------
 
1.检查端口
netstat -lntup |grep 80
2.检查进程
ps -ef|grep [n]ginx
3.检查nginx版本
nginx -v
4.检查安装的模块
nginx -V
 
------------------------ nginx相关操作 --------------------------
 
# systemd管理
1.停止
systemctl stop nginx
2.启动
systemctl start nginx
3.重启
systemctl restart nginx
4.重新加载配置文件
systemctl reload nginx
# 二进制程序管理
1.停止
nginx -s stop
2.启动
nginx
3.重新加载
nginx -s reload
~~~

### nginx相关文件
~~~bash
# 1.nginx的主配置文件
ll /etc/nginx/nginx.conf
 
# 2.nginx代理文件
-rw-r--r-- 1 root root 1007 Apr 21 23:07 fastcgi_params # php
-rw-r--r-- 1 root root 636 Apr 21 23:07 scgi_params #AJAX前后分离
-rw-r--r-- 1 root root 664 Apr 21 23:07 uwsgi_params #Python
 
# 3.字符编码文件
-rw-r--r-- 1 root root 3610 Apr 21 23:07 win-utf
-rw-r--r-- 1 root root 2837 Apr 21 23:07 koi-utf
-rw-r--r-- 1 root root 2223 Apr 21 23:07 koi-win
 
# 4.浏览器支持的直接打开文件格式
-rw-r--r-- 1 root root 5231 Apr 21 23:07 mime.types
 
# 5.nginx相关命令文件
-rwxr-xr-x 1 root root 1342640 Apr 21 23:07 /usr/sbin/nginx
-rwxr-xr-x 1 root root 1461544 Apr 21 23:07 /usr/sbin/nginx-debug
 
# 6.日志相关文件
-rw-r--r-- 1 root root 351 Apr 21 23:05 /etc/logrotate.d/nginx
-rw-r----- 1 nginx adm 1654 May 14 11:35 access.log
-rw-r----- 1 nginx adm 4143 May 14 11:52 error.log
~~~

### nginx配置文件
~~~bash
# cat /etc/nginx/nginx.conf 
 
user  nginx;
worker_processes  1;
 
error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;
 
 
events {
    worker_connections  1024;
}
 
 
http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
 
    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';
 
    access_log  /var/log/nginx/access.log  main;
 
    sendfile        on;
    #tcp_nopush     on;
 
    keepalive_timeout  65;
 
    #gzip  on;
 
    include /etc/nginx/conf.d/*.conf;
}
{1}
{1}
##########################   模块划分 #########################
{1}
--------------------------------------   核心模块   --------------------------------------------
user  nginx;
worker_processes  1;
{1}
error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;
{1}
------------------------------------   事件驱动模块   -------------------------------------------
events {
    worker_connections  1024;
}
{1}
--------------------------------------   HTTP模块   ---------------------------------------------
http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
{1}
    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';
{1}
    access_log  /var/log/nginx/access.log  main;
{1}
    sendfile        on;
    #tcp_nopush     on;
{1}
    keepalive_timeout  65;
{1}
    #gzip  on;
{1}
    include /etc/nginx/conf.d/*.conf;
}
~~~
~~~bash

## 核心模块
 
# nginx启动用户
user  nginx;
# worker进程数
worker_processes  1;
# 错误日志的路径和级别
error_log  /var/log/nginx/error.log warn;
# pid文件的路径
pid        /var/run/nginx.pid;
 
## 事件驱动模块
 
events {
# 每一个worker进程允许连接数量
    worker_connections  1024;
}
 
## HTTP模块
 
http {
# 包含指定文件的内容，该文件是nginx浏览器允许访问的文件类型
include /etc/nginx/mime.types;
# 默认需要下载类型的格式
default_type application/octet-stream;
# 日志格式 （指定日志内容格式） 
 log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';
# 日志路径 和 指定格式 （指定日志名和路径）
access_log /var/log/nginx/access.log main;
access_log /var/log/nginx/zls_access.log zidingyi;
# 高效传输文件
sendfile on;
#tcp_nopush on;
# 长连接的超时时间
keepalive_timeout 65;
# 开启gzip压缩
#gzip on;
# 包含所有下面路径下conf结尾的文件（指定配置文件后缀）
include /etc/nginx/conf.d/*.conf;
}
~~~