# Service

`将运行在一组 Pods 上的应用程序公开为网络服务的抽象方法。`



## 一、service资源

### service的实现

`随着pod的更新，pod的ip始终在变化，不可能使用人工去维护,所以需要引入一个负载均衡，去代理常常变化的pod`

~~~bash
apiVersion: v1
kind: Service
metadata:
  labels:
    name: nginx
  name: nginx-proxy
  namespace: dnet-integration-test
spec:
  type: ClusterIP
  ports:
  - name: nginx
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    name: nginx
---
apiVersion: v1
kind: Endpoints
metadata:
  labels:
    name: nginx
  name: nginx-proxy
  namespace: dnet-integration-test
subsets:
- addresses:
  - ip: 172.16.162.109
    nodeName: cn-hangzhou.172.16.246.154
    targetRef:
      kind: Pod
      name: nginx-proxy-5b4fbdb4f-x7htb
      namespace: dnet-integration-test
  ports:
  - name: nginx
    port: 80
    protocol: TCP

1.service并不是直接关联pod，而是通过endpoints去维护
2.service是由kube-proxy实现了一种 VIP（虚拟 IP），而不是一个实体服务
3.Endpoint在service使用标签选择器时自动创建，如果创建service时未指定标签 则需要手动创建endpoint 并且endpoint名称必须与service完全一致
# 从Kubernetes v1.21 [stable]开始，Endpoint Slices替代了Endpoints，为了应对大集群中endpoints记录的IP个数超出1000的限制值
~~~

### 服务发现

`使用 Kubernetes API 查询 API 服务器 的 Endpoints 资源进行服务发现，Endpoints会根据Pod的变更动态更新`

~~~bash
# 服务发现的两种方式
1.环境变量，在pod启动的时候kubelet 在pod中为每个活跃的 Service 添加一组环境变量，分别为Docker links兼容变量和{SVCNAME}_SERVICE_HOST、{SVCNAME}_SERVICE_PORT变量
  REDIS_MASTER_SERVICE_HOST=10.0.0.11
  REDIS_MASTER_SERVICE_PORT=6379
  REDIS_MASTER_PORT=tcp://10.0.0.11:6379
  REDIS_MASTER_PORT_6379_TCP=tcp://10.0.0.11:6379
  REDIS_MASTER_PORT_6379_TCP_PROTO=tcp
  REDIS_MASTER_PORT_6379_TCP_PORT=6379
  REDIS_MASTER_PORT_6379_TCP_ADDR=10.0.0.11

2.DNS，CoreDNS监视Kubernetes API中的新服务，并为每个服务创建一组DNS记录

# 补充知识
FQDN 全限定域名
  1) 查看DNS配置
  root@omp-portal-77b7674f97-wwznw:/data# cat /etc/resolv.conf 
  nameserver 10.96.0.10
  search dev.svc.cluster.local svc.cluster.local cluster.local
  options ndots:5
   
  2) 在pod中ping service falcon-schedule 
  root@omp-portal-77b7674f97-wwznw:/data# ping falcon-schedule
  PING falcon-schedule.dev.svc.cluster.local (10.103.136.14) 56(84) bytes of data.
  64 bytes from falcon-schedule.dev.svc.cluster.local (10.103.136.14): icmp_seq=1 ttl=64 time=0.117 ms
  64 bytes from falcon-schedule.dev.svc.cluster.local (10.103.136.14): icmp_seq=2 ttl=64 time=0.066 ms
  64 bytes from falcon-schedule.dev.svc.cluster.local (10.103.136.14): icmp_seq=3 ttl=64 time=0.127 ms
   
  3) DNS服务器找到其全限定域名 falcon-schedule.dev.svc.cluster.local 并解析IP
    falcon-schedule  主机名 
    dev 名称空间
    svc.cluster.local 集群域后缀 也可以理解为域名
   
  # 例：
  root@omp-portal-77b7674f97-wwznw:/data# curl http://gateway:8000
  Welcome to api-gateway
   
  root@omp-portal-77b7674f97-wwznw:/data# curl http://gateway.dev:8000
  Welcome to api-gateway
   
  root@omp-portal-77b7674f97-wwznw:/data# curl http://gateway.dev.svc.cluster.local:8000
  Welcome to api-gateway
~~~

## 二、service的工作模式

### userspace 代理模式

![](https://d33wubrfki0l68.cloudfront.net/e351b830334b8622a700a8da6568cb081c464a9b/13020/images/docs/services-userspace-overview.svg)

- 1.kube-proxy 会监视 Kubernetes 控制平面对 Service 对象和 Endpoints 对象的添加和移除操作。

- 2.对每个 Service，它会在本地 Node 上打开一个端口（随机选择）。 任何连接到“代理端口”的请求，都会被代理到 Service 的后端 Pods 中的某个上面（如 Endpoints 所报告的一样）。 使用哪个后端 Pod，是 kube-proxy 基于 SessionAffinity 来确定的。

- 3.最后，它配置 iptables 规则，捕获到达该 Service 的 clusterIP（是虚拟 IP） 和 Port 的请求，并重定向到代理端口，代理端口再代理请求到后端Pod。

- 4.默认情况下，用户空间模式下的 kube-proxy 通过轮转算法选择后端。

### iptables 代理模式

![](https://d33wubrfki0l68.cloudfront.net/27b2978647a8d7bdc2a96b213f0c0d3242ef9ce0/e8c9b/images/docs/services-iptables-overview.svg)

- 1.kube-proxy 会监视 Kubernetes 控制节点对 Service 对象和 Endpoints 对象的添加和移除。 对每个 Service，它会配置 iptables 规则，从而捕获到达该 Service 的 clusterIP 和端口的请求，进而将请求重定向到 Service 的一组后端中的某个 Pod 上面。 

- 2.对于每个 Endpoints 对象，它也会配置 iptables 规则，这个规则会选择一个后端组合。默认的策略是，kube-proxy 在 iptables 模式下随机选择一个后端。

- 3.使用 iptables 处理流量具有较低的系统开销，因为流量由 Linux netfilter 处理， 而无需在用户空间和内核空间之间切换。 这种方法也可能更可靠。

- 4.如果 kube-proxy 在 iptables 模式下运行，并且所选的第一个 Pod 没有响应， 则连接失败。 这与用户空间模式不同：在这种情况下，kube-proxy 将检测到与第一个 Pod 的连接已失败， 并会自动使用其他后端 Pod 重试。

### IPVS 代理模式

![](https://d33wubrfki0l68.cloudfront.net/2d3d2b521cf7f9ff83238218dac1c019c270b1ed/9ac5c/images/docs/services-ipvs-overview.svg)

`FEATURE STATE: Kubernetes v1.11 [stable]`

- 1.kube-proxy 监视 Kubernetes 服务和端点，调用 netlink 接口相应地创建 IPVS 规则， 并定期将 IPVS 规则与 Kubernetes 服务和端点同步。 该控制循环可确保IPVS 状态与所需状态匹配。访问服务时，IPVS 将流量定向到后端Pod之一。

- 2.IPVS代理模式基于类似于 iptables 模式的 netfilter 挂钩函数， 但是使用哈希表作为基础数据结构，并且在内核空间中工作。 这意味着，与 iptables 模式下的 kube-proxy 相比，IPVS 模式下的 kube-proxy 重定向通信的延迟要短，并且在同步代理规则时具有更好的性能。 与其他代理模式相比，IPVS 模式还支持更高的网络流量吞吐量。

~~~bash
IPVS 提供了更多选项来平衡后端 Pod 的流量。 这些是：

  - rr：轮替（Round-Robin）
  - lc：最少链接（Least Connection），即打开链接数量最少者优先
  - dh：目标地址哈希（Destination Hashing）
  - sh：源地址哈希（Source Hashing）
  - sed：最短预期延迟（Shortest Expected Delay）
  - nq：从不排队（Never Queue）
~~~


## 三、service的类型

### 无头服务（Headless Services）
- None

~~~bash
apiVersion: v1
kind: Service
metadata:
  labels:
    name: nginx
  name: nginx-proxy
  namespace: dnet-integration-test
spec:
  type: None # 无头服务
  ports:
  - name: nginx
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    name: nginx

1.None类型的无头service，kube-proxy不会处理他们，但是依然在集群内通过服务发现进行后端pod的代理
2.即使service配置中没有定义标签选择器，DNS依然会试图查找同名的endpoints资源
~~~

- ExternalName

~~~bash
apiVersion: v1
kind: Service
metadata:
  name: my-service
  namespace: prod
spec:
  type: ExternalName
  externalName: my.database.example.com

当集群请求my-service时返回的不是IP而是my.database.example.com这个域名
~~~

### 集群IP
- ClusterIP

~~~bash
apiVersion: v1
kind: Service
metadata:
  labels:
    name: nginx
  name: nginx-proxy
  namespace: dnet-integration-test
spec:
  ports:
  - name: nginx
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    name: nginx

1.ClusterIP是默认类型，配置ClusterIP或者没有配置service类型时都是使用ClusterIP
2.通过集群的内部 IP 暴露服务，选择该值时服务只能够在集群内部访问。
~~~

### 节点端口
- NodePort

~~~bash
apiVersion: v1
kind: Service
metadata:
  labels:
    name: nginx
  name: nginx-proxy
  namespace: dnet-integration-test
spec:
  tpye: NodePort # nodeport类型
  ports:
  - name: nginx
    port: 80
    protocol: TCP
    targetPort: 80
    nodePort: 30002 # 定义主机上暴露的端口
  selector:
    name: nginx

1.使用nodeport类型时，如果定义了端口，那么会在集群中的所有主机上暴露该端口，如果没有定义则会生成随机端口（默认值：30000-32767）
2.使用任意nodeIP+nodeport生成的端口即可访问集群中的服务
~~~

### 负载均衡
- LoadBalancer

~~~bash
apiVersion: v1
kind: Service
metadata:
  labels:
    name: nginx
  name: nginx-proxy
  namespace: dnet-integration-test
spec:
  tpye: LoadBalancer # LoadBalancer类型
  ports:
  - name: nginx
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    name: nginx

1.使用云提供商的负载均衡器向外部暴露服务。
2.外部负载均衡器可以将流量路由到自动创建的 NodePort 服务和 ClusterIP 服务上。
3.kube-apiserver 启用了 MixedProtocolLBService 配置选项， 则当定义了多个端口时，允许使用不同的协议，但协议类型仍然由云提供商决定。

~~~

## 四、细节配置

### 会话保持
~~~bash
apiVersion: v1
kind: Service
spec:
  sessionAffinity: ClientIP
 
- None  不设置此项时的默认值
- ClientIP  配置会话保持

- 为什么没有基于cookie的会话保持？
  cookie是HTTP协议中的一部分，service不是在HTTP层面上工作。service处理TCP和UDP包，并不关心其中的载荷内容
~~~

### 真实IP
- 外部流量策略

~~~bash
apiVersion: v1
kind: Service
spec:
  externalTrafficPolicy: Local

- Local
- Cluster

1.默认配置为cluster，kube-proxy会将流量分配至所有节点，再由节点代理至当前节点或者其他节点的pod上
2.配置为Local时，kube-proxy将直接爸流量分配至有pod的节点，不会在经过其他节点路由

# 补充
在这种情况下在本地pod终止或短暂失败的情况下会出现流量丢失的情况，在Kubernetes v1.22 [alpha] 启用了 kube-proxy 的 ProxyTerminatingEndpoints，可以在这种情况下将流量转发至健康的节点。
~~~

- 内部流量策略
`Kubernetes v1.22 [beta]`

~~~bash
apiVersion: v1
kind: Service
spec:
  internalTrafficPolicy: Local

- Local
- Cluster

1.将字段设置为 Cluster 会将内部流量路由到所有就绪端点，设置为 Local 只会路由到当前节点上就绪的端点。 
2.如果流量策略是 Local，而且当前节点上没有就绪的端点，那么 kube-proxy 会丢弃流量。
~~~

## 五、CNI
`k8s中的网络同样是以插件的方式部署在集群中的`

### 常见的CNI插件
~~~bash
1.Flannel 高性能，低开销，不支持网络策略，不支持加密，适合中小集群。

2.Calico  高性能，开销中等，支持网络策略，加密性能好，适合中大型集群。

3.Cilium  高性能，开销中等，支持网络策略，加密性能一般，适合中大型集群。

# 需要更换其他cni时需要清空/etc/cni/net.d/目录
~~~

### pod内部调用路径
`以Flannel为例`

~~~bash

~~~
