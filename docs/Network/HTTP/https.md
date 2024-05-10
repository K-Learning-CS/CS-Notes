端口回顾
ssh:22
telnet:23
ftp:21
dns:53
rsync:873
http:80
mysql:3306
redis:6379
php:9000
tomcat:8080
---------------
https:443
### HTTPS

# HTTP:超文本传输协议（明文，不安全）
# HTTPS：加密超文本传输协议（安全）
 
1.使用https的原因：
    因为HTTP不安全，当我们使用http网站时，会遭到劫持和篡改，如果采用https协议，那么数据在传输过程中是加密的，所以黑客无法窃取或者篡改数据报文信息，同时也避免网站传输时信息泄露。
 
2.加密方法：
    实现https时，需要了解ssl协议，但我们现在使用的更多的是TLS加密协议，在OSI七层模型中，http协议在最顶部的应用层，在经过下一级的表示层时，ssl协议发挥作用，他通过（握手、交换秘钥、告警、加密）等方式，在应用层http协议没有感知的情况下做到了数据的安全加密。
 
3.验证机制：
    在数据进行加密与解密过程中，需要有一个权威机构CA来验证双方身份。我们首先需要申请证书，先去登记机构进行身份登记，我是谁，我是干嘛的，我想做什么，然后登记机构再通过CSR发给CA，CA中心通过后会生成一堆公钥和私钥，公钥会在CA证书链中保存，公钥和私钥证书我们拿到后，会将其部署在WEB服务器上。
 
4.验证流程：
    1）当浏览器访问我们的https站点时，他会去请求我们的证书
    2）Nginx这样的web服务器会将我们的公钥证书发给浏览器
    3）浏览器会去验证我们的证书是否合法有效
    4）.CA机构会将过期的证书放置在CRL服务器，CRL服务的验证效率是非常差的，所以CA有推出了OCSP响应程序，OCSP响应程序可以查询指定的一个证书是否过去，所以浏览器可以直接查询OSCP响应程序，但OSCP响应程序性能还不是很高
    5）Nginx会有一个OCSP的开关，当我们开启后，Nginx会主动上OCSP上查询，这样大量的客户端直接从Nginx获取证书是否有效
    证书类型
对比	域名型 DV	企业型 OV	增强型 EV
绿色地址栏	img小锁标记+https	img小锁标记+https	img小锁标记+企业名称+https
一般用途	个人站点和应用； 简单的https加密需求	电子商务站点和应用； 中小型企业站点	大型金融平台； 大型企业和政府机构站点
审核内容	域名所有权验证	全面的企业身份验证； 域名所有权验证	最高等级的企业身份验证； 域名所有权验证
颁发时长	10分钟-24小时	3-5个工作日	5-7个工作日
单次申请年限	1年	1-2年	1-2年
赔付保障金	——	125-175万美金	150-175万美金
# 证书选择
1.单域名  不要钱 一年一续
2.多域名  稍微便宜点
3.通配符  贼贵
# 注意事项
1.https不支持续费，证书到期需要重新申请并进行替换
2.https不支持三级域名，例如：test.m.driverzeng.com
# https状态
1.绿色状态：网站是安全的
2.黄色状态：代码中，带有http的不安链接
3.红色状态：网站内部有其它不安全连接
    假装配置ssl证书
#nginx必须有ssl模块
[root@web03 ~]# nginx -V
--with-http_ssl_module
 
 
## 创建证书存放路径
[root@web01 nginx]# mkdir /etc/nginx/ssl
Generating RSA private key, 2048 bit long modulus
..+++
...............................+++
e is 65537 (0x10001)
 
### 密码是1234
Enter pass phrase for /etc/nginx/ssl/20200603105245_www.linux.com.key: 1234
Verifying - Enter pass phrase for /etc/nginx/ssl/20200603105245_www.linux.com.key: 1234
{1}
## CA颁发证书
[root@web01 ssl]# openssl genrsa -idea -out /etc/nginx/ssl/$(date +%Y%m%d%H%M%S)_www.linux.com.key 2048
{1}
## 自签证书
[root@web01 ssl]# openssl req -days 36500 -x509 -sha256 -nodes -newkey rsa:2048 -keyout /etc/nginx/ssl/20200603105245_www.linux.com.key -out /etc/nginx/ssl/20200603105245_www.linux.com.crt
Generating a 2048 bit RSA private key
......+++
.................................................+++
writing new private key to '/etc/nginx/ssl/20200603105245_www.linux.com.key'
-----
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:China
string is too long, it needs to be less than  2 bytes long
Country Name (2 letter code) [XX]:CN
State or Province Name (full name) []:Shanghai
Locality Name (eg, city) [Default City]:Shanghai
Organization Name (eg, company) [Default Company Ltd]:oldboy
Organizational Unit Name (eg, section) []:ops
Common Name (eg, your name or your server's hostname) []:linux.com                                              
Email Address []:123@qq.com
{1}
{1}
# 为了用户体验，必须80强转443
server{
        listen 80;
        server_name www.linux.com;
        #return 302 https://$server_name$request_uri;
        rewrite (.*) https://$server_name$request_uri redirect;
}
HTTPS参数优化
upstream web {
server 172.16.1.7;
server 172.16.1.8;
}
server {
listen 80;
server_name blog.linux.com;
return 302 https://$server_name$request_uri;
}
server {
listen 443 ssl;
server_name blog.linux.com;
ssl_certificate /etc/nginx/ssl/20200603105245_www.linux.com.crt;
ssl_certificate_key /etc/nginx/ssl/20200603105245_www.linux.com.key;
# 下面这一段
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 1440m;
ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
ssl_prefer_server_ciphers on;
 
location / {
proxy_pass http://web;
proxy_set_header Host $host;
}
}
测试
主机名	Wan IP	Lan IP	搭建服务
lb01	10.0.0.5	172.16.1.5	负载均衡
web01	10.0.0.7	172.16.1.7	nginx和php
web02	10.0.0.8	172.16.1.8	nginx和php
web03	10.0.0.9	172.16.1.9	nginx和php
nfs	10.0.0.31	172.16.1.31	nfs和sersync
backup	10.0.0.41	172.16.1.41	rsync
db01	10.0.0.51	172.16.1.51	MySQL
backup
1.写脚本
vim /root/rsync.sh
 
#!/bin/bash
install=`yum install -y rsync`
cat >/etc/rsyncd.conf<<'EOF'
#!/bin/bash
uid = rsync
gid = rsync
port = 873
fake super = yes
use chroot = no
max connections = 200
timeout = 600
ignore errors
read only = false
list = false
auth users = kang_bak
secrets file = /etc/rsync.passwd
log file = /var/log/rsyncd.log
 
[backup]
comment = welcome to oldboyedu backup!
path = /backup
EOF
 
useradd rsync -s /sbin/nologin -M
mkdir /backup
chown rsync.rsync /backup/ -R
echo 'kang_bak:123' > /etc/rsync.passwd
chmod 600  /etc/rsync.passwd
systemctl start rsyncd
systemctl enable rsyncd
 
2.一键部署
sh /root/rsync.sh
nfs
1.写脚本
vim /root/sersync.sh
 
#!/bin/bash
 
install=`yum install -y rsync nfs-utils inotify-tools`
echo "/code/wp 172.16.1.0/24(rw,sync,all_squash,anonuid=666,anongid=666)" >> /etc/exports
echo "/code/zh 172.16.1.0/24(rw,sync,all_squash,anonuid=666,anongid=666)" >> /etc/exports
groupadd www -g 666
useradd www -u 666 -g 666 -s /sbin/nologin -M
mkdir -p /code/{wp,zh}
chown www.www /code/
systemctl start rpcbind nfs-server
systemctl enable rpcbind nfs-server
 
download=`wget http://test.driverzeng.com/other/sersync2.5.4_64bit_binary_stable_final.tar.gz`
tar xf sersync2.5.4_64bit_binary_stable_final.tar.gz
mv GNU-Linux-x86 /usr/local/sersync
 
cat >/usr/local/sersync/confxml.xml<<'EOF'
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
    <localpath watch="/code">
 
        <!-- rsync服务端的IP 和 name:模块 -->
        <remote ip="10.0.0.41" name="backup"/>
        <!--<remote ip="192.168.8.39" name="tongbu"/>-->
        <!--<remote ip="192.168.8.40" name="tongbu"/>-->
    </localpath>
    <rsync>
        <!-- rsync命令执行的参数 -->
        <commonParams params="-az"/>
            <!-- rsync认证start="true" users="rsync指定的匿名用户" passwordfile="指定一个密码文件的位置权限必须600" -->
        <auth start="true" users="kang_bak" passwordfile="/etc/rsync.passwd"/>
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
    <param prefix="/bin/sh" suffix="" ignoreError="true"/>  <!--prefix /opt/tongbu/mmm.sh suffix-->
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
EOF
 
echo '123' > /etc/rsync.passwd
chmod 600  /etc/rsync.passwd
/usr/local/sersync/sersync2 -rdo /usr/local/sersync/confxml.xml
 
2.一键部署
sh /root/sersync.sh
db01
1.安装MySQL
yum install -y mariadb-server
2.启动服务，并加入开机自启
systemctl start mariadb && systemctl enable mariadb
3.给root用户密码
mysqladmin -uroot password '123'
4.连接数据库
mysql -uroot -p123
5.创建数据库
create database wp;
create database zh;
6.查看是否创建成功
show databases;
7.创建WordPress连接数据库的用户和密码
grant all on *.* to php_user@'%' identified by '111';
web01
1.更换官方源
cat>>/etc/yum.repos.d/nginx.repo<<'EOF'
[nginx-stable]
name=nginx stable repo
baseurl=http://nginx.org/packages/centos/$releasever/$basearch/
gpgcheck=1
enabled=1
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true
EOF
2.安装nginx
yum install -y nginx
3.新建并修改nginx用户
groupadd -g 666 www
useradd -u 666 -g 666 -s /sbin/nologin -M www
sed -i 's#^user  nginx#user  www#' /etc/nginx/nginx.conf
4.加入开机自启
systemctl enable nginx
5.添加nginx配置文件
vim /etc/nginx/conf.d/wp.conf
server {
        listen 80;
        server_name wp.kang.com;
        root /code/wp;
        index index.php index.html;
 
        location ~ \.php$ {
                root /code/wp;
         
                fastcgi_pass localhost:9000;
                fastcgi_index index.php;
                fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
                include        fastcgi_params;
        }
}
 
vim /etc/nginx/conf.d/zh.conf
server {
        listen 80;
        server_name zh.kang.com;
        root /code/zh;
        index index.php index.html;
 
        location ~ \.php$ {
                root /code/zh;
         
                fastcgi_pass localhost:9000;
                fastcgi_index index.php;
                fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
                include        fastcgi_params;
        }
}
 
# 让所有的访问都走https
echo 'fastcgi_param HTTPS on;' >> /etc/nginx/fastcgi_params
 
6.创建对应文件
mkdir -p /code/{wp,zh}
7.将对应文件解压并放入对应文件夹
8.授权
chown -R www.www /code/
9.创建并挂载图片目录
mkdir -p /code/wordpress/wp-content/uploads/
mount -t nfs 172.16.1.31:/code/wp /code/wordpress/wp-content/uploads/
mkdir -p /code/zh/uploads
mount -t nfs 172.16.1.31:/code/zh /code/zh/uploads
10.安装php，先卸载自带
yum remove php-mysql-5.4 php php-fpm php-common
11.更换php源
vim /etc/yum.repos.d/php.repo
[php-webtatic]
name = PHP Repository
baseurl = http://us-east.repo.webtatic.com/yum/el7/x86_64/
gpgcheck = 0
enabled = 1
12.安装php
yum -y install php71w php71w-cli php71w-common php71w-devel php71w-embedded php71w-gd php71w-mcrypt php71w-mbstring php71w-pdo php71w-xml php71w-fpm php71w-mysqlnd php71w-opcache php71w-pecl-memcached php71w-pecl-redis php71w-pecl-mongodb
13.更改php用户和用户组
sed -i 's#^user = apache#user = www#' /etc/php-fpm.d/www.conf
sed -i 's#^group = apache#group = www#' /etc/php-fpm.d/www.conf
14.启动php并加入开机自启
systemctl start php-fpm && systemctl enable php-fpm
15.启动nginx
systemctl start nginx
16.在windows的 hosts文件中加入域名解析
17.浏览器打开wp.com
数据库名	  wp
用户名		  php_user
密码		   111
数据库主机	10.0.0.51
表前缀	 	  wp_
18.浏览器打开zh.com
数据库名称     zh
数据库用户名   php_user
数据库密码    111
数据库地址    10.0.0.51
表前缀        zh_  
 
19.wp后台把网站改https
web02
1.更换官方源
vim /etc/yum.repos.d/nginx.repo
[nginx-stable]
name=nginx stable repo
baseurl=http://nginx.org/packages/centos/$releasever/$basearch/
gpgcheck=1
enabled=1
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true
2.安装nginx
yum install -y nginx
3.新建并修改nginx用户
groupadd -g 666 www
useradd -u 666 -g 666 -s /sbin/nologin -M www
sed -i 's#^user  nginx#user  www#' /etc/nginx/nginx.conf
4.从web01添加
rm -rf /etc/php-fpm.d/
rm -rf /etc/php.ini 
rsync -avz /etc/php-fpm.d/   root@10.0.0.8:/etc/php-fpm.d/
rsync -avz /etc/php.ini   root@10.0.0.8:/etc/php.ini
rsync -avz /code/   root@10.0.0.8:/code
rsync -avz /etc/nginx/conf.d/   root@10.0.0.8:/etc/nginx/conf.d/
5.挂载图片目录
mount -t nfs 172.16.1.31:/code/wp /code/wordpress/wp-content/uploads/
mount -t nfs 172.16.1.31:/code/zh /code/zh/uploads
6.安装php，先卸载自带
yum remove php-mysql-5.4 php php-fpm php-common
7.更换php源
vim /etc/yum.repos.d/php.repo
[php-webtatic]
name = PHP Repository
baseurl = http://us-east.repo.webtatic.com/yum/el7/x86_64/
gpgcheck = 0
enabled = 1
8.安装php
yum -y install php71w php71w-cli php71w-common php71w-devel php71w-embedded php71w-gd php71w-mcrypt php71w-mbstring php71w-pdo php71w-xml php71w-fpm php71w-mysqlnd php71w-opcache php71w-pecl-memcached php71w-pecl-redis php71w-pecl-mongodb
9.更改php用户和用户组
sed -i 's#^user = apache#user = www#' /etc/php-fpm.d/www.conf
sed -i 's#^group = apache#group = www#' /etc/php-fpm.d/www.conf
10.启动php并加入开机自启
systemctl start php-fpm && systemctl enable php-fpm
11.启动nginx并加入开机自启
systemctl start nginx && systemctl enable nginx
web03
1.更换官方源
vim /etc/yum.repos.d/nginx.repo
[nginx-stable]
name=nginx stable repo
baseurl=http://nginx.org/packages/centos/$releasever/$basearch/
gpgcheck=1
enabled=1
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true
2.安装nginx
yum install -y nginx
3.新建并修改nginx用户
groupadd -g 666 www
useradd -u 666 -g 666 -s /sbin/nologin -M www
sed -i 's#^user  nginx#user  www#' /etc/nginx/nginx.conf
4.从web01添加
rm -rf /etc/php-fpm.d/
rm -rf /etc/php.ini
rsync -avz /etc/php-fpm.d/   root@10.0.0.9:/etc/php-fpm.d/
rsync -avz /etc/php.ini   root@10.0.0.9:/etc/php.ini
rsync -avz /code/   root@10.0.0.9:/code
rsync -avz /etc/nginx/conf.d/   root@10.0.0.9:/etc/nginx/conf.d/
5.挂载图片目录
mount -t nfs 172.16.1.31:/code/wp /code/wordpress/wp-content/uploads/
mount -t nfs 172.16.1.31:/code/zh /code/zh/uploads
6.安装php，先卸载自带
yum remove php-mysql-5.4 php php-fpm php-common
7.更换php源
vim /etc/yum.repos.d/php.repo
[php-webtatic]
name = PHP Repository
baseurl = http://us-east.repo.webtatic.com/yum/el7/x86_64/
gpgcheck = 0
enabled = 1
8.安装php
yum -y install php71w php71w-cli php71w-common php71w-devel php71w-embedded php71w-gd php71w-mcrypt php71w-mbstring php71w-pdo php71w-xml php71w-fpm php71w-mysqlnd php71w-opcache php71w-pecl-memcached php71w-pecl-redis php71w-pecl-mongodb
9.更改php用户和用户组
sed -i 's#^user = apache#user = www#' /etc/php-fpm.d/www.conf
sed -i 's#^group = apache#group = www#' /etc/php-fpm.d/www.conf
10.启动php并加入开机自启
systemctl start php-fpm && systemctl enable php-fpm
11.启动nginx并加入开机自启
systemctl start nginx && systemctl enable nginx
lb01
1.更换官方源
cat>>/etc/yum.repos.d/nginx.repo<<'EOF'
[nginx-stable]
name=nginx stable repo
baseurl=http://nginx.org/packages/centos/$releasever/$basearch/
gpgcheck=1
enabled=1
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true
EOF
2.安装nginx
yum install -y nginx
3.新建并修改nginx用户
groupadd -g 666 www
useradd -u 666 -g 666 -s /sbin/nologin -M www
sed -i 's#^user  nginx#user  www#' /etc/nginx/nginx.conf
4.配置证书
mkdir /etc/nginx/ssl
cd /etc/nginx/ssl
openssl req -days 36500 -x509 \
> -sha256 -nodes -newkey rsa:2048 -keyout server.key -out server.crt
Generating a 2048 bit RSA private key
 
添加nginx配置文件
vim /etc/nginx/conf.d/zh.conf
upstream zh {
        server 172.16.1.7;
        server 172.16.1.8;
        server 172.16.1.9;
}
 
server {
        listen 80;
        server_name zh.kang.com;
        return 302 https://$server_name$request_uri;
}
 
server {
        listen 443 ssl;
        server_name zh.kang.com;
        ssl_certificate     /etc/nginx/ssl/server.crt;
        ssl_certificate_key /etc/nginx/ssl/server.key;
 
        location / {
                proxy_pass http://zh;
                proxy_set_header Host $host;
        }
}
 
vim /etc/nginx/conf.d/wp.conf
upstream wp {
        server 172.16.1.7;
        server 172.16.1.8;
        server 172.16.1.9;
}
 
server {
        listen 80;
        server_name wp.kang.com;
        return 302 https://$server_name$request_uri;
}
 
server {
        listen 443 ssl;
        server_name wp.kang.com;
        ssl_certificate     /etc/nginx/ssl/server.crt;
        ssl_certificate_key /etc/nginx/ssl/server.key;
 
        location / {
                proxy_pass http://wp;
                proxy_set_header Host $host;
        }
}
 
5.启动并加入开机自启
systemctl start nginx && systemctl enable nginx
7.物理机hosts解析
10.0.0.5    zh.kang.com wp.kang.com
8.浏览器访问
zh.kang.com