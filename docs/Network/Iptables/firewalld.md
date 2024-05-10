防火墙安全概述
    在CentOS7系统中集成了多款防火墙管理工具，默认启用的是firewalld（动态防火墙管理器）防火墙管理工具，Firewalld支持CLI（命令行）以及GUI（图形）的两种管理方式。
    对于接触Linux较早的人员对Iptables比较熟悉，但由于Iptables的规则比较的麻烦，并且对网络有一定要求，所以学习成本较高。但firewalld的学习对网络并没有那么高的要求，相对iptables来说要简单不少，所以建议刚接触CentOS7系统的人员直接学习Firewalld。
    需要注意的是：如果开启防火墙工具，并且没有配置任何允许的规则，那么从外部访问防火墙设备默认会被阻止，但是如果直接从防火墙内部往外部流出的流量默认会被允许。
    firewalld 只能做和IP/Port相关的限制，web相关的限制无法实现。
防火墙使用区域管理
    防火墙区域

| 区域选项    |                         默认规则策略                         |
| ----------- | :----------------------------------------------------------: |
| **trusted** |                   允许所有的数据包流入流出                   |
| home        | 拒绝流入的流量，除非与流出的流量相关；而如果流量与ssh、mdns、ipp-client、amba-client与dhcpv6-client服务相关，则允许流量 |
| internal    |                        等同于home区域                        |
| work        | 拒绝流入的流量，除非与流出的流量相关；而如果流量与ssh、ipp-client、dhcpv6-client服务相关，则允许流量 |
| **public**  | 拒绝流入的流量，除非与流出的流量相关；而如果流量与ssh、dhcpv6-client服务相关，则允许流量 |
| external    | 拒绝流入的流量，除非与流出的流量相关；而如果流量与ssh服务相关，则允许流量 |
| dmz         | 拒绝流入的流量，除非与流出的流量相关；而如果流量与ssh服务相关，则允许流量 |
| block       |             拒绝流入的流量，除非与流出的流量相关             |
| **drop**    |             拒绝流入的流量，除非与流出的流量相关             |
    防火墙参数

|             参数              |                         作用                         |
| :---------------------------: | :--------------------------------------------------: |
|     **zone区域相关指令**      |                                                      |
|      --get-default-zone       |                  获取默认的区域名称                  |
| --set-default-zone=<区域名称> |             设置默认的区域，使其永久生效             |
|      --get-active-zones       |           显示当前正在使用的区域与网卡名称           |
|          --get-zones          |                  显示总共可用的区域                  |
|                               |                                                      |
|   **services服务相关命令**    |                                                      |
|        --get-services         |            列出服务列表中所有可管理的服务            |
|        --add-service=         |           设置默认区域允许该填加服务的流量           |
|       --remove-service=       |          设置默认区域不允许该删除服务的流量          |
|     **Port端口相关指令**      |                                                      |
|   --add-port=<端口号/协议>    |           设置默认区域允许该填加端口的流量           |
|  --remove-port=<端口号/协议>  |           置默认区域不允许该删除端口的流量           |
|   **Interface网站相关指令**   |                                                      |
|  --add-interface=<网卡名称>   |       将源自该网卡的所有流量都导向某个指定区域       |
| --change-interface=<网卡名称> |               将某个网卡与区域进行关联               |
|       **其他相关指令**        |                                                      |
|          --list-all           | 显示当前区域的网卡配置参数、资源、端口以及服务等信息 |
|           --reload            | 让“永久生效”的配置规则立即生效，并覆盖当前的配置规则 |
防火墙配置策略
为了能正常使用firewalld服务和相关工具去管理防火墙，必须启动firewalld服务，同时关闭以前旧的防火墙相关服务，需要注意firewalld的规则分为两种状态：
 
runtime运行时: 修改规则马上生效，但如果重启服务则马上失效，测试建议。
permanent持久配置: 修改规则后需要reload重载服务才会生效，生产建议。
 
#永久生效使用：
firewall-cmd --add-port=80/udp
success
	
firewall-cmd --add-port=80/udp --permanent
success
    禁用防火墙
#1. 禁用旧版防火墙服务或保证没启动
[root@web01 ~]# systemctl mask iptables
Created symlink from /etc/systemd/system/iptables.service to /dev/null.
[root@web01 ~]# systemctl mask ip6tables
Created symlink from /etc/systemd/system/ip6tables.service to /dev/null.
 
#2. 启动firewalld防火墙，并加入开机自启动服务
[root@m01 ~]# systemctl start firewalld
[root@m01 ~]# systemctl enable firewalld
 
#3.取消禁用防火墙
[root@web01 ~]# systemctl unmask iptables
    启动防火墙
[root@web01 services]# systemctl start firewalld
[root@web01 services]# systemctl enable firewalld
Created symlink from /etc/systemd/system/dbus-org.fedoraproject.FirewallD1.service to /usr/lib/systemd/system/firewalld.service.
Created symlink from /etc/systemd/system/multi-user.target.wants/firewalld.service to /usr/lib/systemd/system/firewalld.service.
    firewalld 常用命令
#1.查看默认使用的区域
[root@web01 services]# firewall-cmd --get-default-zone
public
 
#2.查看默认区域的规则
[root@web01 services]# firewall-cmd --list-all
public (active)						#区域名字（活跃的在使用的）
  target: default					#状态：默认
  icmp-block-inversion: no			#ICMP
  interfaces: eth0 eth1				#区域绑定的网卡
  sources: 							#允许的源IP
  services: ssh dhcpv6-client		#允许的服务
  ports: 							#允许的端口
  protocols: 						#允许的协议
  masquerade: no					#IP伪装
  forward-ports: 					#端口转发
  source-ports: 					#来源端口
  icmp-blocks: 						#icmp块
  rich rules: 						#富规则
 
#3.查看指定区域的规则
[root@web01 services]# firewall-cmd --list-all --zone=drop
drop
  target: DROP
  icmp-block-inversion: no
  interfaces: 
  sources: 
  services: 
  ports: 
  protocols: 
  masquerade: no
  forward-ports: 
  source-ports: 
  icmp-blocks: 
  rich rules:
  
#4.查询某区域是否允许某服务
[root@web01 services]# firewall-cmd --zone=public --query-service=ssh
yes
 
#5.重启防火墙
[root@web01 services]# firewall-cmd --reload
success
 
#6.同时配置多个服务
[root@web01 services]# firewall-cmd --add-service={ssh,httpd}
success
    配置测试
## 配置要求：调整防火墙，默认区域拒绝所有的流量，如果来源IP是10.0.0.0/24则允许
 
#移除public区域所有操作
[root@web01 services]# firewall-cmd --remove-service=ssh
success
[root@web01 services]# firewall-cmd --remove-service=dhcpv6-client
success
 
#配置允许的网段到trusted区域
[root@web01 services]# firewall-cmd --add-source=10.0.0.0/24 --zone=trusted 
success
防火墙配置放行策略
    firewalld放行服务
#放行mysql服务（firewalld已存在的服务）
[root@web01 ~]# firewall-cmd --add-service=mysql
success
 
#放行nginx服务（firewalld中不存在的服务）
[root@web01 ~]# cd /usr/lib/firewalld/services
[root@web01 services]# cp mysql.xml nginx.xml
[root@web01 services]# vim nginx.xml 
<?xml version="1.0" encoding="utf-8"?>
<service>
  <short>Nginx</short>
  <description>Nginx Server</description>
  <port protocol="tcp" port="80"/>
</service>
[root@web01 services]# firewall-cmd --reload
success
[root@web01 services]# firewall-cmd --add-service=nginx
success
[root@web01 services]#
    firewalld放行端口
[root@web01 services]# firewall-cmd --add-port=80/tcp
 
[root@web01 services]# firewall-cmd --list-all
public (active)
  target: default
  icmp-block-inversion: no
  interfaces: eth0 eth1
  sources: 
  services: ssh dhcpv6-client
  ports: 80/tcp 80/udp
  protocols: 
  masquerade: no
  forward-ports: 
  source-ports: 
  icmp-blocks: 
  rich rules:
    firewalld放行网段
[root@web01 services]# firewall-cmd --add-source=10.0.0.0/24 --zone=trusted
防火墙端口转发策略
端口转发是指传统的目标地址映射，实现外网访问内网资源，firewalld转发命令格式为：
firewalld-cmd --permanent --zone=<区域> --add-forward-port=port=<源端口号>:proto=<协议>:toport=<目标端口号>:toaddr=<目标IP地址>
 
如果需要将本地的10.0.0.7:5555端口转发至后端172.16.1.8:22端口
 
#1.配置端口转发
[root@web01 ~]# firewall-cmd --permanent --zone=public --add-forward-port=port=5555:proto=tcp:toport=22:toaddr=172.16.1.8
success
[root@web01 ~]# firewall-cmd --reload
 
#2.开启IP伪装
[root@web01 ~]# firewall-cmd --add-masquerade 
success
[root@web01 ~]# firewall-cmd --add-masquerade --permanent 
success
 
#3.测试访问
[root@m01 ~]# ssh 10.0.0.7 -p5555
root@10.0.0.7's password: 
Last login: Tue Jul  7 01:06:01 2020 from 172.16.1.7
防火墙富规则
    firewalld中的富语言规则表示更细致，更详细的防火墙策略配置，他可以针对系统服务、端口号、原地址和目标地址等诸多信息进行更有针对性的策略配置，优先级在所有的防火墙策略中也是最高的，下面为firewalld富语言规则帮助手册
    富规则语法
[root@web01 ~]# man firewalld.richlanguage
           rule
             [source]
             [destination]
             service|port|protocol|icmp-block|icmp-type|masquerade|forward-port|source-port
             [log]
             [audit]
             [accept|reject|drop|mark]
             
rule [family="ipv4|ipv6"]
source address="address[/mask]" [invert="True"]
service name="service name"
port port="port value" protocol="tcp|udp"
protocol value="protocol value"
forward-port port="port value" protocol="tcp|udp" to-port="port value" to-addr="address"
accept | reject [type="reject type"] | drop
    实例一：允许10.0.0.1主机能够访问http服务，允许172.16.1.0/24能访问22端口
#允许10.0.0.1主机能够访问http服务
[root@web01 ~]# firewall-cmd --add-rich-rule='rule family=ipv4 source address=10.0.0.1 service name=http accept'
success
 
#允许172.16.1.0/24能访问22端口
[root@web01 ~]# firewall-cmd --add-rich-rule='rule family=ipv4 source address=172.16.1.0/24 port port=22 protocol=tcp accept'
success
    实例二：默认public区域对外开放所有人能通过ssh服务连接，但拒绝172.16.1.0/24网段通过ssh连接服务器
firewall-cmd --add-rich-rule='rule family=ipv4 source address=172.16.1.0/24 service name=ssh drop' --permanent
 
firewall-cmd --reload
    实例三：使用firewalld，允许所有人能访问http,https服务，但只有10.0.0.1主机可以访问ssh服务
firewall-cmd --add-service={http,https} --permanent
 
firewall-cmd --remove-service=ssh --permanent
 
firewall-cmd --add-rich-rule='rule family=ipv4 source address=10.0.0.1  service name=ssh  accept' --permanent
firewall-cmd --reload
    实例四：当用户来源IP地址是10.0.0.1主机，则将用户请求的5555端口转发至后端172.16.1.7的22端口
firewall-cmd --permanent --zone=public --add-forward-port=port=5555:proto=tcp:toport=22:toaddr=172.16.1.7
 
firewall-cmd --add-masquerade --permanent 
 
firewall-cmd --reload
防火墙规则备份
我们所有针对public区域编写的永久添加的规则都会写入备份文件(--permanent) /etc/firewalld/zones/public.xml
 
#我们防火墙的配置，永久生效后会保存在 /etc/firewalld/zones/目录下，所以，以后进行服务器扩展，或者配置相同防火墙时，只需要拷贝该目录下的文件即可
 
备份也备份以上目录下的文件即可
防火墙内部共享上网
    开启防火墙IP伪装
[root@web01 ~]# firewall-cmd --add-masquerade --permanent 
success
[root@web01 ~]# firewall-cmd --add-masquerade
success
    防火墙开启内核转发（如果是Centos6，需要配置）
[root@web01 ~]# sysctl -a | grep net.ipv4.ip_forward
net.ipv4.ip_forward = 1
 
#配置内核转发
[root@m01 ~]# vim /etc/sysctl.conf
net.ipv4.ip_forward = 1
 
#在CentOS6中开启之后生效命令
[root@m01 ~]# sysctl -p
 
#查看内核转发是否开启
[root@m01 ~]# sysctl -a|grep net.ipv4.ip_forward
net.ipv4.ip_forward = 1
    配置没有外网的机器网关地址
[root@web01 ~]# vim /etc/sysconfig/network-scripts/ifcfg-eth1
#添加配置
GATEWAY=172.16.1.7    #防火墙内网地址
DNS1=223.5.5.5
 
#重启网卡
[root@web01 ~]# ifdown eth1
[root@web01 ~]# ifup eth1