# argo

目录
=================

* [argo](#argo)
   * [argo-rollouts](#argo-rollouts)
      * [一、安装](#一安装)
      * [二、命令使用](#二命令使用)
      * [三、灰度发布](#三灰度发布)
   * [argo cd](#argo-cd)
      * [一、安装](#一安装-1)
      * [二、命令使用](#二命令使用-1)
      * [三、使用argo cd](#三使用argo-cd)
         * [配置git仓库](#配置git仓库)
         * [创建应用](#创建应用)
         * [同步应用](#同步应用)
         * [创建项目](#创建项目)
         * [自定义资源](#自定义资源)
      * [用户及权限管理](#用户及权限管理)
         * [用户及kustomize版本参数配置](#用户及kustomize版本参数配置)
         * [用户权限管理](#用户权限管理)
      * [argocd web页面访问及代理](#argocd-web页面访问及代理)
         * [设计](#设计)
         * [配置](#配置)
## argo-rollouts
[官方网站](https://argoproj.github.io/argo-rollouts/)

### 一、安装
~~~bash
# argo rollout为定制资源 所以在使用前需要在目标集群中安装才能使用
1.定制资源及提供者安装
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
# 高可用：将副本数改为3

2.cli安装
wget http://download.init0.cn/kubectl-argo-rollouts
chmod +x kubectl-argo-rollouts
mv kubectl-argo-rollouts /usr/local/bin
~~~

### 二、命令使用
`  <> 中的部分需要依据实际情况替换 `

`  详情见 kubectl argo rollouts --help `

`  在引入argo rollout时主要考虑的是版本升级时平滑过度的问题 所以基本不会涉及到单独使用argo rollout的情况 `

- 版本检查

~~~bash
kubectl argo rollouts version
~~~

- 开启图形界面

~~~bash
kubectl argo rollouts dashboard
~~~

- 创建

~~~bash
kubectl argo rollouts create -f <rollout_resources>
~~~

- 列出已有rollout

~~~bash
kubectl argo rollouts list rollouts
~~~

- 获取rollout状态

~~~bash
kubectl argo rollouts get rollout <rollout_name>
~~~

- 更新镜像

~~~bash
kubectl argo rollouts set image <rollout_name> <name:version>
~~~

- 暂停更新

~~~bash
kubectl argo rollouts pause <rollout_name>
~~~

- 推动暂停的rollout更新

~~~bash
kubectl argo rollouts promote <rollout_name>
kubectl argo rollouts promote <rollout_name> --full #忽略更新策略全面更新
~~~

- 重启rollout部署的pod

~~~bash
kubectl argo rollouts restart <rollout_name>
~~~

- 重新发布

~~~bash
kubectl argo rollouts retry rollout <rollout_name>
~~~

- 回滚

~~~bash
kubectl argo rollouts undo <rollout_name>
kubectl argo rollouts undo <rollout_name> --to-revision=3 # 回滚至指定版本
~~~


### 三、灰度发布
~~~bash
cat > base/deployment.yaml << 'EOF'
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: nginx
spec:
  replicas: 1
  strategy:
    canary:
      steps:
      - setWeight: 20
      - pause: {duration: 10}
  selector:
    matchLabels:
      name: nginx
  template:
    metadata:
      labels:
        name: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.18
          ports:
            - containerPort: 80
              protocol: TCP
          volumeMounts:
            - mountPath: /usr/share/nginx/html/
              name: config
      volumes:
      - name: config
        configMap:
          defaultMode: 0644
          name: nginx-index
EOF

cat > overlays/test/deployment.yaml << 'EOF'
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: nginx
spec:
  replicas: 5
  strategy:
    canary:
      steps:
      - setWeight: 20
      - pause: {duration: 10}
      - setWeight: 40
      - pause: {duration: 10}
      - setWeight: 60
      - pause: {duration: 10}
      - setWeight: 80
      - pause: {duration: 10}
EOF
~~~

- 部署

~~~bash
kustomize build overlays/test/ | kubectl apply -f -
~~~

~~~bash
kubectl get rollouts -n kustomize-kang
~~~

~~~bash
kubectl argo rollouts get rollout test-nginx-suffix -n kustomize-kang -w
~~~

~~~bash
kustomize build overlays/test/
~~~

- 引用argo-rollouts的kustomize语法解释文件

~~~bash
wget https://argoproj.github.io/argo-rollouts/features/kustomize/rollout-transform.yaml

mv rollout-transform.yaml base/

cat > base/kustomization.yaml << 'EOF'
resources:
  - deployment.yaml
  - service.yaml
configMapGenerator:
- name: nginx-index
  files:
  - index.html
configurations: # 使用此配置引用
  - rollout-transform.yaml
EOF
~~~

~~~bash
kustomize build overlays/test/
~~~

~~~bash
kustomize build overlays/test/ | kubectl apply -f -
~~~

~~~bash
kubectl argo rollouts get rollout test-nginx-suffix -n kustomize-kang -w
~~~

- 滚动更新

~~~bash
cat > overlays/test/index.html << 'EOF'
This is argo!
EOF
~~~

~~~bash
kustomize build overlays/test/ | kubectl apply -f -
~~~

~~~bash
kubectl argo rollouts get rollout test-nginx-suffix -n kustomize-kang -w
~~~

- 灰度发布细节
~~~yaml
# 暂停时间的定义
spec:
  strategy:
    canary:
      steps:
        - pause: { duration: 10 }  # 10 seconds
        - pause: { duration: 10s } # 10 seconds
        - pause: { duration: 10m } # 10 minutes
        - pause: { duration: 10h } # 10 hours
        - pause: {}                # pause indefinitely
~~~

~~~bash

~~~

## argo cd

[官方网站](https://argoproj.github.io/argo-cd/)

### 一、安装
[安装文档](https://github.com/argoproj/argo-cd/tree/master/manifests)

- 1.单节点本地集群部署

~~~bash
kubectl create ns argocd
kubectl apply -f https://github.com/argoproj/argo-cd/blob/master/manifests/install.yaml -n argocd
~~~

- 2.单节点部署

~~~bash
kubectl create ns argocd
kubectl apply -f https://github.com/argoproj/argo-cd/blob/master/manifests/namespace-install.yaml
~~~

- 3.高可用本地集群部署

~~~bash
kubectl create ns argocd
kubectl apply -f https://github.com/argoproj/argo-cd/blob/master/manifests/ha/install.yaml
~~~

- 4.高可用部署

~~~bash
kubectl create ns argocd
kubectl apply -f https://github.com/argoproj/argo-cd/blob/master/manifests/ha/namespace-install.yaml
~~~

- 5.目标集群部署定制资源

~~~bash
kubectl apply -k https://github.com/argoproj/argo-cd/manifests/crds\?ref\=stable
~~~

- 6.cli安装

~~~bash
wget http://download.init0.cn/argocd
chmod +x argocd
mv argocd /usr/local/bin
~~~

### 二、命令使用

` <> 中的部分需要依据实际情况替换`

` 详情见 argocd --help`

- 获取admin登陆密码

~~~bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
~~~

- 登陆argo cd

~~~bash
argocd login <argocd_server_ip:port> --username <user_name> --password <password>   --insecure
~~~

- 退出登陆

~~~bash
argocd logout <argocd_server_ip:port>
~~~

- 修改密码

~~~bash
argocd account update-password --account <user_name> --new-password <want_password> --current-password <admin_password>
~~~

- 添加repo

~~~bash
argocd repo add <git_addr> --username <git_user> --password <password> --insecure-skip-server-verification
~~~

- 添加[project](https://argoproj.github.io/argo-cd/user-guide/projects/)

~~~bash
argocd proj create <proj_name>
~~~

- 创建app

~~~bash
argocd app create <app_name> --repo <git_addr> --path <git_dir_path> --revision <git_branch>  --dest-server <deploy_to_k8s_cluster> --dest-namespace <k8s_namespace> --project <project>  --kustomize-version v4.1.3
# 自动同步 可选 --sync-policy auto --auto-prune
~~~

- 同步app

~~~bash
argocd app sync <app_name> --prune
~~~

- 图形化界面

~~~bash
cat > argocd-service.yaml << 'EOF'
apiVersion: v1
kind: Service
metadata:
  namespace: argocd
  labels:
    app.kubernetes.io/component: server
    app.kubernetes.io/name: argocd-server
    app.kubernetes.io/part-of: argocd
  name: argocd-server
spec:
  type: NodePort
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8080
    nodePort: 30002
  - name: https
    port: 443
    protocol: TCP
    targetPort: 8080
  selector:
    app.kubernetes.io/name: argocd-server
EOF

kubectl apply -f argocd-service.yaml

# 访问 http://ip:30002
# 操作部分都可以在图形化界面上完成 但是关于argo cd的配置就只能使用cli或者yaml来修改
~~~

### 三、使用argo cd

- 当我们创建了应用以后 argocd会去对应的git repo中拉取代码

- 根据我们定义的目录 使用kustomize生成yaml 并apply至指定的集群中指定的名称空间

- 如果定义了自动同步 argocd会以三分钟一次的速度比对git的变化

#### 配置git仓库
~~~bash
argocd repo add <git_addr> --username <git_user> --password <password> --insecure-skip-server-verification
~~~

#### 创建应用
~~~bash
argocd app create <app_name> --repo <git_addr> --path <git_dir_path> --revision <git_branch>  --dest-server <deploy_to_k8s_cluster> --dest-namespace <k8s_namespace> --project <project>  --kustomize-version v4.1.3
# 自动同步 可选 --sync-policy auto --auto-prune
~~~

#### 同步应用 
~~~bash
argocd app sync <app_name> --prune
~~~

#### 创建项目
~~~bash
argocd proj create <proj_name>
~~~

#### 自定义资源
~~~bash
    1） project
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: qianfan-test
  namespace: argocd
spec:
  clusterResourceWhitelist:
  - group: '*'
    kind: '*'
  destinations:
  - namespace: '*'
    server: '*'
  namespaceResourceWhitelist:
  - group: '*'
    kind: '*'
  orphanedResources:
    ignore:
    - {}
    warn: false
  sourceRepos:
  - '*'

    2） applications
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: s02c-dpos-rs
  namespace: argocd
spec:
  destination:
    namespace: default
    server: https://kubernetes.default.svc
  project: qianfan-test
  source:
    kustomize:
      version: v4.1.3
    path: qianfan/dpos-rs/overlays/test/s02c
    repoURL: https://gitlab.hd123.com/qianfanops/kubernetes.git
    targetRevision: develop
  syncPolicy:
    automated:
      prune: true
# 无论是project还是app都是k8s中的定制资源 并不仅仅只是运行在argocd程序内部 我们使用cli或者图形界面生成的应用都会被转化成定制资源的形式生成在k8s中 所以所有的应用部署的数据都是存放在etcd中的
~~~

### 用户及权限管理
####  用户及kustomize版本参数配置
~~~bash
vim argocd-cm.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
  namespace: argocd
  labels:
    app.kubernetes.io/name: argocd-cm
    app.kubernetes.io/part-of: argocd
data:
  kustomize.path.v4.1.3: /usr/local/bin/kustomize # 配置kustomize命令的引用
  kustomize.buildOptions.v4.1.3: --load-restrictor LoadRestrictionsNone # 配置kustomize默认参数

  admin.enabled: "false" # 禁用admin
  accounts.qianfan: apiKey, login # 千帆用户账户权限
  accounts.qianfan.enabled: "true" # 启动千帆用户
  # 使用内置的dex服务进行LDAP配置
  url: https://47.118.34.109:9876
  dex.config: |
    connectors:
      - type: ldap
        id: ldap
        name: LDAP
        config:
          host: "ldap.hddomain.cn:389"
          insecureNoSSL: true
          insecureSkipVerify: true
          bindDN: "*@domain.cn"
          bindPW: "password"
          userSearch:
            baseDN: "OU=HDUsers,dc=hddomain,dc=cn"
            filter: ""
            username: "sAMAccountName"
            idAttr: distinguishedName
            emailAttr: mail
            nameAttr: displayName
          groupSearch:
            baseDN: "OU=HDGroups,dc=hddomain,dc=cn"
            filter: ""
            userAttr: distinguishedName
            groupAttr: member
            nameAttr: name

kubectl apply -f argocd-cm.yaml
~~~

#### 用户权限管理
~~~bash
vim argocd-rbac.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-rbac-cm
  namespace: argocd
data:
  #policy.default: role:readonly # 默认策略 只读
  scopes: '[groups,email]' 在ldap组信息中显示邮箱
  policy.csv: |
    # 定义角色权限
    #p, role:org-admin, applications, *,/, allow
    p, role:octopus, applications, *, octopus/*, allow
    #p, role:org-admin, clusters, get, *, allow
    #p, role:org-admin, projects, get, *, allow
    #p, role:org-admin, repositories, get, *, allow
    # 绑定角色
    g, kangpeiwen@hd123.com, role:admin # 邮箱绑定角色
    g, dept_000700110002, role:octopus # 通过组绑定角色
    g, qianfan, role:admin # 用户绑定角色

kubectl apply -f argocd-rbac.yaml
~~~
### argocd web页面访问及代理

#### 设计
- 在heading环境中，我们统一使用 argocd.hd123.com 这个域名来进行访问

~~~bash
1.使用不用的 URL 来区分集群，背后使用 nginx localtion 匹配来进行跳转，所以实际的访问路径为 argocd.hd123.com --> nginx --> argocd addr
2.这里 URL 的设定为 产品-环境，例如 dnet-int dnet-prd pay-prd，所以当我们访问千帆测试的argocd时，访问地址为 argocd.hd123.com/dnet-int/
3.实现上述功能需要添加 nginx pod ，以及配置 argocd 的 argocd-server
~~~

#### 配置

~~~bash
# argocd 配置
kubectl edit deploy argocd-server -n argocd
...
        - --rootpath
        - /dnet-int 
        - --insecure
...

- 前两行为定义 argocd 初始 URL 配置
- 后一行为采用 http 协议，不使用 https 配置

# nginx 配置
- 此代理 nginx 部署在 dnet-prd 集群中，由 argocd 管理
- 其配置存放在 dnet 的 git 中，地址为 http://gitlab.app.hd123.cn:10080/qianfanops/toolset 分支为 k8s_dnet，路径为 associated-resources/argo/prd/argocd-proxy.yaml

# k8s集群配置
在部署 argocd 时，默认使用 nodeport 方式在主机上暴露 30002 端口，所以需要在 ingress 对应的 slb 上将 30002 端口暴露到公网，而后修改 dnet-prd 环境下的 nginx 配置进行代理
~~~