# Kubeadm 安装


## 一、环境准备（最低双核2G内存）

| **主机名** |    **IP**     | **角色** |
| :--------: | :-----------: | :------: |
|  master-1  | 172.17.11.150 | Master-1 |
|  master-2  | 172.17.11.151 | master-2 |
|  master-3  | 172.17.11.152 | master-3 |
|   node-3   | 172.17.11.153 |  Node-3  |
|            | 172.17.11.155 |   VIP    |

## 二、基础优化

~~~bash
1.关闭selinux、firewalld
setenforce 0
vim /etc/sysconfig/selinux
vim /etc/selinux/config
systemctl stop firewalld
systemctl disable firewalld

2.换源
mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo

3.安装依赖
yum install wget expect vim net-tools ntp bash-completion ipvsadm ipset jq iptables conntrack sysstat libseccomp -y
 
 
4.关闭swap
swapoff -a
#完全关闭 vim /etc/fstab

5.修改主机名
hostnamectl set-hostname 


6.配置域名解析
cat >> /etc/hosts << EOF
172.17.11.150 master-1
172.17.11.151 node-1
172.17.11.152 node-2
172.17.11.153 node-3
EOF


7.时间同步
yum install -y ntpdate
crontab -l
*/5 * * * * ntpdate ntp.aliyun.com &> /dev/null


8.修改内核参数
cat > /etc/sysctl.d/kubernetes.conf <<EOF
net.bridge.bridge-nf-call-iptables=1
net.bridge.bridge-nf-call-ip6tables=1
net.ipv4.ip_forward=1
net.ipv4.tcp_tw_recycle=0
vm.swappiness=0
vm.overcommit_memory=1
vm.panic_on_oom=0
fs.inotify.max_user_instances=8192
fs.inotify.max_user_watches=1048576
fs.file-max=52706963
fs.nr_open=52706963
net.ipv6.conf.all.disable_ipv6=1
net.netfilter.nf_conntrack_max=2310720
EOF

sysctl -p /etc/sysctl.d/kubernetes.conf
# 如有报错可以忽略

9.安装IPVS
cat > /etc/sysconfig/modules/k8s.modules <<EOF
modprobe -- ip_vs
modprobe -- ip_vs_rr
modprobe -- ip_vs_wrr
modprobe -- ip_vs_sh
modprobe -- nf_conntrack_ipv4
modprobe -- ip_tables
modprobe -- ip_set
modprobe -- xt_set
modprobe -- ipt_set
modprobe -- ipt_rpfilter
modprobe -- ipt_REJECT
modprobe -- ipip
EOF

# 查看是否加载,如未看到信息可以重启试试。
chmod 755 /etc/sysconfig/modules/k8s.modules && bash /etc/sysconfig/modules/k8s.modules && lsmod | grep -e ip_vs -e nf_conntrack_ipv4


10.安装Docker
yum install -y yum-utils device-mapper-persistent-data lvm2
yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum makecache fast
yum -y install docker-ce
systemctl enable --now docker


11.配置docker加速并修改驱动
cat > daemon.json <<EOF
{
    "exec-opts": ["native.cgroupdriver=systemd"],
    "registry-mirrors": [
        "http://hub-mirror.c.163.com",
        "https://docker.mirrors.ustc.edu.cn",
        "http://f1361db2.m.daocloud.io",
        "https://registry.docker-cn.com"
    ]
}
EOF
systemctl restart docker


12.配置kubernetes源
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
~~~

## 三、镜像准备

~~~bash
# 由于k8s镜像资源都在海外站点 直接使用时会出现`镜像拉取失败`的报错 所以需要先行手动将镜像下载至本地 再初始化

1.安装k8s相关组件
yum install -y kubelet kubeadm kubectl
2.设置开机自启并启动kubelet
systemctl enable --now kubelet.service
3.配置初始化文件 #master-1操作
#生成
kubeadm config print init-defaults > kubeadm.yaml
#配置 （此处可以看到当前版本号）
vim kubeadm.yaml
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
kubernetesVersion: v1.21.1
controlPlaneEndpoint: 172.17.11.155:8443
imageRepository: k8s.gcr.io
apiServer:
  certSANs:
   - 172.17.11.150
   - 172.17.11.151
   - 172.17.11.152
   - 172.17.11.153
   - 172.17.11.155
dns:
  type: CoreDNS
networking:
 podSubnet: 10.244.0.0/16
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: ipvs

4.根据版本号查找所需镜像
kubeadm config images list --kubernetes-version=1.20.0
# 结果
k8s.gcr.io/kube-apiserver:v1.21.1
k8s.gcr.io/kube-controller-manager:v1.21.1
k8s.gcr.io/kube-scheduler:v1.21.1
k8s.gcr.io/kube-proxy:v1.21.1
k8s.gcr.io/pause:3.4.1
k8s.gcr.io/etcd:3.4.13-0
k8s.gcr.io/coredns/coredns:v1.8.0

5.编写镜像拉取脚本并执行
cat >images.txt<<EOF 
kube-apiserver
kube-controller-manager
kube-scheduler
kube-proxy
pause
etcd
coredns
EOF

cat >images.sh<<'EOF' 
k8s_version="v1.21.1"
pause_version="v3.4.1"
etcd_version="v3.4.13-0"
coredns_version="v1.8.0"

for image in `cat images.txt`
do
    if [[ ${image} == "pause" ]];then
    docker pull registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image
    docker tag registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image k8s.gcr.io/$image:${pause_version}
    docker rmi registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image
    elif [[ ${image} == "etcd" ]];then
    docker pull registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image
    docker tag registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image k8s.gcr.io/$image:${etcd_version}
    docker rmi registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image
    elif [[ ${image} == "coredns" ]];then
    docker pull registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image
    docker tag registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image k8s.gcr.io/$image/coredns:${coredns_version}
    docker rmi registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image
    else
    docker pull registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image
    docker tag registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image k8s.gcr.io/$image:${k8s_version}
    docker rmi registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:$image
   fi
done
EOF

sh images.sh

6.确认
docker images

# 如果需要下载特定版本的镜像而别人的镜像仓库没有 则可以使用阿里云`容器镜像服务`构建自己的镜像
~~~

## 四、keepalive（master节点）

~~~bash
1.安装
yum install -y socat keepalived ipvsadm conntrack

2.修改配置文件
vim /etc/keepalived/keepalived.conf
global_defs {
   router_id LVS_DEVEL
}
vrrp_instance VI_1 {
    state BACKUP
    nopreempt
    interface eth1
    virtual_router_id 80
    # 其余master节点配置于此一致 唯一需要修改此优先级
    priority 100
    advert_int 1
    authentication {
        auth_type PASS
        auth_pass just0kk
    }
    virtual_ipaddress {
        172.17.11.155
    }
}
virtual_server 172.17.11.155 6443 {
    delay_loop 6
    lb_algo loadbalance
    lb_kind DR
    net_mask 255.255.255.0
    persistence_timeout 0
    protocol TCP
    real_server 172.17.11.150 6443 {
        weight 1
        SSL_GET {
            url {
              path /healthz
              status_code 200
            }
            connect_timeout 3
            nb_get_retry 3
            delay_before_retry 3
        }
    }
    real_server 172.17.11.151 6443 {
        weight 1
        SSL_GET {
            url {
              path /healthz
              status_code 200
            }
            connect_timeout 3
            nb_get_retry 3
            delay_before_retry 3
        }
    }
    real_server 172.17.11.152 6443 {
        weight 1
        SSL_GET {
            url {
              path /healthz
              status_code 200
            }
            connect_timeout 3
            nb_get_retry 3
            delay_before_retry 3
        }
    }
}

3.启动keepalive
systemctl enable --now keepalived  && systemctl status keepalived

# 此keepalive配置为非抢占式 原因是当节点宕机后 重启节点 优先级高的抢占式节点会抢回VIP 此时apiserver并未提供服务 但是流量已经接入 导致数据丢失
~~~

## 五、haproxy（master节点）

~~~bash
1.安装haproxy
yum install -y haproxy

2.配置haproxy
vim /etc/haproxy/haproxy.cfg
global
  maxconn  2000
  ulimit-n  16384
  log  127.0.0.1 local0 err
  stats timeout 30s

defaults
  log global
  mode  http
  option  httplog
  timeout connect 5000
  timeout client  50000
  timeout server  50000
  timeout http-request 15s
  timeout http-keep-alive 15s

frontend monitor-in
  bind *:33305
  mode http
  option httplog
  monitor-uri /monitor

listen stats
  bind    *:8006
  mode    http
  stats   enable
  stats   hide-version
  stats   uri       /stats
  stats   refresh   30s
  stats   realm     Haproxy\ Statistics
  stats   auth      admin:admin

frontend k8s-master
  bind 0.0.0.0:8443
  bind 127.0.0.1:8443
  mode tcp
  option tcplog
  tcp-request inspect-delay 5s
  default_backend k8s-master

backend k8s-master
  mode tcp
  option tcplog
  option tcp-check
  balance roundrobin
  default-server inter 10s downinter 5s rise 2 fall 2 slowstart 60s maxconn 250 maxqueue 256 weight 100
  server master-1    172.17.11.150:6443  check inter 2000 fall 2 rise 2 weight 100
  server master-2    172.17.11.151:6443  check inter 2000 fall 2 rise 2 weight 100
  server master-3    172.17.11.152:6443  check inter 2000 fall 2 rise 2 weight 100

3.启动haproxy
systemctl enable --now haproxy

~~~

## 六、创建集群

~~~bash
# master-1上操作
1.根据配置文件初始化
kubeadm init --config=kubeadm.yaml

2.配置kubectl与kube-apiserver交互
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config


3.安装网络组件
wget https://docs.projectcalico.org/v3.19/manifests/calico.yaml
kubectl apply -f calico.yaml


4.查看集群状态（稍等片刻）
kubectl get nodes

5.查组件状态
kubectl get cs
# 结果：
NAME                 STATUS      MESSAGE                                                                                     ERROR
scheduler            Unhealthy   Get http://127.0.0.1:10251/healthz: dial tcp 127.0.0.1:10251: connect: connection refused   
controller-manager   Unhealthy   Get http://127.0.0.1:10252/healthz: dial tcp 127.0.0.1:10252: connect: connection refused   
etcd-0               Healthy     {"health":"true"}

#解决方法：
vim /etc/kubernetes/manifests/kube-controller-manager.yaml
#    - --port=0
vim /etc/kubernetes/manifests/kube-scheduler.yaml
#    - --port=0


6.查组件状态（稍等片刻）
kubectl get cs

NAME                 STATUS    MESSAGE             ERROR
scheduler            Healthy   ok                  
controller-manager   Healthy   ok                  
etcd-0               Healthy   {"health":"true"}


7.将证书同步至其他master
for i in master-2 master-3;do
scp /etc/kubernetes/pki/ca.* ${i}:/etc/kubernetes/pki/
scp /etc/kubernetes/pki/sa.* ${i}:/etc/kubernetes/pki/
scp /etc/kubernetes/pki/front-proxy-ca.* ${i}:/etc/kubernetes/pki/
scp /etc/kubernetes/pki/etcd/ca.* ${i}:/etc/kubernetes/pki/etcd/
done

8.将其余master添加至集群
# 将提示中的此段命令复制至其他master上执行
kubeadm join 172.16.1.80:6443 --token abcdef.0123456789abcdef \
    --discovery-token-ca-cert-hash sha256:0f55ea798f2c733574fa0ed0ac4a2b489671cc6b93345162966f026056cc566a \
    --control-plane

# 根据提示操作
    mkdir -p $HOME/.kube
	  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
	  sudo chown $(id -u):$(id -g) $HOME/.kube/config

9.将node节点添加至集群
# 将提示中的此段命令复制至其他node上执行
kubeadm join 172.16.1.80:6443 --token abcdef.0123456789abcdef \
    --discovery-token-ca-cert-hash sha256:0f55ea798f2c733574fa0ed0ac4a2b489671cc6b93345162966f026056cc566a
~~~

## 七、检查

~~~bash
# master输入：
kubectl get nodes
kubectl get pods -A
~~~

## 八、附加组件安装（可选）

~~~bash
# metrics-server (HPA依赖此组件)
1.安装helm
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash

2.添加helm源
helm repo add apphub https://apphub.aliyuncs.com

3.在repo中搜索metrics-server
helm search repo metrics-server

4.安装
helm install metrics-server  apphub/metrics-server -n kube-system --set apiService.create=true

5.修改metrics-server deployment
kubectl edit deployment.apps/metrics-server -n kube-system
  在secure-port后加入两条参数
...
    spec:
      containers:
      - command:
        - metrics-server
        - --secure-port=8443
        - --kubelet-insecure-tls # 解决无法解析主机IP问题
        - --kubelet-preferred-address-types=InternalIP # 跳过证书验证
...


~~~





