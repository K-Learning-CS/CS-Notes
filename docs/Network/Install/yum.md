rpm
rpm包属性
 
[root@oldboy /mnt/Packages]# rpm -q bash
bash-4.2.46-33.el7.x86_64.rpm
软件名-发行版本-当前版本发布次数-适用的linux版本-硬件平台  扩展名

#分类：
rpm包      预先编译打包 redhat的规范    软件版本低
源码包      手动编译打包 软件官网下载    版本随意
二进制包    解压即可使用                 不能修改源码
 
#获取途径：
1.各官网
2.光盘镜像
3.rpm包查询网站
4.各大国内源镜像站
rpm + 选项 + 参数
#选项：
-i 安装
-v 显示详细信息
-h 显示进度
-f 仅升级
-U 没有则安装 有责升级
-e 删除
-q 查询
-q --scripts  查询rpm包安装前后执行的脚本
-qa 查询所有已安装的rpm包
-qd 查询rpm包的man帮助
-qc 查询当前rpm包的所有配置文件
-ql 查询当前rpm包的所有目录及文件
-qi 查询当前rpm包的详细信息
-qf 查询属于当前rpm包的命令
-qip 查询查询当前未安装rpm包的详细信息
-qlp 查询当前未安装rpm包的所有目录及文件
--test 安装前测试
--force 强制安装
--nodeps 忽略依赖关系
 
#可以将链接作为参数使用
yum
#特性：
1.必须联网
2.管理rpm包
3.自动解决依赖关系
4.命令简单 
5.生产环境常用
 
#yum命令：
yum install    安装
yum reinstall  重装
yum remove     卸载
yum list       列出yum库中所有可用rpm包
yum repolist   列出yum库及包的数量
yum provides   查询当前命令属于哪个rpm包
yum check-update  检查更新
yum update        接包名时更新软件  不接时全部更新
yum clean all     清楚缓存
yum makecache     建立缓存
yum history       yum历史记录
yum group         yum组管理
本地yum库
1.下载nginx相关的rpm包
mkdir /nginx/base
yum install nginx --downloadonly --downloaddir=/nginx/base
2.创建仓库
yum install -y createrepo
createrepo /nginx/base
3.新建配置文件
mv /etc/yum.repos.d/* ./
vim /etc/yum.repos.d/nginx.repo
[nginx]
name=nginx_repo
baseurl=file:///nginx/base
gpgcheck=0
enabled=1
4.测试
yum clean all
yum repolist
yum install -y nginx
ftp yum仓库
#服务器：
1.安装vsftpd
rm -rf  /etc/yum.repos.d/*
mv *.repo /etc/yum.repos.d/
yum install -y vsftpd
2.创建仓库
mkdir /var/ftp/pub/nginx
mv /nginx/base/*.rpm /var/ftp/pub/nginx
createrepo /var/ftp/pub/nginx
#客户机：
3.创建yum配置文件
mv /etc/yum.repos.d/* /opt/
vim /etc/yum.repos.d/nginx.repo
[nginx]
name=nginx_repo
baseurl=ftp://10.0.0.5/pub/nginx/
gpgcheck=0
enabled=1
4.测试
yum clean all
yum repolist
yum install -y nginx
nginx yum仓库
#服务器：
1.安装nginx
yum install -y nginx
2.添加配置文件
vim /etc/nginx/conf.d/yum.conf
server {
	listen 80;
	server_name www.kang.com;
	location / {
	root	/nginx;
	autoindex on;
	access_log off;
	}
}
3.创建仓库
rm -rf /nginx/base/*
mv /var/ftp/pub/nginx/*.rpm /nginx/base/
createrepo /nginx/base/
4.启动nginx
systemctl start nginx
#客户机：
5.配置域名解析
echo "10.0.0.5 www.kang.com" >> /etc/hosts
6.新建yum仓库配置文件
rm -rf /etc/yum.repos.d/*
vim /etc/yum.repos.d/nginx.repo
[nginx]
name=nginx_repo
baseurl=http://www.kang.com/base/
gpgcheck=0
enbaled=1
7.测试
yum clean all
yum repolist
yum install -y nginx
源码安装
1.下载源码包
wget http://nginx.org/download/nginx-1.16.1.tar.gz
2.安装依赖
yum install -y gcc gcc-c++ glibc zlib-devel pcre-devel openssl-devel
3.解压并进入
tar xf nginx-1.16.1.tar.gz 
cd nginx-1.16.1
4.创建用户并指定安装位置
useradd nginx -s /sbin/nologin -M
mkdir /app
./configure --prefix=/app/nginx-1.16.1 --user=nginx --group=nginx
5.编译安装
make && make install
6.启动nginx并检测80端口
/app/nginx-1.16.1/sbin/nginx
netstat -lntup | grep 80
7.更新软连接
ln -s /app/nginx-1.16.1 /app/nginx
8.修改nginx默认页面
vim /app/nginx-1.16.1/html/index.html
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8"/>
<title>欢迎来到DSB的nginx页面</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>欢迎来到DSB的nginx的web页面</h1>
<p><em>Thank you for using nginx.</em></p>
</body>
</html>
自制rpm包
1.下载解压
mkdir fpm
cd /root/fpm/
wget http://test.driverzeng.com/other/fpm-1.3.3.x86_64.tar.gz
tar xf fpm-1.3.3.x86_64.tar.gz 
2.安装ruby
yum -y install ruby rubygems ruby-devel rpm-build
3.查看并更换gem源
gem sources --list
gem sources --remove https://rubygems.org/
gem sources -a https://mirrors.aliyun.com/rubygems/
4.安装fpm
gem install *.gem
5.写出rpm包的执行脚本
vim nginx.sh
#!/bin/bash
useradd nginx -s /sbin/nologin -M
ln -s /app/nginx-1.16.1 /app/nginx
6.使用fpm制作rpm包
fpm -s dir -t rpm -n nginx -v 1.16.1 -d 'zlib-devel,gcc,gcc-c++,glibc,pcre-devel,openssl-devel' --post-install /root/nginx.sh -f /app/nginx-1.16.1/