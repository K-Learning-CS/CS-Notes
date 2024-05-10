# Containerd

## 一、简介

- 在 kubernetes v1.20 版本后移除了对 docker 的支持，迫使用户转向 containerd。对于 k8s 来讲，docker 并不支持 CRI 标准，于是得维护一个名为 dockershim 的插件来进行对 docker 的调用，而 containerd 是在 docker 的妥协下开发出的符合 CRI 标准的组件，k8s 可以直接调用。去掉 docker，缩短了调用链路，极大的提升了效率，同时对于开发者来讲，dockershim 的弃用也减少了一部分维护的压力。

## 二、变化

- 对于用户来讲，这种转变几乎是无感的，用于 OCI 标准的存在，你使用 docker 构建的镜像可以使用 containerd 无缝运行，用户并不需要改变构建的习惯，但是当我们需要操作 containerd 时，需要转向使用 ctr CMD 来操作 containerd 而不是 docker ，其中的操作大部分都类似，没有太多学习成本。

## 三、操作

### ctr 命令概览

```bash
[root@iZbp1d9vn80jvwkpw2j7naZ ~]# ctr
NAME:
   ctr -
        __
  _____/ /______
 / ___/ __/ ___/
/ /__/ /_/ /
\___/\__/_/

containerd CLI


USAGE:
   ctr [global options] command [command options] [arguments...]

VERSION:
   1.4.8

DESCRIPTION:

ctr is an unsupported debug and administrative client for interacting
with the containerd daemon. Because it is unsupported, the commands,
options, and operations are not guaranteed to be backward compatible or
stable from release to release of the containerd project.

COMMANDS:
   plugins, plugin            provides information about containerd plugins
   version                    print the client and server versions
   containers, c, container   manage containers
   content                    manage content
   events, event              display containerd events
   images, image, i           manage images
   leases                     manage leases
   namespaces, namespace, ns  manage namespaces
   pprof                      provide golang pprof outputs for containerd
   run                        run a container
   snapshots, snapshot        manage snapshots
   tasks, t, task             manage tasks
   install                    install a new package
   oci                        OCI tools
   shim                       interact with a shim directly
   help, h                    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug                      enable debug output in logs
   --address value, -a value    address for containerd's GRPC server (default: "/run/containerd/containerd.sock") [$CONTAINERD_ADDRESS]
   --timeout value              total timeout for ctr commands (default: 0s)
   --connect-timeout value      timeout for connecting to containerd (default: 0s)
   --namespace value, -n value  namespace to use with commands (default: "default") [$CONTAINERD_NAMESPACE]
   --help, -h                   show help
   --version, -v                print the version
```

### 镜像操作

```
[root@iZbp1d9vn80jvwkpw2j7naZ ~]# ctr images
NAME:
   ctr images - manage images

USAGE:
   ctr images command [command options] [arguments...]

COMMANDS:
   check       check that an image has all content available locally
   export      export images
   import      import images
   list, ls    list images known to containerd
   mount       mount an image to a target path
   unmount     unmount the image from the target
   pull        pull an image from a remote
   push        push an image to a remote
   remove, rm  remove one or more images by reference
   tag         tag an image
   label       set and clear labels for an image

OPTIONS:
   --help, -h  show help
```

- 在containerd 中，多了一个namespace 的概念，所有的资源都需要指定 ns 才能查看。

```
查看ns
ctr ns ls
```

```
1.查看镜像
ctr -n k8s.io images ls

2.拉取镜像
ctr -n k8s.io images pull docker.io/library/nginx:alpine
--user # 指定用户密码

3.打标签
ctr -n k8s.io images tag docker.io/library/nginx:alpine harbor.k8s.local/course/nginx:alpine

4.推送镜像
ctr -n k8s.io images push harbor.k8s.local/course/nginx:alpine

5.删除镜像
ctr -n k8s.io images rm harbor.k8s.local/course/nginx:alpine
```

### 容器操作

```
[root@iZbp1d9vn80jvwkpw2j7naZ ~]# ctr container
NAME:
   ctr containers - manage containers

USAGE:
   ctr containers command [command options] [arguments...]

COMMANDS:
   create           create container
   delete, del, rm  delete one or more existing containers
   info             get info about a container
   list, ls         list containers
   label            set and clear labels for a container
   checkpoint       checkpoint a container
   restore          restore a container from checkpoint

OPTIONS:
   --help, -h  show help
```

```
1.创建容器
ctr -n k8s.io containers create docker.io/library/nginx:alpine nginx

2.查看容器
ctr -n k8s.io containers ls

3.容器详情
ctr -n k8s.io containers info nginx

4.删除容器
ctr -n k8s.io containers rm nginx
```

### 容器运行

```
[root@iZbp1d9vn80jvwkpw2j7naZ ~]# ctr -n k8s.io task
NAME:
   ctr tasks - manage tasks

USAGE:
   ctr tasks command [command options] [arguments...]

COMMANDS:
   attach           attach to the IO of a running container
   checkpoint       checkpoint a container
   delete, rm       delete one or more tasks
   exec             execute additional processes in an existing container
   list, ls         list tasks
   kill             signal a container (default: SIGTERM)
   pause            pause an existing container
   ps               list processes for container
   resume           resume a paused container
   start            start a container that has been created
   metrics, metric  get a single data point of metrics for a task with the built-in Linux runtime

OPTIONS:
   --help, -h  show help
```

```
1.启动容器
ctr -n k8s.io task start -d nginx

2.查看运行中的容器
ctr -n k8s.io task ls

3.进入容器
ctr -n k8s.io task exec --exec-id 0 -t nginx sh
--exec-id # 必须且唯一

4.停止容器
ctr -n k8s.io task kill nginx

5.删除容器
ctr -n k8s.io task rm nginx # 必须先kill

6.获取容器的内存、CPU 和 PID 的限额与使用量
ctr -n k8s.io task metrics nginx
```
