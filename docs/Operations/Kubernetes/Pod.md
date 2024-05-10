# Pod


## 一、pod的概念

### pod是什么

- Pod直译过来就是豌豆荚，很形象的表达了pod的概念，而其中包含的豌豆就是容器。

- Pod是在K8s集群中运行部署应用或服务的最小单元，它是可以支持一个或多容器。

- Pod中的容器通过pause容器共享网络、存储、等资源。

- 通常情况下我们会使用控制器来创建Pod，而不是直接创建。

### pod带来的好处

- Pod做为一个可以独立运行的服务单元，简化了应用部署的难度，以更高的抽象层次为应用部署管提供了极大的方便。

- 做为最小的应用实例可以独立运行，因此可以方便的进行部署、水平扩展和收缩、方便进行调度管理与资源的分配。

- Pod中的容器共享相同的数据和网络地址空间，Pod之间也进行了统一的资源管理与分配。

### pause容器
- 每个Pod里运行着一个特殊的被称之为Pause的容器，其他容器则为业务容器，这些业务容器共享Pause容器的网络栈和Volume挂载卷，因此他们之间通信和数据交换更为高效，在设计时我们可以充分利用这一特性将一组密切相关的服务进程放入同一个Pod中。同一个Pod里的容器之间仅需通过localhost就能互相通信

~~~bash
- PID命名空间：Pod中的不同应用程序可以看到其他应用程序的进程ID。
   这里pause的pid为1，承担systemd进程的责任，创建子进程，孤儿进程管理，垃圾回收。

- 网络命名空间：Pod中的多个容器能够访问同一个IP和端口范围。

- IPC命名空间：Pod中的多个容器能够使用SystemV IPC或POSIX消息队列进行通信。

- UTS命名空间：Pod中的多个容器共享一个主机名；Volumes（共享存储卷）

- Pod中的各个容器可以访问在Pod级别定义的Volumes。

~~~

### pod的生命周期分析
~~~bash
1.create/apply pod

# Waiting
# Pending
2.绑定node(资源)、pod ip

3.拉取镜像

4.挂载存储、配置、密钥

# Running
5.初始化容器阶段初始化pod中每一个容器,他们是串行执行的，执行完成后就退出了

6.启动主容器main container

7.在main container刚刚启动之后可以执行post start命令

8.在整个main container执行的过程中可以做两类探测：startupliveness probe(存活探测)和readiness probe(就绪探测)

# Terminated
9.在main container结束前可以执行pre stop命令

10.deleted
~~~

## 二、相关配置分析

### yaml文件的定义
~~~bash
- 绝大多数由这五个字段构成
apiVersion: #api version || api group/api version
kind: # api类型
metadata: # 元数据 名称、标签、注释等
spec: # 规格 详细的定义
status: # 状态 由集群自动生成
~~~

### 最基础的nginx pod为例
~~~bash
apiVersion: apps/v1 #pod 属于app组 v1版本
kind: Pod # api类型为pod
metadata: # 元数据
  labels: # 标签
    name: nginx
  name: nginx # pod 名称
  namespace: default # 名称空间
spec: # 规格
  containers:
  - name: nginx # 容器名称
    image: nginx:1.19 # 容器使用的镜像
    ports: # 定义容器端口和协议
    - containerPort: 5678
      protocol: TCP

1.labels为非必须项，但没有labels将无法被控制器或资源进行关联管理
2.ports为非必需项，如果没有声明端口容器网络将为none
~~~

### 容器资源
~~~bash
apiVersion: v1
kind: Pod
metadata:
  labels:
    name: nginx
  name: nginx
  namespace: default
spec:
  containers:
  - name: nginx
    image: nginx:1.19
    ports:
    - containerPort: 5678
      protocol: TCP
    resources: # 定义资源限制
      limits: # 容器最多可以使用多少
        cpu: "1"
        memory: 2Gi
      requests: # 容器至少需要多少
        cpu: 50m
        memory: 1Gi

1.requests决定了pod的调度 如果所有node节点资源不满足 pod会处于Pending状态
2.当设置limits而没有设置requests时，Kubernetes 默认requests等于limits
3.对于limits，cpu使用会被限制，但是如果内存超过限制值则会被kernel OOM kill,此时kubernetes会重启该 container 或者在本机或其他节点上重新创建一个 pod。
~~~

### 镜像拉取
~~~bash
apiVersion: v1
kind: Pod
metadata:
  labels:
    name: nginx
  name: nginx
  namespace: default
spec:
  containers:
  - name: nginx
    image: nginx:1.19
    ports:
    - containerPort: 5678
      protocol: TCP
    resources:
      limits:
        cpu: "1"
        memory: 2Gi
      requests:
        cpu: 50m
        memory: 1Gi
    imagePullPolicy: Always # 镜像拉取策略
  imagePullSecrets: # harbor仓库密钥
  - name: harbor

1.镜像拉取策略有 IfNotPresent(如果本地没有则拉取)，Always(总是拉取)，Never(从不拉取) 三种，默认为IfNotPresent
2.harbor仓库密钥是全局的，不是容器级别，这里使用的是提前创建好的secret资源
~~~


### 存储及配置挂载 
~~~bash
apiVersion: v1
kind: Pod
metadata:
  labels:
    name: nginx
  name: nginx
  namespace: default
spec:
  imagePullSecrets:
  - name: harbor
  containers:
  - name: nginx
    image: nginx:1.19
    ports:
    - containerPort: 5678
      protocol: TCP
    resources:
      limits:
        cpu: "1"
        memory: 2Gi
      requests:
        cpu: 50m
        memory: 1Gi
    imagePullPolicy: Always
    volumeMounts: # 挂载到容器
      - mountPath: /etc/nginx/conf.d/
        name: proxy-config
      - name: vol
        mountPath: /tmp
      - name: nginx-html
        mountPath: /usr/share/nginx/html
  volumes: # 声明需要的卷
  - name: proxy-config
    configMap:
      defaultMode: 0644
      name: nginx-proxy-config
  - name: vol
    hostPath:
      path: /tmp
  - name: nginx-html
    persistentVolumeClaim:
      claimName: my-pvc
~~~

### 初始化容器
~~~bash
apiVersion: v1
kind: Pod
metadata:
  labels:
    name: nginx
  name: nginx
  namespace: default
spec:
  imagePullSecrets:
  - name: harbor
  containers:
  - name: nginx
    image: nginx:1.19
    ports:
    - containerPort: 5678
      protocol: TCP
    resources:
      limits:
        cpu: "1"
        memory: 2Gi
      requests:
        cpu: 50m
        memory: 1Gi
    imagePullPolicy: Always
    volumeMounts:
      - mountPath: /etc/nginx/conf.d/
        name: proxy-config
      - name: vol
        mountPath: /tmp
      - name: nginx-html
        mountPath: /usr/share/nginx/html
  initContainers: # 定义初始化容器
    - name: init-myservice # 等待nginx service的创建
      image: busybox:1.28
      command: ['sh', '-c', 'until nslookup nginx; do echo waiting for nginx; sleep 2; done;']
    - name: init-mydb
      image: busybox:1.28 # 检测 mydb service的存在
      command: ['sh', '-c', 'until nslookup mydb; do echo waiting for mydb; sleep 2; done;']
  volumes:
  - name: proxy-config
    configMap:
      defaultMode: 0644
      name: nginx-proxy-config
  - name: vol
    hostPath:
      path: /tmp
  - name: nginx-html
    persistentVolumeClaim:
      claimName: my-pvc

1.初始化容器为串行启动，前面结束后面开始
2.当初始化容器退出状态为错误时，kubelet 会根据Pod的restartPolicy策略进行重试
3.初始化容器也可以配置资源限制，但会和业务容器的分开计算，以较大的一方为pod的资源限制
4.初始化容器主要做:
  - 等待其他模块Ready,用来解决服务之间的依赖问题
  - 检测所有已经存在的成员节点，为主容器准备好集群的配置信息，这样主容器起来后就能用这个配置信息加入集群
  - 将节点注册到配置中心
~~~

### 容器生命周期回调(Hooks)
~~~bash
apiVersion: v1
kind: Pod
metadata:
  labels:
    name: nginx
  name: nginx
  namespace: default
spec:
  imagePullSecrets:
  - name: harbor
  containers:
  - name: nginx
    image: nginx:1.19
    ports:
    - containerPort: 5678
      protocol: TCP
    resources:
      limits:
        cpu: "1"
        memory: 2Gi
      requests:
        cpu: 50m
        memory: 1Gi
      lifecycle: # 生命周期
        postStart: # 启动钩子
          exec:
            command: ["/bin/sh", "-c", "echo Hello from the postStart handler > /usr/share/message"]
        preStop: # 终止钩子
          exec:
            command: ["/bin/sh","-c","nginx -s quit; while killall -0 nginx; do sleep 1; done"]
    imagePullPolicy: Always
    volumeMounts:
      - mountPath: /etc/nginx/conf.d/
        name: proxy-config
      - name: vol
        mountPath: /tmp
      - name: nginx-html
        mountPath: /usr/share/nginx/html
  initContainers:
    - name: init-myservice
      image: busybox:1.28
      command: ['sh', '-c', 'until nslookup nginx; do echo waiting for nginx; sleep 2; done;']
    - name: init-mydb
      image: busybox:1.28
      command: ['sh', '-c', 'until nslookup mydb; do echo waiting for mydb; sleep 2; done;']
  volumes:
  - name: proxy-config
    configMap:
      defaultMode: 0644
      name: nginx-proxy-config
  - name: vol
    hostPath:
      path: /tmp
  - name: nginx-html
    persistentVolumeClaim:
      claimName: my-pvc

1.生命周期钩子分为启动钩子(PostStart)和终止钩子(PreStop)，都是和主容器一起并行执行的
2.生命周期钩子的实现有两种:
  - exec  在容器的cgroups和名称空间中执行特定的命令(例如pre-stop.sh),命令所消耗的资源计入容器的资源消耗
  - http  对容器上的特定端点执行HTTP请求
3.生命周期钩子是属于容器的一部分，它的状态直接决定着容器的状态，容器结束时不会考虑PreStop的运行状态
~~~

### 健康检查
~~~bash
apiVersion: v1
kind: Pod
metadata:
  labels:
    name: nginx
  name: nginx
  namespace: default
spec:
  imagePullSecrets:
  - name: harbor
  containers:
  - name: nginx
    image: nginx:1.19
    ports:
    - containerPort: 5678
      protocol: TCP
    resources:
      limits:
        cpu: "1"
        memory: 2Gi
      requests:
        cpu: 50m
        memory: 1Gi
      lifecycle:
        postStart:
          exec:
            command: ["/bin/sh", "-c", "echo Hello from the postStart handler > /usr/share/message"]
        preStop:
          exec:
            command: ["/bin/sh","-c","nginx -s quit; while killall -0 nginx; do sleep 1; done"]
      startupProbe: # 启动探测
        exec: # 执行命令
          command:
          - cat
          - /tmp/healthy
        failureThreshold: 30 # 判断为失败次数
        periodSeconds: 10 # 间隔时间
        successThreshold: 1 # 判断为成功的次数
        timeoutSeconds: 1 # 超时时间
      livenessProbe: # 存活性探测
        tcpSocket: # tcp探测
          port: 8080 # 端口 8080
        failureThreshold: 3
        periodSeconds: 10
        successThreshold: 1
        timeoutSeconds: 1
      readinessProbe: # 就绪性探测
        httpGet: # http探测
          path: /adapter-weixin/s/about # 请求的路径
          port: 8080 # 端口
          scheme: HTTP # 请求的协议
        failureThreshold: 1
        periodSeconds: 6
        successThreshold: 3
        timeoutSeconds: 1
    restartPolicy: Always # 容器重启策略
    imagePullPolicy: Always
    volumeMounts:
      - mountPath: /etc/nginx/conf.d/
        name: proxy-config
      - name: vol
        mountPath: /tmp
      - name: nginx-html
        mountPath: /usr/share/nginx/html
  initContainers:
    - name: init-myservice
      image: busybox:1.28
      command: ['sh', '-c', 'until nslookup nginx; do echo waiting for nginx; sleep 2; done;']
    - name: init-mydb
      image: busybox:1.28
      command: ['sh', '-c', 'until nslookup mydb; do echo waiting for mydb; sleep 2; done;']
  volumes:
  - name: proxy-config
    configMap:
      defaultMode: 0644
      name: nginx-proxy-config
  - name: vol
    hostPath:
      path: /tmp
  - name: nginx-html
    persistentVolumeClaim:
      claimName: my-pvc

1.健康检查由三种探测器实现:
  - startupProbe    启动探测会阻塞存活性探测，其作用与存活性探测相同，使用启动探测的目的是保护满启动容器
  - livenessProbe   存活性探测依靠定义的检测方式定期检查容器状态，成功时容器状态为Running，失败后依据容器的重启策略对容器进行处理
  - readinessProbe  就绪性探测负责容器流量的接入，成功时容器状态Ready为true，失败时service从endpoints中移除pod ip，存活性探测失败也会停止流量接入
2.容器探测的三种方式:
  - HTTP  http请求 返回值200到400之间为正常
  - TCP   端口 端口存在为正常
  - Exec  命令执行 命令执行的返回值为0则正常
3.容器重启的三种策略:
  - Always     当容器终止退出后，总是重启容器，默认策略。
  - Onfailure  当容器异常退出（退出码非0）时，才重启容器。
  - Never      当容器终止退出时，才不重启容器。
~~~

## 三、pod与容器

### 定位容器的位置
~~~bash
...
status:
  conditions:
  - lastProbeTime: null
    lastTransitionTime: "2021-08-31T11:04:14Z"
    status: "True"
    type: Initialized
  - lastProbeTime: null
    lastTransitionTime: "2021-08-31T11:04:16Z"
    status: "True"
    type: Ready
  - lastProbeTime: null
    lastTransitionTime: "2021-08-31T11:04:16Z"
    status: "True"
    type: ContainersReady
  - lastProbeTime: null
    lastTransitionTime: "2021-08-31T11:04:14Z"
    status: "True"
    type: PodScheduled
  containerStatuses:
  - containerID: docker://7de54dfb03976d78ab3dda086c8cafbb6a28cbd276876ae60c815830763c10ff # 容器ID
    image: nginx:1.19
    imageID: docker-pullable://nginx@sha256:df13abe416e37eb3db4722840dd479b00ba193ac6606e7902331dcea50f4f1f2
    lastState: {}
    name: nginx
    ready: true
    restartCount: 0
    started: true
    state:
      running:
        startedAt: "2021-08-31T11:04:15Z"
  hostIP: 172.16.246.154 # 容器所在主机
  phase: Running
  podIP: 172.16.162.109
  podIPs:
  - ip: 172.16.162.109
  qosClass: BestEffort
  startTime: "2021-08-31T11:04:14Z"
~~~

### 查看容器配置
~~~bash
1.登陆指定机器
2.进入指定目录 cd /var/lib/docker/containers/{容器ID}
3.查看容器配置 cat config.v2.json | jq
~~~
