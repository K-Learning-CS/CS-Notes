# Volume

## 一、k8s中的存储类型

- volume：持久化存储卷，可以对数据进行持久化存储。
- 在 k8s 中，pod 的生命周期是短暂的，在没有使用存储的情况下，pod 产生的数据会随着 pod 的重启或销毁释放。


### 1）查看k8s支持哪些存储

~~~bash 
# kubectl explain pods.spec.volumes
 
awsElasticBlockStore
azureDisk
azureFile
cephfs
cinder
configMap
csi
downwardAPI
emptyDir
fc
flexVolume
flocker
gcePersistentDisk
gitRepo
glusterfs
hostPath
iscsi
nfs
persistentVolumeClaim
photonPersistentDisk
portworxVolume
projected
quobyte
rbd
scaleIO
secret
storageos
vsphereVolume
 
# 常用的如下：
emptyDir
hostPath
nfs
persistentVolumeClaim
glusterfs
cephfs
configMap
~~~

### 2）如何使用存储卷：
~~~bash
1.定义pod的volume，这个volume指明它要关联到哪个存储上的
2.在容器中要使用volume mounts（挂载存储）
~~~
## 二、常用的存储类型 
 
### 1）EmptyDir
- emptyDir类型的Volume在Pod分配到Node上时被创建，Kubernetes会在Node上自动分配一个目录，因此无需指定宿主机Node上对应的目录文件。 
- 这个目录的初始内容为空，当Pod从Node上移除时，emptyDir中的数据会被永久删除。
- emptyDir Volume主要用于某些应用程序无需永久保存的临时目录，多个容器的共享目录等。

~~~bash
# 例如我们在 k8s 的使用中使用的日志处理 sidecar 便是使用 EmptyDir 关联主容器，并共享日志目录的
...
        - name: clear-log
          image: harbor.qianfan123.com/toolset/logging-clean:no-prod
          resources:
            limits:
              cpu: 50m
              memory: 20Mi
            requests:
              cpu: 50m
              memory: 20Mi
          volumeMounts:
          - name: logs-storage
            mountPath: /opt/heading/tomcat/logs
      volumes:
      - name: config
        secret:
          secretName: filebeat-secret
      - name: logs-storage
        emptyDir: {}
...
~~~

### 2）HostPath
- hostPath Volume为Pod挂载宿主机上的目录或文件。 
- hostPath Volume使得容器可以使用宿主机的高速文件系统进行存储；
- hostpath（宿主机路径）：节点级别的存储卷，在pod被删除，这个存储卷还是存在的，不会被删除，所以只要同一个pod被调度到同一个节点上来，在pod被删除重新被调度到这个节点之后，对应的数据依然是存在的。
 
查看hostpath存储卷的使用
~~~bash
# kubectl explain pods.spec.volumes.hostPath

KIND:     Pod
VERSION:  v1

RESOURCE: hostPath <Object>

DESCRIPTION:
     HostPath represents a pre-existing file or directory on the host machine
     that is directly exposed to the container. This is generally used for
     system agents or other privileged things that are allowed to see the host
     machine. Most containers will NOT need this. More info:
     https://kubernetes.io/docs/concepts/storage/volumes#hostpath

     Represents a host path mapped into a pod. Host path volumes do not support
     ownership management or SELinux relabeling.

FIELDS:
   path	<string> -required-
     Path of the directory on the host. If the path is a symlink, it will follow
     the link to the real path. More info:
     https://kubernetes.io/docs/concepts/storage/volumes#hostpath

   type	<string>
     Type for HostPath Volume Defaults to "" More info:
     https://kubernetes.io/docs/concepts/storage/volumes#hostpath
~~~

~~~bash
cat pod_volume_host.yaml
 
apiVersion: v1
kind: Pod
metadata:
  name: test-hostpath
spec:
  containers:
  - image: nginx
    name: test-nginx
    volumeMounts:
    - mountPath: /test-nginx
      name: test-volume
  - image: tomcat
    name: test-tomcat
    volumeMounts:
    - mountPath: /test-tomcat
      name: test-volume
  volumes:
  - name: test-volume
    hostPath:
      path: /data1
      type: DirectoryOrCreate
 
kubectl apply -f pod_volume_host.yaml 
如果pod处于running状态，那么可以执行如下步骤测试存储卷是否可以被正常使用

- kubectl exec -it test-hostpath -c test-nginx -- /bin/bash 登录到nginx容器 
  查看否存在目录 /test-nginx/，如果存在，说明存储卷挂载成功
  
- kubectl exec -it test-hostpath -c test-tomcat-- /bin/bash 登录到tomcat容器
  查看是否存在目录 /test-tomcat/，如果存在，说明存储卷挂载成功
 
- kubectl get pods -o wide
  连接到 pod 调度的节点上，查看节点的 /data1 目录，如果存在 pod 中的文件则验证成功
~~~
 
- hostpath存储卷缺点： 单节点 pod删除之后重新创建必须调度到同一个node节点，数据才不会丢失
 
### 3）NFS
~~~bash
#搭建nfs

1.以master节点作为nfs服务端：
安装nfs：
yum install nfs-utils -y

2.创建共享目录：
mkdir /data/volumes -pv 
cat /etc/exports
/data/volumes 192.168.0.0/24(rw,no_root_squash) 
 
3.启动nfs
exportfs -arv
systemctl start nfs

4.测试
1.在node上手动挂载：
mount -t nfs 192.168.0.0:/data/volumes /mnt
df -h 可以看到已经挂载了

2.手动卸载：
umount /mnt


# 使用pod挂载nfs

1.创建pod
cat pod-nfs.yaml
 
apiVersion: v1
kind: Pod
metadata:
 name: test-nfs-volume
spec:
 containers:
 - name: test-nfs
   image: nginx
   ports:
   - containerPort: 80
     protocol: TCP
   volumeMounts:
   - name: nfs-volumes
     mountPath: /usr/share/nginx/html
 volumes:
 - name: nfs-volumes
   nfs:
    path: /data/volumes
    server: 192.168.0.6 # 此地址为 master 节点 ip
 
2.测试
在 master 节点的 nfs 挂载目录中创建 index.html ，输入任意内容，访问 nginx 显示内容则成功
 
 
- nfs支持多个客户端挂载，可以在多创建几个pod，挂载同一个nfs服务器
- 但是nfs如果宕机了，数据也就丢失了，所以需要使用分布式存储，常见的分布式存储有glusterfs和cephfs
~~~ 
 
 
### 4）PV、PVC

- PersistentVolume（PV）是 k8s 中的集群资源，他的作用是将外部非标准储存，描述成为集群内部的标准存储，方便统一管理。比如一块硬盘或者网络存储等，在集群中都可以描述成 PV
- PersistentVolumeClaim（PVC）是 k8s 中的一个持久化存储卷，仅在需要使用的时候定义，比如我要多大的磁盘，这个磁盘的读写模式是什么，当我们声明了 PVC 后，集群中如果有满足 PVC 需求的 PV，则会自动绑定上。
- 这里的绑定类似于我们使用 mount 把当前目录（PVC）挂载到磁盘（PV），只是这里绑定的不是目录，而是 k8s 中的资源对象

#### 1. pv和pvc的生命周期

~~~bash
（1）pv的供应方式
- 静态的
  集群管理员创建了许多PV。它们包含可供群集用户使用的实际存储的详细信息。它们存在于Kubernetes API中，可供使用。
- 动态的
  当管理员创建的静态PV都不匹配用户的PersistentVolumeClaim时，群集可能会尝试为PVC专门动态配置卷。
  此配置基于StorageClasses：PVC必须请求存储类，管理员必须已创建并配置该类，以便进行动态配置。

（2）绑定
- 用户创建pvc并指定需要的资源和访问模式。在找到可用pv之前，pvc会保持未绑定状态。

（3）使用
- 需要找一个存储服务器，把它划分成多个存储空间；
- k8s管理员可以把这些存储空间定义成多个pv；
- 在pod中使用pvc类型的存储卷之前需要先创建pvc，通过定义需要使用的pv的大小和对应的访问模式，找到合适的pv；
- pvc被创建之后，就可以当成存储卷来使用了，我们在定义pod时就可以使用这个pvc的存储卷
- pvc和pv它们是一一对应的关系，pv如果被被pvc绑定了，就不能被其他pvc使用了；
- 我们在创建pvc的时候，应该确保和底下的pv能绑定，如果没有合适的pv，那么pvc就会处于pending状态。

（4）回收策略
- 当我们创建pod时如果使用pvc做为存储卷，那么它会和pv绑定。
- 当删除pod，pvc和pv绑定就会解除，解除之后和pvc绑定的pv卷里的数据需要怎么处理，目前，卷可以保留，回收或删除。
- 他们的回收策略分别为：
  Retain
    当删除pvc的时候，pv仍然存在，处于released状态，但是它不能被其他pvc绑定使用，里面的数据还是存在的，当我们下次再使用的时候，数据还是存在的，这个是默认的回收策略。
    管理员能够通过下面的步骤手工回收存储卷：
    1）删除PV：在PV被删除后，在外部设施中相关的存储资产仍然还在；
    2）手工删除遗留在外部存储中的数据；
    3）手工删除存储资产，如果需要重用这些存储资产，则需要创建新的PV。
  Delete
    删除pvc时即会从Kubernetes中移除PV，也会从相关的外部设施中删除存储资产。
~~~
#### 2. 创建PV
~~~bash
（1）在nfs中导出多个存储目录，在nfs服务器上操作（这里是k8s的master节点）
mkdir -p /data/volume_test/v{1,2,3,4,5,6,7,8,9,10}
cat /etc/exports
/data/volume_test/v1 192.168.0.0/24(rw,no_root_squash)
/data/volume_test/v2 192.168.0.0/24(rw,no_root_squash)
/data/volume_test/v3 192.168.0.0/24(rw,no_root_squash)
/data/volume_test/v4 192.168.0.0/24(rw,no_root_squash)
exportfs -arv 使配置文件生效
service nfs restart

（2）把上面的存储目录做成pv
kubectl explain pv 查看pv的创建方法
kubectl explain pv.spec.nfs 查看怎么把nfs定义成pv
参考:https://kubernetes.io/docs/concepts/storage/persistent-volumes

（3）创建pv（pv是集群级别的资源，不需要定义namespace）
# cat pv.yaml 
apiVersion: v1
kind: PersistentVolume
metadata:
  name:  v1
spec:
  capacity:
    storage: 1Gi
  accessModes: ["ReadWriteOnce"]
  nfs:
    path: /data/volume_test/v1
    server: master1
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name:  v2
spec:
  capacity:
      storage: 2Gi
  accessModes: ["ReadWriteMany"]
  nfs:
    path: /data/volume_test/v2
    server: master1
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name:  v3
spec:
  capacity:
      storage: 3Gi
  accessModes: ["ReadOnlyMany"]
  nfs:
    path: /data/volume_test/v3
    server: master1
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name:  v4
spec:
  capacity:
      storage: 4Gi
  accessModes: ["ReadWriteOnce","ReadWriteMany"]
  nfs:
    path: /data/volume_test/v4
    server: master1
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name:  v5
spec:
  capacity: #pv的存储空间容量
      storage: 5Gi
  accessModes: ["ReadWriteOnce","ReadWriteMany"]
  nfs:
    path: /data/volume_test/v5 #把nfs的存储空间创建成pv
    server: master1

kubectl  apply  -f pv.yaml
kubectl get pv

（4）创建pvc
cat pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc
spec:
  accessModes: ["ReadWriteMany"]
  resources:
    requests:
      storage: 2Gi

（5）创建pod
cat pod-pvc.yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-pvc
spec:
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: nginx-html
      mountPath: /usr/share/nginx/html
  volumes:
  - name: nginx-html
    persistentVolumeClaim:
      claimName: my-pvc

kubectl  apply -f pod-pvc.yaml
 
注：
（1）我们每次创建pvc的时候，需要事先有划分好的pv，可能不方便，那么可以在创建pvc的时候直接动态创建一个pv这个存储类，pv事先是不存在的
（2）pvc和pv绑定，如果使用默认的回收策略retain，那么删除pvc之后，pv会处于released状态，我们想要继续使用这个pv，需要手动删除pv，kubectl delete pv pv_name，删除pv，不会删除pv里的数据，当我们重新创建pvc时还会和这个最匹配的pv绑定，数据还是原来数据，不会丢失.
~~~

### 5）StorageClass（动态存储）
- storageclass 是一个存储类，动态存储可以根据 PVC 的请求动态创建 PV ，让其完美绑定。
- StorageClass对象的名称很重要，用户定义的 PVC 描述中，通过动态存储的名称来请求创建 PV。
- 管理员在首次创建StorageClass对象时设置类的名称和其他参数，并且在创建对象后无法更新这些对象。

#### 1.创建动态存储（以 nfs 为例）
- 在使用 nfs 时，我们需要借助 nfs-provisioner 插件来实现对 nfs 资源的调用，所以需要启用一个 pod 来做这件事情
- 同样的，这个 pod 需要一些集群权限才能查看到创建请求，所以需要对其进行 RBAC 授权
~~~bash
（1）RBAC 授权
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-provisioner
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: nfs-provisioner-runner
rules:
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "update", "patch"]
  - apiGroups: [""]
    resources: ["services", "endpoints"]
    verbs: ["get"]
  - apiGroups: ["extensions"]
    resources: ["podsecuritypolicies"]
    resourceNames: ["nfs-provisioner"]
    verbs: ["use"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: run-nfs-provisioner
subjects:
  - kind: ServiceAccount
    name: nfs-provisioner
    namespace: default
roleRef:
  kind: ClusterRole
  name: nfs-provisioner-runner
  apiGroup: rbac.authorization.k8s.io
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-provisioner
rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-provisioner
subjects:
  - kind: ServiceAccount
    name: nfs-provisioner
    namespace: default
roleRef:
  kind: Role
  name: leader-locking-nfs-provisioner
  apiGroup: rbac.authorization.k8s.io

（2）创建 nfs-provisioner
kind: Deployment
apiVersion: apps/v1
metadata:
  name: nfs-provisioner
spec:
  selector:
    matchLabels:
       app: nfs-provisioner
  replicas: 1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: nfs-provisioner
    spec:
      serviceAccount: nfs-provisioner
      containers:
        - name: nfs-provisioner
          image: groundhog2k/nfs-subdir-external-provisioner:v3.2.0
          volumeMounts:
            - name: nfs-client-root
              mountPath: /persistentvolumes
          env:
            - name: PROVISIONER_NAME
              value: example.com/nfs
            - name: NFS_SERVER
              value: 172.16.1.81
            - name: NFS_PATH
              value: /data/volume_test/sc
      volumes:
        - name: nfs-client-root
          nfs:
            server: 172.16.1.81
            path: /data/volume_test/sc
---
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: nfs
provisioner: example.com/nfs
~~~

#### 2.默认动态存储
- 在默认情况下，集群中没有默认的动态存储，这时候创建 PVC 时需要指定动态存储才能创建，当我们设置默认动态存储后，则不再需要指定。
~~~bash
kubectl get storageclass
 
NAME                 PROVISIONER               AGE
standard (default)   kubernetes.io/gce-pd      1d
gold                 kubernetes.io/gce-pd      1d
 
# 设定默认storageclass
kubectl patch storageclass <your-class-name> -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'

# 取消默认storageclass
kubectl patch storageclass standard -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}'
~~~