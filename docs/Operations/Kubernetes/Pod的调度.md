# Pod 的调度


## 概述

- 在默认情况下，kube-scheduler，会根据默认的调度策略将 pod 调度到节点上，这个时候我们不需要关心具体的节点，只需要关注资源池即可。但是有些时候我们需要将pod调度或者不调度至特定节点，则可以通过以下方式实现。

## 常用方式

### 一、NodeSelector

- 但在某些情况下，我们需要将 pod 调度至指定节点，这个时候我们可以使用 nodeSelector 来简单实现。

~~~yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    env: test
spec:
  containers:
  - name: nginx
    image: nginx
    imagePullPolicy: IfNotPresent
  nodeSelector: # 选择 node 标签
    disktype: ssd
~~~

- nodeSelector 通过 node 节点的 labels 来镜像关联，我们可以打标签或者默认标签的方式来独立出 node
- 常用内置的节点标签有：
  - kubernetes.io/hostname 
  - failure-domain.beta.kubernetes.io/zone
  - failure-domain.beta.kubernetes.io/region
  - topology.kubernetes.io/zone
  - topology.kubernetes.io/region
  - beta.kubernetes.io/instance-type
  - node.kubernetes.io/instance-type
  - kubernetes.io/os
  - kubernetes.io/arch

### 二、污点(Taint) 与 容忍度(Toleration)

- 污点(Taint) 是应用在节点之上的, 从这个名字就可以看出来, 是为了排斥 pod 所存在的.
- 容忍度(Toleration)是应用于 Pod 上的, 允许(但并不要求) Pod 调度到带有与之匹配的污点的节点上. 

#### 作用

- Taint(污点) 和 Toleration(容忍) 可以作用于 node 和 pod 上, 其目的是优化 pod 在集群间的调度,  具有 taint 的 node 和 pod 是互斥关系。
- Taint(污点) 和 toleration(容忍) 相互配合, 可以用来避免 pod 被分配到不合适的节点上. 每个节点上都可以应用一个或多个taint, 这表示对于那些不能容忍这些 taint 的 pod, 是不会被该节点接受的. 如果将 toleration 应用于 pod 上, 则表示这些 pod 可以(但不要求)被调度到具有相应 taint 的节点上.

#### 使用

- 添加和移除污点

```bash
1.添加
kubectl taint nodes node1 key1=value1:NoSchedule

2.移除（尾部的 - ）
kubectl taint nodes node1 key1=value1:NoSchedule-


这里的 key1=value1:NoSchedule 指 给 node 打上 一个污点标签 key1=value1 并且此标签的污点策略为 NoSchedule

污点策略：
NoSchedule： Kubernetes 不会将 Pod 分配到该节点。
PreferNoSchedule： Kubernetes 会尝试不将 Pod 分配到该节点。
NoExecute： Kubernetes 将 Pod 从该节点驱逐。
```

- 给 pod 设置容忍

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    env: test
spec:
  containers:
  - name: nginx
    image: nginx
    imagePullPolicy: IfNotPresent
  tolerations: # 设置容忍
  - key: "key1"
    operator: "Equal" # 这里使用 Equal 则下方必须写 key1 对应的值
    value: "value1"
    effect: "NoSchedule" # 污点策略
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    env: test
spec:
  containers:
  - name: nginx
    image: nginx
    imagePullPolicy: IfNotPresent
  tolerations: # 设置容忍
  - key: "key1"
    operator: "Exists" # 这里使用 Exists 则无需填写 key1 对应的值，存在 key1 标签即可
    effect: "NoSchedule"
```

- 存在两种特殊情况： 
  - 如果一个容忍度的 key 为空且 operator 为 Exists， 表示这个容忍度与任意的 key 、value 和 effect 都匹配，即这个容忍度能容忍任意 taint。 
  - 如果 effect 为空，则可以与所有键名 key1 的效果相匹配。

- 典型案例 
  - 专用节点：如果将某些节点专门分配给特定的一组用户使用，可以给这些节点添加一个 NoSchedule 污点，并在 pod 上配置亲和性即可
  - 配备了特殊硬件的节点：在部分节点配备了特殊硬件（比如 GPU）的集群中， 我们希望不需要这类硬件的 Pod 不要被分配到这些特殊节点，可以先给配备了特殊硬件的节点添加 NoSchedule/PreferNoSchedule 污点，并在 pod 上配置亲和性即可
  - 基于污点的驱逐: 这是在每个 Pod 中配置的在节点出现问题时的驱逐行为，接下来的章节会描述这个特性。

#### 基于污点的驱逐

当污点的 effect 值 NoExecute会影响已经在节点上运行的 Pod
- 如果 Pod 不能忍受 effect 值为 NoExecute 的污点，那么 Pod 将马上被驱逐 
- 如果 Pod 能够忍受 effect 值为 NoExecute 的污点，但是在容忍度定义中没有指定 tolerationSeconds，则 Pod 还会一直在这个节点上运行。 
- 如果 Pod 能够忍受 effect 值为 NoExecute 的污点，而且指定了 tolerationSeconds， 则 Pod 还能在这个节点上继续运行这个指定的时间长度。

当某种条件为真时，节点控制器会自动给节点添加一个污点。当前内置的污点包括：
- node.kubernetes.io/not-ready：节点未准备好。这相当于节点状态 Ready 的值为 "False"。 
- node.kubernetes.io/unreachable：节点控制器访问不到节点. 这相当于节点状态 Ready 的值为 "Unknown"。 
- node.kubernetes.io/memory-pressure：节点存在内存压力。 
- node.kubernetes.io/disk-pressure：节点存在磁盘压力。 
- node.kubernetes.io/pid-pressure: 节点的 PID 压力。 
- node.kubernetes.io/network-unavailable：节点网络不可用。 
- node.kubernetes.io/unschedulable: 节点不可调度。 
- node.cloudprovider.kubernetes.io/uninitialized：如果 kubelet 启动时指定了一个 "外部" 云平台驱动， 它将给当前节点添加一个污点将其标志为不可用。在 cloud-controller-manager 的一个控制器初始化这个节点后，kubelet 将删除这个污点。


### 三、亲和性与反亲和性

- 亲和性/反亲和性功能极大地扩展了你可以表达约束的类型。关键的增强点包括：
  - 语言表达能力更强（不仅仅是“对完全匹配规则的 AND”） 
  - 你可以发现规则是“软需求”/“偏好”，而不是硬性要求，因此， 如果调度器无法满足该要求，仍然调度该 Pod 
  - 你可以使用节点上（或其他拓扑域中）的 Pod 的标签来约束，而不是使用 节点本身的标签，来允许哪些 pod 可以或者不可以被放置在一起。

#### 分类

- 目前有两种类型的节点亲和性，分别为 
  - requiredDuringSchedulingIgnoredDuringExecution（硬亲和性），pod的调度必须满足亲和性的要求
  - preferredDuringSchedulingIgnoredDuringExecution（软亲和性），pod的调度尝试满足亲和性的要求
  
节点亲和性通过 PodSpec 的 affinity 字段下的 nodeAffinity 字段进行指定，例如：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: with-node-affinity
spec:
  affinity:
    nodeAffinity: # 配置亲和性
      requiredDuringSchedulingIgnoredDuringExecution: # 硬亲和性
        nodeSelectorTerms:
        - matchExpressions:
          - key: kubernetes.io/e2e-az-name
            operator: In
            values:
            - e2e-az1
            - e2e-az2
      preferredDuringSchedulingIgnoredDuringExecution: # 软亲和性
      - weight: 1
        preference:
          matchExpressions:
          - key: another-node-label-key
            operator: In
            values:
            - another-node-label-value
  containers:
  - name: with-node-affinity
    image: k8s.gcr.io/pause:2.0

这里同样也是对标签进行关联，而对标签值的匹配也就是 operator 有：
In # 标签值在下列值中即可
NotIn # 标签值不在下列值中即可
Exists # 标签值存在即可
DoesNotExist # 标签值不存在即可
Gt # 标签值大于下列值中即可
Lt # 标签值小于下列值中即可

# 亲和性与反亲和性是通过选择合适的上述操作符来实现

preferredDuringSchedulingIgnoredDuringExecution 中的 weight 字段值的 范围是 1-100。 pod 调度时会根据此权重对节点进行评分，调度至分最高的节点。
```

- 如果同时指定了 nodeSelector 和 nodeAffinity，两者必须都要满足， 才能将 Pod 调度到候选节点上
- 如果在 nodeAffinity 中指定多个 nodeSelectorTerms，满足其中一个 nodeSelectorTerms，就可以将 pod 调度到节点上
- 如果在 nodeSelectorTerms 中指定多个 matchExpressions，只有当所有 matchExpressions 满足的情况下，Pod 才会可以调度到节点上
- 如果你修改或删除了 pod 所调度到的节点的标签，Pod 不会被删除。 换句话说，亲和性选择只在 Pod 调度期间有效。

### 四、优先级和抢占

- Pod 可以有优先级。 优先级表示一个 Pod 相对于其他 Pod 的重要性。 
- 如果一个 Pod 无法被调度，调度程序会尝试抢占（驱逐）较低优先级的 Pod， 以使悬决 Pod 可以被调度。

#### PriorityClass

- 创建一个 PriorityClass 资源

```yaml
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass # 优先级类型资源
metadata:
  name: high-priority
value: 1000000 # 优先级
globalDefault: false # 是否为全局默认
description: "此优先级类应仅用于 test 服务 Pod。" # 注释
```

- PriorityClass 是一个无名称空间对象，优先级值在必填的 value 字段中指定。值越大，优先级越高。 
- PriorityClass 对象的名称必须是有效的 DNS 子域名， 并且它不能以 system- 为前缀。
- PriorityClass 对象可以设置任何小于或等于 10 亿的 32 位整数值。 较大的数字是为通常不应被抢占或驱逐的关键的系统 Pod 所保留的。 集群管理员应该为这类映射分别创建独立的 PriorityClass 对象。
- PriorityClass 还有两个可选字段：globalDefault 和 description。 globalDefault 字段表示这个 PriorityClass 的值应该用于没有 priorityClassName 的 Pod。 系统中只能存在一个 globalDefault 设置为 true 的 PriorityClass。 如果不存在设置了 globalDefault 的 PriorityClass， 则没有 priorityClassName 的 Pod 的优先级为零。
- description 字段是一个任意字符串。 它用来告诉集群用户何时应该使用此 PriorityClass。


#### 配置 pod 优先级

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    env: test
spec:
  containers:
  - name: nginx
    image: nginx
    imagePullPolicy: IfNotPresent
  priorityClassName: high-priority # 此处引用定义的优先级资源
```

### 五、节点保留资源

- Kubernetes 的节点可以按照节点的资源容量进行调度，默认情况下 Pod 能够使用节点全部可用容量。
- 这样就会造成一个问题，因为节点自己通常运行了不少驱动 OS 和 Kubernetes 的系统守护进程。
- 除非为这些系统守护进程留出资源，否则它们将与 Pod 争夺资源并导致节点资源短缺问题。 
- 当我们在线上使用 Kubernetes 集群的时候，如果没有对节点配置正确的资源预留，我们可以考虑一个场景，由于某个应用无限制的使用节点的 CPU 资源，导致节点上 CPU 使用持续100%运行，而且压榨到了 kubelet 组件的 CPU 使用，这样就会导致 kubelet 和 apiserver 的心跳出问题，节点就会出现 Not Ready 状况了。
- 默认情况下节点 Not Ready 过后，5分钟后会驱逐应用到其他节点，当这个应用跑到其他节点上的时候同样100%的使用 CPU，是不是也会把这个节点搞挂掉，同样的情况继续下去，也就导致了整个集群的雪崩，集群内的节点一个一个的 Not Ready 了，后果是非常严重的。
- 要解决这个问题就需要为 Kubernetes 集群配置资源预留，kubelet 暴露了一个名为 Node Allocatable 的特性，有助于为系统守护进程预留计算资源，Kubernetes 也是推荐集群管理员按照每个节点上的工作负载来配置 Node Allocatable。


#### Node Allocatable

- Kubernetes 节点上的 Allocatable 被定义为 Pod 可用计算资源量，调度器不会超额申请 Allocatable,目前支持 CPU, memory 和 ephemeral-storage 这几个参数。

我们可以通过 kubectl describe node 命令查看节点可分配资源的数据：

```yaml
kubectl describe node cn-hangzhou.172.16.246.144
...
Capacity:
  cpu:                4
  ephemeral-storage:  51474024Ki
  hugepages-1Gi:      0
  hugepages-2Mi:      0
  memory:             16266180Ki
  pods:               23
Allocatable:
  cpu:                4
  ephemeral-storage:  47438460440
  hugepages-1Gi:      0
  hugepages-2Mi:      0
  memory:             15242180Ki
  pods:               23
...
```

- 可以看到其中有 Capacity 与 Allocatable 两项内容，其中的 Allocatable 就是节点可被分配的资源，我们这里没有配置资源预留，所以默认情况下 Capacity 与 Allocatable 的值基本上是一致的。下图显示了可分配资源和资源预留之间的关系：


- Kubelet Node Allocatable 用来为 Kube 组件和 System 进程预留资源，从而保证当节点出现满负荷时也能保证 Kube 和 System 进程有足够的资源。
- 目前支持 cpu, memory, ephemeral-storage 三种资源预留。
- Node Capacity 是节点的所有硬件资源，kube-reserved 是给 kube 组件预留的资源，system-reserved 是给系统进程预留的资源，eviction-threshold 是 kubelet 驱逐的阈值设定，allocatable 才是真正调度器调度 Pod 时的参考值（保证节点上所有 Pods 的 request 资源不超过Allocatable）。


- 节点可分配资源的计算方式为： 
  - Node Allocatable Resource = Node Capacity - Kube-reserved - system-reserved - eviction-threshold

#### 配置资源预留

- 首先我们来配置 Kube 预留值，kube-reserved 是为了给诸如 kubelet、容器运行时、node problem detector 等 kubernetes 系统守护进程争取资源预留。要配置 Kube 预留，需要把 kubelet 的 --kube-reserved-cgroup 标志的值设置为 kube 守护进程的父控制组。
- 不过需要注意，如果 --kube-reserved-cgroup 不存在，Kubelet 不会创建它，启动 Kubelet 将会失败。


- 比如我们这里修改节点的 Kube 资源预留，我们可以直接修改 /var/lib/kubelet/config.yaml 文件来动态配置 kubelet，添加如下所示的资源预留配置：

```yaml
enforceNodeAllocatable:
- pods
- kube-reserved  # 开启 kube 资源预留
kubeReserved:
  cpu: 500m
  memory: 1Gi
  ephemeral-storage: 1Gi
kubeReservedCgroup: /kubelet.slice  # 指定 kube 资源预留的 cgroup
```

- 修改完成后，重启 kubelet，如果没有创建上面的 kubelet 的 cgroup，启动会失败：
```bash
systemctl restart kubelet
journalctl -u kubelet -f

Aug 11 15:04:13 ydzs-node4 kubelet[28843]: F0811 15:04:13.653476   28843 kubelet.go:1380] Failed to start ContainerManager Failed to enforce Kube Reserved Cgroup Limits on "/kubelet.slice": ["kubelet"] cgroup does not exist
```

- 上面的提示信息很明显，我们指定的 kubelet 这个 cgroup 不存在，但是由于子系统较多，具体是哪一个子系统不存在不好定位，我们可以将 kubelet 的日志级别调整为 v=4，就可以看到具体丢失的 cgroup 路径：

```bash
vi /var/lib/kubelet/kubeadm-flags.env
KUBELET_KUBEADM_ARGS="--v=4 --cgroup-driver=systemd --network-plugin=cni"

然后再次重启 kubelet：

systemctl daemon-reload
systemctl restart kubelet
```


- 再次查看 kubelet 日志：


```bash
journalctl -u kubelet -f

Sep 09 17:57:36 ydzs-node4 kubelet[20427]: I0909 17:57:36.382811   20427 cgroup_manager_linux.go:273] The Cgroup [kubelet] has some missing paths: [/sys/fs/cgroup/cpu,cpuacct/kubelet.slice /sys/fs/cgroup/memory/kubelet.slice /sys/fs/cgroup/systemd/kubelet.slice /sys/fs/cgroup/pids/kubelet.slice /sys/fs/cgroup/cpu,cpuacct/kubelet.slice /sys/fs/cgroup/cpuset/kubelet.slice]
Sep 09 17:57:36 ydzs-node4 kubelet[20427]: I0909 17:57:36.383002   20427 factory.go:170] Factory "systemd" can handle container "/system.slice/run-docker-netns-db100461211c.mount", but ignoring.
Sep 09 17:57:36 ydzs-node4 kubelet[20427]: I0909 17:57:36.383025   20427 manager.go:908] ignoring container "/system.slice/run-docker-netns-db100461211c.mount"
Sep 09 17:57:36 ydzs-node4 kubelet[20427]: F0909 17:57:36.383046   20427 kubelet.go:1381] Failed to start ContainerManager Failed to enforce Kube Reserved Cgroup Limits on "/kubelet.slice": ["kubelet"] cgroup does not exist
```

- 注意：systemd 的 cgroup 驱动对应的 cgroup 名称是以 .slice 结尾的，比如如果你把 cgroup 名称配置成 kubelet.service，那么对应的创建的 cgroup 名称应该为 kubelet.service.slice。如果你配置的是 cgroupfs 的驱动，则用配置的值即可。无论哪种方式，通过查看错误日志都是排查问题最好的方式。


- 现在可以看到具体的 cgroup 不存在的路径信息了：
```
The Cgroup [kubelet] has some missing paths: [/sys/fs/cgroup/cpu,cpuacct/kubelet.slice /sys/fs/cgroup/memory/kubelet.slice /sys/fs/cgroup/systemd/kubelet.slice /sys/fs/cgroup/pids/kubelet.slice /sys/fs/cgroup/cpu,cpuacct/kubelet.slice /sys/fs/cgroup/cpuset/kubelet.slice]

所以要解决这个问题也很简单，我们只需要创建上面的几个路径即可：

mkdir -p /sys/fs/cgroup/cpu,cpuacct/kubelet.slice
mkdir -p /sys/fs/cgroup/memory/kubelet.slice
mkdir -p /sys/fs/cgroup/systemd/kubelet.slice
mkdir -p /sys/fs/cgroup/pids/kubelet.slice
mkdir -p /sys/fs/cgroup/cpu,cpuacct/kubelet.slice
mkdir -p /sys/fs/cgroup/cpuset/kubelet.slice
mkdir -p /sys/fs/cgroup/hugetlb/kubelet.slice
```

- 创建完成后，再次重启：

```bash
systemctl restart kubelet
journalctl -u kubelet -f
```


- 启动完成后我们可以通过查看 cgroup 里面的限制信息校验是否配置成功，比如我们查看内存的限制信息：

```bash
cat /sys/fs/cgroup/memory/kubelet.slice/memory.limit_in_bytes
1073741824  # 1Gi
```

- 现在再次查看节点的信息：

```yaml
kubectl describe node ydzs-node4

Addresses:
  InternalIP:  10.151.30.59
  Hostname:    ydzs-node4
Capacity:
  cpu:                4
  ephemeral-storage:  17921Mi
  hugepages-2Mi:      0
  memory:             8008820Ki
  pods:               110
Allocatable:
  cpu:                3500m
  ephemeral-storage:  15838635595
  hugepages-2Mi:      0
  memory:             6857844Ki
  pods:               110
```

- 可以看到可以分配的 Allocatable 值就变成了 Kube 预留过后的值了，证明我们的 Kube 预留成功了。

#### 系统预留值

- 我们也可以用同样的方式为系统配置预留值，system-reserved 用于为诸如 sshd、udev 等系统守护进程争取资源预留，system-reserved 也应该为 kernel 预留 内存，因为目前 kernel 使用的内存并不记在 Kubernetes 的 pod 上。但是在执行 system-reserved 预留操作时请加倍小心，因为它可能导致节点上的关键系统服务 CPU 资源短缺或因为内存不足而被终止，所以如果不是自己非常清楚如何配置，可以不用配置系统预留值。

- 同样通过 kubelet 的参数 --system-reserved 配置系统预留值，但是也需要配置 --system-reserved-cgroup 参数为系统进程设置 cgroup。
请注意，如果 --system-reserved-cgroup 不存在，kubelet 不会创建它，kubelet 会启动失败。

#### 驱逐阈值

- 上面我们还提到可分配的资源还和 kubelet 驱逐的阈值有关。节点级别的内存压力将导致系统内存不足，这将影响到整个节点及其上运行的所有 Pod，节点可以暂时离线直到内存已经回收为止，我们可以通过配置 kubelet 驱逐阈值来防止系统内存不足。驱逐操作只支持内存和 ephemeral-storage 两种不可压缩资源。当出现内存不足时，调度器不会调度新的 Best-Effort QoS Pods 到此节点，当出现磁盘压力时，调度器不会调度任何新 Pods 到此节点。


- 我们这里为 ydzs-node4 节点配置如下所示的硬驱逐阈值：
```yaml
# /var/lib/kubelet/config.yaml

evictionHard:  # 配置硬驱逐阈值
  memory.available: "300Mi"
  nodefs.available: "10%"
enforceNodeAllocatable:
- pods
- kube-reserved
kubeReserved:
  cpu: 500m
  memory: 1Gi
  ephemeral-storage: 1Gi
kubeReservedCgroup: /kubelet.slice
```


- 我们通过 --eviction-hard 预留一些内存后，当节点上的可用内存降至保留值以下时，kubelet 将尝试驱逐 Pod，
```yaml
kubectl describe node ydzs-node4

Addresses:
  InternalIP:  10.151.30.59
  Hostname:    ydzs-node4
Capacity:
  cpu:                4
  ephemeral-storage:  17921Mi
  hugepages-2Mi:      0
  memory:             8008820Ki
  pods:               110
Allocatable:
  cpu:                3500m
  ephemeral-storage:  15838635595
  hugepages-2Mi:      0
  memory:             6653044Ki
  pods:               110
```

- 配置生效后再次查看节点可分配的资源可以看到内存减少了，临时存储没有变化是因为硬驱逐的默认值就是 10%。也是符合可分配资源的计算公式的：
```bash
Node Allocatable Resource = Node Capacity - Kube-reserved - system-reserved - eviction-threshold
```

## 调度流程

- API Server接受客户端提交 Pod 对象创建请求后的操作过程中，有一个重要的步骤就是由调度器程序 kube-scheduler 从当前集群中选择一个可用的最佳节点来接收并运行它，通常是默认的调度器 kube-scheduler 负责执行此类任务。

- 对于每个待创建的Pod对象来说，调度过程通常分为两个阶段: 过滤 —> 打分，过滤阶段用来过滤掉不符合调度规则的Node，打分阶段建立在过滤阶段之上，为每个符合调度的Node进行打分，分值越高则被调度到该Node的机率越大。


### kube-scheduler

[官方文档](https://kubernetes.io/zh/docs/concepts/scheduling-eviction/kube-scheduler/)

#### kube-scheduler 调度介绍

- kube-scheduler 是 Kubernetes 集群的默认调度器，并且是集群控制面 master 的一部分。对每一个新创建的Pod或者是未被调度的 Pod，kube-scheduler 会选择一个最优的 Node 去运行 Pod。

- 然而，Pod 内的每一个容器对资源都有不同的需求，而且 Pod 本身也有不同的资源需求。因此，Pod 在被调度到 Node 上之前，根据这些特定的资源调度需求，需要对集群中的 Node 进行一次过滤。

- 在一个集群中，满足一个 Pod 调度请求的所有 Node 称之为可调度节点。如果没有任何一个 Node 能满足 Pod 的资源请求，那么这个 Pod 将一直停留在未调度状态直到调度器能够找到合适的 Node。


#### kube-scheduler 调度流程

- kube-scheduler 给一个 pod 做调度选择包含两个步骤： 
  - 1.过滤（Predicates 预选策略）
  - 2.打分（Priorities 优选策略）


- 过滤阶段：过滤阶段会将所有满足 Pod 调度需求的 Node 选出来。例如，PodFitsResources 过滤函数会检查候选 Node 的可用资源能否满足 Pod 的资源请求。在过滤之后，得出一个 Node 列表，里面包含了所有可调度节点；通常情况下，这个 Node 列表包含不止一个 Node。如果这个列表是空的，代表这个 Pod 不可调度。

- 打分阶段：在过滤阶段后调度器会为 Pod 从所有可调度节点中选取一个最合适的 Node。根据当前启用的打分规则，调度器会给每一个可调度节点进行打分。最后，kube-scheduler 会将 Pod 调度到得分最高的 Node 上。如果存在多个得分最高的 Node，kube-scheduler 会从中随机选取一个。


#### 过滤阶段

[官方文档](https://kubernetes.io/docs/reference/scheduling/policies/)

```bash
1.PodFitsHostPorts：检查Node上是否不存在当前被调度Pod的端口（如果被调度Pod用的端口已被占用，则此Node被Pass）。

2.PodFitsHost：检查Pod是否通过主机名指定了特性的Node (是否在Pod中定义了nodeName)

3.PodFitsResources：检查Node是否有空闲资源(如CPU和内存)以满足Pod的需求。

4.PodMatchNodeSelector：检查Pod是否通过节点选择器选择了特定的Node (是否在Pod中定义了nodeSelector)。

5.NoVolumeZoneConflict：检查Pod请求的卷在Node上是否可用 (不可用的Node被Pass)。

6.NoDiskConflict：根据Pod请求的卷和已挂载的卷，检查Pod是否合适于某个Node (例如Pod要挂载/data到容器中，Node上/data/已经被其它Pod挂载，那么此Pod则不适合此Node)

7.MaxCSIVolumeCount：：决定应该附加多少CSI卷，以及是否超过了配置的限制。

8.CheckNodeMemoryPressure：对于内存有压力的Node，则不会被调度Pod。

9.CheckNodePIDPressure：对于进程ID不足的Node，则不会调度Pod

10.CheckNodeDiskPressure：对于磁盘存储已满或者接近满的Node，则不会调度Pod。

11.CheckNodeCondition：Node报告给API Server说自己文件系统不足，网络有写问题或者kubelet还没有准备好运行Pods等问题，则不会调度Pod。

12.PodToleratesNodeTaints：检查Pod的容忍度是否能承受被打上污点的Node。

13.CheckVolumeBinding：根据一个Pod并发流量来评估它是否合适（这适用于结合型和非结合型PVCs）。
```

#### 打分阶段

[官方文档](https://kubernetes.io/docs/reference/scheduling/policies/)

- 当过滤阶段执行后满足过滤条件的Node，将进行打分阶段。

```bash
1.SelectorSpreadPriority：优先减少节点上属于同一个 Service 或 Replication Controller 的 Pod 数量

2.InterPodAffinityPriority：优先将 Pod 调度到相同的拓扑上（如同一个节点、Rack、Zone 等）

3.LeastRequestedPriority：节点上放置的Pod越多，这些Pod使用的资源越多，这个Node给出的打分就越低，所以优先调度到Pod少及资源使用少的节点上。

4.MostRequestedPriority：尽量调度到已经使用过的 Node 上，将把计划的Pods放到运行整个工作负载所需的最小节点数量上。

5.RequestedToCapacityRatioPriority：使用默认资源评分函数形状创建基于requestedToCapacity的
ResourceAllocationPriority。

6.BalancedResourceAllocation：优先平衡各节点的资源使用。

7.NodePreferAvoidPodsPriority：根据节点注释对节点进行优先级排序，以使用它来提示两个不同的 Pod 不应在同一节点上运行。
scheduler.alpha.kubernetes.io/preferAvoidPods。

8.NodeAffinityPriority：优先调度到匹配 NodeAffinity （Node亲和性调度）的节点上。

9.TaintTolerationPriority：优先调度到匹配 TaintToleration (污点) 的节点上

10.ImageLocalityPriority：尽量将使用大镜像的容器调度到已经下拉了该镜像的节点上。

11.ServiceSpreadingPriority：尽量将同一个 service 的 Pod 分布到不同节点上，服务对单个节点故障更具弹性。

12.EqualPriority：将所有节点的权重设置为 1。

13.EvenPodsSpreadPriority：实现首选pod拓扑扩展约束。
```

- 得分最高的 node 节点将与 pod 绑定