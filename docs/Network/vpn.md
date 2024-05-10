# vpn

## shadowsocks 服务端安装


1. 基础目录
```bash
mkdir -p /export/pids
mkdir -p /export/logs/php
mkdir -p /export/logs/nginx
mkdir -p /export/init.d
mkdir -p /export/shell
mkdir -p /export/tools
chmod 757 /export/pids
```


2. 安装Shadowsocket

```bash
yum install -y epel-release
yum install -y python-pip 

// 2.8.x 版本
pip install shadowsocks
或者
// 3.0.0 版本
pip install https://github.com/shadowsocks/shadowsocks/archive/master.zip -U
```

3. 配置shadowsocks.json文件

```bash
mkdir -p /export/shadowsocks/config
vim /export/shadowsocks/config/shadowsocks.json

{
    "server":"0.0.0.0",
    "port_password":{
        "8300":"NWY5MDk5N2NiY",
        "8301":"Y5MDk5N2NiY2U"
    },
    "local_address":"127.0.0.1",
    "local_port":1080,
    "timeout":300,
    "method":"aes-256-cfb",
    "fast_open":false
}
```

4. 配置启动脚本

```bash
vim /export/init.d/shadowsocks

#! /bin/bash
// $* set [start/stop/restart]
ssserver -c /export/shadowsocks/config/shadowsocks.json --log-file /export/logs/shadowsocks.log --pid-file=/export/pids/shadowsocks.pid -d $*

```

5. 启动服务

```bash
chmod 757 /export/init.d/shadowsocks
/export/init.d/shadowsocks start
```


## shadowsocks客户端安装

1. 基础目录
```bash
mkdir -p /export/pids
mkdir -p /export/logs/php
mkdir -p /export/logs/nginx
mkdir -p /export/init.d
mkdir -p /export/shell
mkdir -p /export/tools
chmod 757 /export/pids
```


2. 安装Shadowsocket

```bash
yum install -y epel-release
yum install -y python-pip 

// 2.8.x 版本
pip install shadowsocks
或者
// 3.0.0 版本
pip install https://github.com/shadowsocks/shadowsocks/archive/master.zip -U
```

3. 配置shadowsocks.json文件

```bash
mkdir -p /export/shadowsocks/config
vim /export/shadowsocks/config/shadowsocks.json

{
  "server":"47.111.177.46", # 服务端地址
  "server_port":38883, # 服务端端口
  "password":"NWY5MDk5N2NiY", # 服务端密码
  "local_address": "0.0.0.0", # 本地环境服务地址
  "local_port":8300, # 本地环境服务端口
  "timeout":300,
  "method":"aes-256-cfb",
  "workers": 1
}
```

4. 配置启动脚本

```bash
vim /export/init.d/shadowsocks

#! /bin/bash
// $* set [start/stop/restart]
sslocal -c /export/shadowsocks/config/shadowsocks.json --log-file /export/logs/shadowsocks.log --pid-file=/export/pids/shadowsocks.pid -d $*
```

5. 启动服务

```bash
chmod 757 /export/init.d/shadowsocks
/export/init.d/shadowsocks start
```

## privoxy代理

1. 安装privoxy

```bash
yum install -y epel-release
yum install privoxy
```

2. 配置privoxy

```bash
# 修改下列项
vim /etc/privoxy/config

listen-address  0.0.0.0:8118
toggle  1
forward-socks5t   unionpay-api.pinpula.com         127.0.0.1:8300 .
forward         192.168.*.*/     .
forward          10.*.*.*/       .
forwarded-connect-retries  1
```

3. 启动

```bash
systemctl enable --now privoxy
```

4. 代理访问

```bash
export http_proxy=http://192.168.108.130:8118 https_proxy=http://192.168.108.130:8118
```

- 相关资料
  - https://github.com/shadowsocks
  - https://brickyang.github.io/2017/01/14/CentOS-7-%E5%AE%89%E8%A3%85-Shadowsocks-%E5%AE%A2%E6%88%B7%E7%AB%AF/

- vpn入门系列
  - https://www.codeleading.com/article/43202129261/
  - https://www.shuzhiduo.com/A/kmzLMoLAJG/
  - https://www.dddisk.com/post/117.html?ivk_sa=1024320u
  - https://mp.weixin.qq.com/s/ElNS6Kw7SN81kZVzuVjQww

- privoxy
  - https://blog.csdn.net/yihuliunian/article/details/105340812?spm=1001.2101.3001.6650.6&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7ERate-6-105340812-blog-132575873.235%5Ev38%5Epc_relevant_anti_vip_base&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7ERate-6-105340812-blog-132575873.235%5Ev38%5Epc_relevant_anti_vip_base&utm_relevant_index=7