# http 超文本传输协议( Hyper Text Transfer Protocol )   是Web基础协议  
 
学过的端口：
ftp:21
ssh:22
telnet:23
rsync:873
http:80
 
# URL 统一资源定位符（唯一标识）
  URL的组成：
# 例： https://www.google.com/doodles/about
协议：					http://
主机:端口			   www.google.com
文件名和路径			  服务器站点目录下的，目录和文件
 
# http工作原理
当在浏览器中输入 https://www.google.com/doodles/about
1.先分析url中的域名是谁   :  www.google.com
2.请求DNS服务器做解析	  ： 172.217.174.196
3.DNS把172.217.174.196 返回给浏览器
4.跟 172.217.174.196的80端口建立连接（建立TCP连接）
5.用GET请求下载/doodles/about
6.172.217.174.196 把页面返回给浏览器
7.断开TCP连接
8.浏览器显示页面
![访问网站分析](https://www.init0.cn/wp-content/uploads/2020/07/image-20200513192120355.png)
页面分析
--------------------     General （ 基本信息 ）     --------------------
# 请求网址
Request URL: https://www.google.com/doodles/about
# 请求方法
Request Method: GET
# 远程地址
Remote Address: 127.0.0.1:1080
# 状态码
Status Code: 200
 
-------------------- Request Headers （ 请求头部 ） --------------------
# 请求的资源类型
accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
# 资源类型压缩
accept-encoding: gzip, deflate, br
# 资源类型的语言
accept-language: zh-CN,zh;q=0.9
# 缓存控制：服务端没有开启缓存
cache-control: max-age=0
# 长连接
Connection: keep-alive
# 访问的主机：www.biadu.com
Host: www.google.com
#  项目缓存：没有开启
Pragma: no-cache
# 客户端优先加密
upgrade-insecure-requests: 1
# 用户访问网站的客户端工具
user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36
 
-------------------- Request Headers （ 响应头部 ） --------------------
# 建立长连接
Connection: keep-alive	
# 解析方式和字符集
Content-Type: text/html;charset=utf-8	     
# 日期
date: Wed, 13 May 2020 10:20:36 GMT	
# 该网站服务器，使用的软件和版本号
server: Google Frontend					       
# 状态码
status: 200	
HTTP请求方法
方法(Method)	含义
GET	请求读取一个Web页面（下载一个页面）
POST	附加一个命名资(如Web页面)（上传）
DELETE	删除Web页面
CONNECT	用于代理服务器
HEAD	请求读取一个Web页面的头部
PUT	请求存储一个Web页面
TRACE	用于测试，要求服务器送回收到的请求
OPTION	查询特定选项
HTTP响应方法
# 2xx和3xx都是网页可以正常访问
# 4xx：Nginx的报错（出错，出在nginx上）去检查nginx服务，或者服务器权限等...
# 5xx：后端报错（nginx后面连接的服务报错：mysql、php、tomcat、redis.......，排除前面的原因，最后就是代码有问题）
状态码	含义
200	成功
301	永久重定向（跳转）
302	临时重定向（跳转）
304	本地缓存（浏览器的缓存）
307	内部重定向（跳转）
400	客户端错误
401	认证失败
403	找不到主页，权限不足
404	找不到页面
405	请求方法不被允许
500	内部错误（MySQL关闭等...）
502	bad gateway 坏了的网关(php tomcat 等服务关闭)
503	服务端请求限制
504	请求超时
referer
# HTTP Referer是header的一部分，当浏览器向web服务器发送请求的时候，一般会带上Referer，告诉服务器该网页是从哪个页面链接过来的，服务器因此可以获得一些信息用于处理。
 
1.从主页上链接到一个朋友那里，他的服务器就能够从HTTP Referer中统计出每天有多少用户点击我主页上的链接访问他的网站。
2.Referer的正确英语拼法是referrer。由于早期HTTP规范的拼写错误，为了保持向后兼容就将错就错了。其它网络技术的规范企图修正此问题，使用正确拼法，所以拼法不统一。
3.Request.ServerVariables("HTTP_REFERER")的用法是防外连接。
详细的http原理
1.用输入域名 - > 浏览器跳转 - > 浏览器缓存 - > Hosts文件 - > DNS解析（递归查询|迭代查询）
    客户端向服务端发起查询 - > 递归查询
    服务端向服务端发起查询 - > 迭代查询
2.由浏览器向服务器发起TCP连接（三次握手）
    客户端     -->请求包连接  syn=1   seq=x             服务端
    服务端     -->响应客户端  syn=1   ack=x+1 seq=y     客户端
    客户端     -->建立连接    ack=y+1 seq=x+1           服务端
3.客户端发起http请求：
    1）请求的方法是什么:      GET获取
    2）请求的Host主机是:      www.google.com
    3）请求的资源是什么:      /doodles/about
    4）请求的端端口是什么:    默认http是80 https是443
    5）请求携带的参数是什么:  属性（请求类型、压缩、认证、浏览器信息、等等）
    6）请求最后的空行
4.服务端响应的内容是
    1）服务端响应使用WEB服务软件
    2）服务端响应请求文件类型
    3）服务端响应请求的文件是否进行压缩
    4）服务端响应请求的主机是否进行长连接
5.客户端向服务端发起TCP断开（四次挥手）
    客户端     --> 断开请求 fin=1 seq=x          -->    服务端
    服务端     --> 响应断开 fin=1 ack=x+1 seq=y  -->    客户端
    服务端     --> 断开连接 fin=1 ack=x+1 seq=z  -->    客户端
    客户端     --> 确认断开 fin=1 ack=x+1 seq=sj -->    服务端
HTTP相关术语
# pv
页面独立浏览量  一次访问可以产生多次请求
# pu
独立设备  开发统计
# IP
独立公网IP
# SOA松耦合架构
由多个网站模块化组成一个网站  