# 简介

### 1.Tailscale

*Tailscale 是一种基于 WireGuard 的虚拟组网工具，和 Netmaker 类似，
最大的区别在于 Tailscale 是在用户态实现了 WireGuard 协议，
而 Netmaker 直接使用了内核态的 WireGuard。
所以 Tailscale 相比于内核态 WireGuard 性能会有所损失，
但与 OpenVPN 之流相比还是能甩好几十条街的，
Tailscale 虽然在性能上做了些许取舍，
但在功能和易用性上绝对是完爆其他工具，
它的主要优势有：*

1. 开箱即用
    - 无需配置防火墙
    - 没有额外的配置

2. 高安全性/私密性
    - 自动密钥轮换 
    - 点对点连接 
    - 支持用户审查端到端的访问记录

3. 在原有的 ICE、STUN 等 UDP 协议外，实现了 DERP TCP 协议来实现 NAT 穿透
4. 基于公网的控制服务器下发 ACL 和配置，实现节点动态更新
5. 通过第三方（如 Google） SSO 服务生成用户和私钥，实现身份认证

简而言之，我们可以将 Tailscale 看成是更为易用、功能更完善的 WireGuard。 

*Tailscale 是一款商业产品，但个人用户是可以白嫖的，
个人用户在接入设备不超过 20 台的情况下是可以免费使用的
（虽然有一些限制，比如子网网段无法自定义，且无法设置多个子网）。
除 Windows 和 macOS 的图形应用程序外，其他 Tailscale 客户端的组件
（包含 Android 客户端）是在 BSD 许可下以开源项目的形式开发的，
你可以在他们的 GitHub 仓库找到各个操作系统的客户端源码。*

### 2.Headscale

*Headscale 由欧洲航天局的 Juan Font 使用 Go 语言开发，在 BSD 许可下发布，
实现了 Tailscale 控制服务器的所有主要功能，可以部署在企业内部，
没有任何设备数量的限制，且所有的网络流量都由自己控制。*


### 3.DERP

*Tailscale 的终极目标是让两台处于网络上的任何位置的机器建立点对点连接（直连），
但现实世界是复杂的，大部份情况下机器都位于 NAT 和防火墙后面，这时候就需要通过打洞来实现直连，
也就是 NAT 穿透。*

NAT 按照 NAT 映射行为和有状态防火墙行为可以分为多种类型，
但对于 NAT 穿透来说根本不需要关心这么多类型，
只需要看 NAT 或者有状态防火墙是否会严格检查目标 Endpoint，
根据这个因素，可以将 NAT 分为 Easy NAT 和 Hard NAT。

1. Easy NAT 及其变种称为 “Endpoint-Independent Mapping” (EIM，终点无关的映射) 
   这里的 Endpoint 指的是目标 Endpoint，也就是说，有状态防火墙只要看到有客户端自己发起的出向包，
   就会允许相应的入向包进入，不管这个入向包是谁发进来的都可以。
2. hard NAT 以及变种称为 “Endpoint-Dependent Mapping”（EDM，终点相关的映射）
   这种 NAT 会针对每个目标 Endpoint 来生成一条相应的映射关系。 在这样的设备上，
   如果客户端向某个目标 Endpoint 发起了出向包，假设客户端的公网 IP 是 2.2.2.2，
   那么有状态防火墙就会打开一个端口，假设是 4242。
   那么只有来自该目标 Endpoint 的入向包才允许通过 2.2.2.2:4242，其他客户端一律不允许。
   这种 NAT 更加严格，所以叫 Hard NAT。

对于 Easy NAT，我们只需要提供一个第三方的服务，
它能够告诉客户端“它看到的客户端的公网 ip:port 是什么”，
然后将这个信息以某种方式告诉通信对端（peer）， 后者就知道该和哪个地址建连了！
这种服务就叫 STUN (Session Traversal Utilities for NAT，NAT会话穿越应用程序)。

对于 Hard NAT 来说，STUN 就不好使了，
即使 STUN 拿到了客户端的公网 ip:port 告诉通信对端也于事无补，
因为防火墙是和 STUN 通信才打开的缺口，这个缺口只允许 STUN 的入向包进入，
其他通信对端知道了这个缺口也进不来。通常企业级 NAT 都属于 Hard NAT。

*这种情况下打洞是不可能了，但也不能就此放弃，可以选择一种折衷的方式：
创建一个中继服务器（relay server），客户端与中继服务器进行通信，
中继服务器再将包中继（relay）给通信对端。*

*DERP 即 Detoured Encrypted Routing Protocol，这是 Tailscale 自研的一个协议：
    它是一个通用目的包中继协议，运行在 HTTP 之上，而大部分网络都是允许 HTTP 通信的。
    它根据目的公钥（destination’s public key）来中继加密的流量（encrypted payloads）。*


Tailscale 使用的算法很有趣，所有客户端之间的连接都是先选择 DERP 模式（中继模式），
这意味着连接立即就能建立（优先级最低但 100% 能成功的模式），用户不用任何等待。
然后开始并行地进行路径发现，通常几秒钟之后，我们就能发现一条更优路径，
然后将现有连接透明升级（upgrade）过去，变成点对点连接（直连）。

因此，DERP 既是 Tailscale 在 NAT 穿透失败时的保底通信方式（此时的角色与 TURN 类似），
也是在其他一些场景下帮助我们完成 NAT 穿透的旁路信道。 换句话说，它既是我们的保底方式，
也是有更好的穿透链路时，帮助我们进行连接升级（upgrade to a peer-to-peer connection）的基础设施。



# 搭建

## 一、前置需求

1. 两个域名：例如 headscale.xxx.com derp.xxx.com 名称可以自定义，这里仅说明用途
2. 域名证书：nginx版本的
3. 拥有固定外网IP的服务器：用于部署 headscale server与 derp server

## 二、服务端安装

*提前安装好 docker 和 docker-compose*

```bash
1.创建数据目录
mkdir /data/ && cd /data/

2.创建 docker-compose 脚本
cat >> docker-compose.yaml <<'EOF'
version: '3.5'
services:
  headscale:
    container_name: headscale
    restart: always
    image: headscale/headscale:latest
    ports:
      - 8080:8080
    networks:
      - headscale
    volumes:
      - /usr/share/zoneinfo/Asia/Shanghai:/etc/localtime:ro
      - /data/config:/etc/headscale
      - /data/data:/var/lib/headscale
      - /data/cert:/var/lib/headscale/cert
    command: headscale serve
    depends_on:
      - derp
  derp:
    container_name: derp
    restart: always
    image: docker.io/fredliang/derper:latest
    ports:
      - 8443:8443
      - 3478:3478
    networks:
      - headscale
    environment:
      - DERP_ADDR=:8443
      - DERP_DOMAIN=derp.rsjitcm.com
      - DERP_CERT_MODE=manual
      - DERP_VERIFY_CLIENTS=true
    volumes:
      - /usr/share/zoneinfo/Asia/Shanghai:/etc/localtime:ro
      - /data/cert:/app/certs
      - /data/tailscale-sock:/var/run/tailscale
    depends_on:
      - tailscale
  tailscale:
    container_name: tailscale
    restart: always
    image: tailscale/tailscale:v1.52.0
    network_mode: host
    privileged: true
    cap_add:
        - NET_ADMIN
        - NET_RAW
    volumes:
      - /usr/share/zoneinfo/Asia/Shanghai:/etc/localtime:ro
      - /data/lib:/var/lib
      - /dev/net/tun:/dev/net/tun
      - /data/tailscale-sock:/var/run/tailscale
    command: sh -c "mkdir -p /var/run/tailscale && ln -s /tmp/tailscaled.sock /var/run/tailscale/tailscaled.sock && tailscaled"
networks:
  headscale:
    external: true
EOF

3.创建docker网桥
docker network create headscale

4.创建证书目录
mkdir /data/cert

5.将两个域名的证书放入该目录
！注意：derp的证书名称必须为 derp.xxx.com.crt derp.xxx.com.key

6.创建config目录，并放入配置
mkdir /data/config

cat >> /data/config/derp.yaml<< 'EOF'
regions:
  900:
    regionid: 900
    regioncode: aliyun # 可用区代码
    regionname: aliyund-derp # 可用区名称
    nodes:
      - name: 900a # 节点名称
        regionid: 900
        hostname: derp.xxx.com # 节点域名
        stunport: 3478 # udp端口
        stunonly: false
        derpport: 8443 # https端口
EOF

cat >> /data/config/config.yaml<< 'EOF'
server_url: https://headscale.xxx.com:8080 # 对外暴露的连接
listen_addr: 0.0.0.0:8080 # 服务器监听
metrics_listen_addr: 127.0.0.1:9090
grpc_listen_addr: 127.0.0.1:50443
grpc_allow_insecure: false
private_key_path: /var/lib/headscale/private.key
noise:
  private_key_path: /var/lib/headscale/noise_private.key
ip_prefixes: # 固定配置不能修改，否则会报错
  - fd7a:115c:a1e0::/48
  - 100.64.0.0/10
derp:
  server:
    enabled: false
    region_id: 999
    region_code: "headscale"
    region_name: "Headscale Embedded DERP"
    stun_listen_addr: "0.0.0.0:3478"
  urls:
  paths:
    - /etc/headscale/derp.yaml # 此处需要引用derp配置
  auto_update_enabled: true
  update_frequency: 24h
disable_check_updates: true
ephemeral_node_inactivity_timeout: 30m
node_update_check_interval: 10s
db_type: sqlite3
db_path: /var/lib/headscale/db.sqlite
acme_url: https://acme-v02.api.letsencrypt.org/directory
acme_email: ""
tls_letsencrypt_hostname: ""
tls_letsencrypt_cache_dir: /var/lib/headscale/cache
tls_letsencrypt_challenge_type: HTTP-01
tls_letsencrypt_listen: ":http"
tls_cert_path: "/var/lib/headscale/cert/xxx.com.crt" # 填写证书文件绝对引用
tls_key_path: "/var/lib/headscale/cert/xxx.com.key"
log:
  format: text
  level: info
acl_policy_path: ""
dns_config:
  override_local_dns: true
  nameservers:
    - 223.6.6.6
  domains: []
  magic_dns: true
  base_domain: xxx.com # 可改可不改
unix_socket: /var/lib/headscale/headscale.sock # 修改位置
unix_socket_permission: "0770"
logtail:
  enabled: false
randomize_client_port: true # 改成当前值
EOF

7.将域名解析到当前服务器，并开放 headscale 的 tcp:8080、derp 的 tcp:8443,udp:3478 端口

8.创建数据目录及必须的文件
mkdir /data/data
touch /data/data/db.sqlite

9.启动服务
docker-compose up -d

# 理论上来说服务端只需要部署 headscale 即可，但是由于免费的 derp 节点都在国外，所以延迟较高，
# 并且流量经过公共服务器对于企业来讲是不安全的，所以需要自建 derp 服务器，derp 和 headscale 一样都是开源的
# DERP_VERIFY_CLIENTS=true 防白嫖模式 启用此模式后只有 tailscale 中的账户才可以使用此 derp，
# 所以需要依赖 tailscale 服务，derp 的其他配置详见 https://hub.docker.com/r/fredliang/derper
# 如果不启用此模式，tailscale 可以不安装
# tailscale docker 版 配置详见 https://hub.docker.com/r/tailscale/tailscale
# 上述的镜像均可以在 dockerhub 找到
```

## 三、服务端操作

*下列为常用操作，完整命令详见`headscale -h` `tailscale -h`*


### 1）tailscale注册至headscale

```bash
1.进入 tailscale 容器
docker exec -it tailscale sh

2.使用启动命令，此时获取到登陆密钥，复制 'key:' 后面的内容
tailscale up --login-server=https://headscale.xxxxx.com:8080 --accept-routes=true --accept-dns=false

3.新开终端进入 headscale 容器
docker exec -it headscale bash

4.创建一个租户，此处类似 k8s 的名称空间，用来做隔离的，名称自定义
headscale users create <名称>

5.注册节点，tailscale 界面显示 'Machine <hostname> registered' 则注册成功
headscale nodes register --user <此处为上面创建的租户名称> --key nodekey:<此处为上面复制的key>
```

### 2）设备操作
```bash
1.查看当前租户
headscale users list

2.查看已注册的节点
headscale node list

3.重命名节点
headscale nodes rename <NEW_NAME> -i <ID>

4.注销节点
headscale nodes expire -i <ID>

4.删除节点
headscale nodes delete -i <ID>
```

### 3）路由操作
```bash
1.进入 tailscale 容器
docker exec -it tailscale sh

2.覆盖注册信息，并添加路由信息
tailscale up --login-server=https://headscale.xxxxx.com:8080 --accept-routes=true --accept-dns=false --advertise-routes=192.168.0.0/16 --reset

3.新开终端进入 headscale 容器
docker exec -it headscale bash

4.查看路由信息
headscale routes list

5.启用路由
headscale routes enable -r <刚注册的路由 ID>

# 由于 headscale 自身的网段原因，`10.188.0.0/16`这个网段无法路由，
# 即使注册路由信息，开启后各客户端也没有相应的路由表
# 此处还有一个坑就是容器部署路由有问题
```


*在 linux 上直接部署 tailscale*
```bash
1.centos7 安装 tailscale
# 其他版本的 linux 参考 https://pkgs.tailscale.com/stable/#centos-7
# Install yum-config-manager if missing
sudo yum install yum-utils
# Add the tailscale repository
sudo yum-config-manager --add-repo https://pkgs.tailscale.com/stable/centos/7/tailscale.repo
# Install Tailscale
sudo yum install tailscale
# Enable and start tailscaled
sudo systemctl enable --now tailscaled

2.在客户端启用 ip 转发
# 参考 https://tailscale.com/kb/1019/subnets/?tab=linux#enable-ip-forwarding
echo 'net.ipv4.ip_forward = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf
echo 'net.ipv6.conf.all.forwarding = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf
sudo sysctl -p /etc/sysctl.d/99-tailscale.conf

3.注册至服务端并携带路由信息，复制 'key:' 后面的内容
tailscale up --login-server=https://headscale.xxxxx.com:8080 --accept-routes=true --accept-dns=false --advertise-routes=192.168.0.0/16

4.新开终端进入 headscale 容器
docker exec -it headscale bash

5.注册节点，tailscale 界面显示 'Machine <hostname> registered' 则注册成功
headscale nodes register --user <此处为上面创建的租户名称> --key nodekey:<此处为上面复制的key>

6.查看路由信息
headscale routes list

7.启用路由
headscale routes enable -r <刚注册的路由 ID>
```

## 四、排错

### 1）derp

*首先访问 derp 的域名 https://derp.xxxxx.com:8443，显示如下内容则服务正常，
如果不能显示，则查看其日志*

```bash
DERP

This is a Tailscale DERP server. 
```

### 2）headscale

*首先访问 headscale 的域名 https://headscale.xxxxx.com:8080/apple ，显示如下内容则服务正常，
如果不能显示，则查看其日志*

```bash
headscale: macOS configuration
Recent Tailscale versions (1.34.0 and higher)
    ......
```

### 3）tailscale

1.网络检查`tailscale netcheck`
```bash
2023/11/22 16:46:34 tlsdial: warning: server cert for "derp.xxxxx.com" is not a Let's Encrypt cert

Report:
	* UDP: false
	* IPv4: (no addr found)
	* IPv6: no, but OS has support
	* MappingVariesByDestIP:
	* HairPinning:
	* PortMapping:
	* Nearest DERP: aliyund-derp
	* DERP latency:
		- aliyun: 10ms    (aliyund-derp) # 显示为我们自定义的 derp 节点
```

2.ping 其他节点 `tailscale ping 10.188.0.3`，
也可以使用路由网段的任意内网 IP，会转发到注册路由的节点上，而后由该节点进行路由。
当现实如下信息基本就大功告成了，现在可以尝试连接内网服务器，或者在浏览器上打开部署
在内网的服务，应该都可以正常的访问。当然，前提是当前机器也注册到了 headscale

```bash
pong from izbp1223nrpyxhfv49wy8zz (fd7a:115c:a1e0::3) via DERP(aliyun) in 20ms
pong from izbp1223nrpyxhfv49wy8zz (fd7a:115c:a1e0::3) via DERP(aliyun) in 19ms
pong from izbp1223nrpyxhfv49wy8zz (fd7a:115c:a1e0::3) via DERP(aliyun) in 21ms
pong from izbp1223nrpyxhfv49wy8zz (fd7a:115c:a1e0::3) via DERP(aliyun) in 18ms
pong from izbp1223nrpyxhfv49wy8zz (fd7a:115c:a1e0::3) via DERP(aliyun) in 20ms
```

## 五、tailscale常用客户端安装

### 1）mac客户端

*打开 https://headscale.xxxxx.com:8080/apple*

```bash
1.客户端安装
- 直接在apple store中安装，需要非国区账号
- 使用文件直接安装 https://pkgs.tailscale.com/stable/#macos

2.描述文件
根据你的安装方式在页面上选择 macOS AppStore profile 或者 macOS Standalone profile 下载，
双击打开描述文件，然后到 设置-->隐私与安全性-->其他-->描述文件-->安装 headscale

3.打开 app，点击登陆，这时候会打开一个页面，上面写着用来认证的 key

4.在 headscale 中注册该节点

5.app 上点点点就完成了
```

### 1）windows客户端

*打开 https://headscale.xxxxx.com:8080/windows*

```bash
1.客户端安装
打开 https://tailscale.com/download/windows 点击下载安装

2.修改注册表
在页面上点击 Windows registry file 下载，然后双击运行

3.打开 app，点击登陆，这时候会打开一个页面，上面写着用来认证的 key

4.在 headscale 中注册该节点

5.app 上点点点就完成了
```

### 3）其他

linux：
   - https://pkgs.tailscale.com/stable/#centos-7
   - https://tailscale.com/download/linux

docker：https://hub.docker.com/r/tailscale/tailscale

ios：https://tailscale.com/download/ios

android：https://tailscale.com/download/android


## 六、参考文档：

[tailscale官网](https://tailscale.com/)

[headscale官网](https://headscale.net/)

[Tailscale-搭建异地局域网开源版中文部署指南](https://blog.csdn.net/github_36665118/article/details/128733646)

[Tailscale 基础教程：Headscale 的部署方法和使用教程](https://icloudnative.io/posts/how-to-set-up-or-migrate-headscale)

[Headscale 端到端直连](https://www.cnblogs.com/Yogile/p/17064031.html)