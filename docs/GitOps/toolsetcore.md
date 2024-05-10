

# Toolsetcore

解析toolsetcore如何实现k8s相关部署等功能的逻辑原理

## 目录

* [Toolsetcore](#Toolsetcore)
   * [配置文件](#配置文件)
   * [环境变量](#环境变量)
   * [功能](#功能)
      * [使用说明](#使用说明)
      * [下载配置](#下载配置)
      * [配置检查](#配置检查)
      * [配置对比](#配置对比)
      * [修改版本](#修改版本)
      * [数据库更新](#数据库更新)
      * [创建Argocd项目](#创建Argocd项目)
      * [部署](#部署)
   * [F&Q](#F&Q)

## 配置文件

用于给toolsetcore读取和使用的配置文件，后续可能会做变化和迭代，

文件名为`config.yaml`存储在`overlay`中和`kustomization.yaml`同级

因为配置并没有被kustomize读取和使用，所以和argo并没有关联

```
# 文件存储坐标
app1
├── base
│   └── kustomization.yaml
└── overlays
    ├── bra
    │   ├── config.yaml
    │   └── kustomization.yaml
    ├── int
    │   ├── config.yaml
    │   └── kustomization.yaml
    └── prd
        ├── config.yaml
        └── kustomization.yaml
```

详细配置解析如下:

```yaml
apiVersion: v0.1    # 作为表示配置的版本
apollo:             # 定义apollo的相关信息，如果没有apollo请不要填写
  addr: http://apollo-portal.hd123.com  # apollo的地址
  app_id: xxxx.xxx        # 该组件对应的配置ID
  cluster: default        # apollo内的集群名，默认都是default
  env: INT                # 配置使用的分支
  namespace: application  # apollo中的名称空间
  token: xxxx             # 连接apollo使用的token信息
argocd:             # 定义argocd的地相关信息，是必填字段
  addr: http://argocd.hd123.com  # argocd的域名地址
  cluster: dnet-int              # argocd的集群名，也是argocd的一级path的值
  kustomize_version: v4.1.3      # 使用的kustomize的版本
  namespace: xxxx                # 在k8s中设置的名称空间
  server: https://kubernetes.default.svc  # argocd连接k8s的apiserver的地址
subsystem:           # 定义组件的相关信息
  version: 1.0.0-SNAPSHOT        # 当前组件的部署版本
```

## 环境变量

| 变量名                           | 含义                                 | 默认值                            |
| -------------------------------- | ------------------------------------ | --------------------------------- |
| DNET_PROJECT                     | 项目名                               | heading                           |
| DNET_PRODUCT                     | 产品名                               | dnet                              |
| DNET_PROFILE                     | 环境名                               | integration_test                  |
| GIT_USER                         | 连接git的账号                        | qianfan                           |
| GIT_PASSWORD                     | 连接git的密码                        | xxxx                              |
| GIT_BASE_URL                     | toolset所在的git的地址               | `https://github-argocd.hd123.com` |
| GIT_BASE_GROUP                   | toolset所在的git的组信息             | qianfanops                        |
| DBUPDRADE_SKIP_ERROR_APPLICATION | 是否跳过应用失败的时候数据库更新报错 | False                             |

- 当前配置文件定位的方式是
  - git地址： `{GIT_BASE_URL}/{GIT_BASE_GROUP}/toolset-{DNET_PROJECT}` 
  - git branch： `k8s_{DNET_PRODUCT}` 
  - 文件位置： `{{subsystem}}/overlays/{DNET_PROFILE}/kustomization.yaml`
    - 其中`subsystem`组件名是脚本的必填项，不是由环境变量注入

## 功能

所有的指令都是用指定的action (`kubernetes`) 进行标识的

### 使用说明

所有功能都是，toolset的 `develop` 分支的`settings.yaml`结合 `k8s_xxx`分支的 `config.yaml` 和`kustomize.yaml`实现完成的

所以执行命令前需要准备配置文件

```shell
DNET_PROFILE=integration_test DNET_PRODUCT=dnet TRUST_PUBLIC_IP=${run_on_public} hdops download_toolset --branch develop -p .

tar zxf toolset.tar.gz -C .

DNET_PROJECT=${DNET_PROJECT} DNET_PRODUCT=${DNET_PRODUCT} DNET_PROFILE=${DNET_PROFILE} GIT_USER=${GITLAB_USER_USR} GIT_PASSWORD=${GITLAB_USER_PSW} hdops kubernetes download_kubernetes -s None
```

### 下载配置

#### 使用方法

`hdops kubernetes download_kubernetes --subsystem None`

- `subsystem` 指代组件名，是必填字段，下载配置的时候填写None即可

#### 程序逻辑

根据环境变量`GIT_BASE_URL`、`GIT_BASE_GROUP`、`DNET_PROJECT`、`DNET_PRODUCT`定位出git地址

默认git repo的命名方式为`toolset-{project}`

- 特例：`heading` 项目的git repo为`toolset`， `miniso`项目的git repo为`toolset_miniso`

git的分支是根据产品区分，命名规范为`k8s_{product}`

[代码位置, kustomize.py(download_kubernetes)](http://github.app.hd123.cn:8080/qianfanops/toolsetcore/blob/develop/hdtoolsetcore/kubernetes/kustomize.py#L88) 

### 配置检查

#### 使用方法

`hdops kubernetes validate --subsystem ${image} --stackids None`

- `subsystem` 指代组件名，是必填字段，检查配置的组件名，有多个用逗号隔开
  - `subsystem` 如果为None，则会检查所有组件的配置是否符合要求
- `stackids` 指代资源栈，非必填字段，会对组件下指定的资源栈进行配置检查

#### 程序逻辑

根据`subsystem`定位到指定环境的`config.yaml`，检查其中字段是否都存在，如果有字段未定义则会报错退出

1. 检查`subsystem`字段是否存在，其中是否存在`["version"]`字段
2. 检查`apollo`字段是否存在
   - 如果存在，则检查其中是否包含`["token", "env", "app_id", "cluster", "namespace"]`字段
   - 如果不存在，则表示该组件没有注册使用apollo记录配置文件
3. 检查`argocd`字段是否存在，其中是否存在`["cluster", "namespace", "kustomize_version"]`字段

[代码位置, validate.py(do_validate)](http://github.app.hd123.cn:8080/qianfanops/toolsetcore/blob/develop/hdtoolsetcore/kubernetes/validate.py#L44) 

### 配置对比

#### 使用方法

`hdops kubernetes compare_config --subsystem ${image} --stackids None`

- `subsystem` 指代组件名，是必填字段，配置比对的组件名，有多个用逗号隔开
  - `subsystem` 如果为None，则会对比所有组件的配置
- `stackids` 指代资源栈，非必填字段，会对组件下指定的资源栈进行配置对比

#### 程序逻辑

根据`subsystem`定位到指定环境的`dev.env、ops.env、ops.j2`，根据`config.yaml`定位到指定的apollo地址

1. 读取apollo的配置信息和`dev.env`比较
2. 读取`ops.j2`的配置和`ops.env`比较
3. 如果配置新增某个词条，输出`+ xxx=zzz`
4. 如果配置删除或修改某个词条，输出`- xxx=zzz`，并且提示警告`存在隐患`

[代码位置, config.py(compare_config)](http://github.app.hd123.cn:8080/qianfanops/toolsetcore/blob/develop/hdtoolsetcore/kubernetes/config.py#L34) 

### 修改版本

#### 使用方法

`hdops kubernetes modify_version -s None --imageversion ${image:version}`

`hdops kubernetes modify_version --subsystem ${image} --version ${version}`

- `imageversion` 指代组合的组件名+版本号，是修改组合版本模式的必填字段

  - 写法模型：`{version1}:{image1.1},{image1.2};{version2}:{image2.1};{version3}:{image3.1}`

  - 例如: `1.31.0-SNAPSHOT:octopus-server,octopus-monitor;1.23.0-SNAPSHOT:octopus-etl`

    修改`octopus-server和octopus-monitor`版本到`1.31.0-SNAPSHOT`，修改`octopus-etl`版本到`1.23.0-SNAPSHOT`

- `subsystem` 指代组件名，是必填字段，配置比对的组件名

- `version` 指代版本，是修改单个版本的必填字段，将指定组件修改为指定版本

- 如果填写了`imageversion`，则只会执行组合版本，如果没填写`imageversion`，则自动执行修改单个版本

#### 程序逻辑

1. (只对修改组合版本的逻辑)，根据写法模型拆解输入，循环执行多个组件版本
2. 遍历制定组件后所有的stackid
3. 读取`config.yaml`中的`subsystem.version`记录为当前版本状态
4. 对比当前版本和输入中的指定版本，如果版本不同则修改`config.yaml`中的`version`字段
5. 提交修改到git仓库，并统一输出版本修改情况
   - 发生修改的输出`{image} 镜像制品版本从 {verison_old} 变更为 {verison_new}`
   - 版本没有发生修改的输出`{image} 版本不需要更新: 当前版本为 {verison}`

[修改组合版本, 代码位置, modify_version.py(do_modify_versions)](http://github.app.hd123.cn:8080/qianfanops/toolsetcore/blob/develop/hdtoolsetcore/kubernetes/modify_version.py#L53) 

[修改单个版本, 代码位置, modify_version.py(do_modify_version)](http://github.app.hd123.cn:8080/qianfanops/toolsetcore/blob/develop/hdtoolsetcore/kubernetes/modify_version.py#L46) 

### 数据库更新

#### 使用方法

`hdops kubernetes db_upgrade --subsystem ${image} --stackids None --version None --dblist None --no-skipversion --no-skipversionupdate --receivers ${receivers} --timeout ${timeout} --threadcount ${threadcount}`

- `subsystem` 指代组件名，是必填字段，执行数据库升级的组件
- `stackids` 指代资源栈，非必填字段，会对组件下指定的资源栈进行数据库升级
- `version` 指代指定具体的数据库版本进行更新，非必填字段，默认是none即组件的版本
- `dblist` 指代指定具体的数据库库名进行更新，非必填字段，默认是none即全部
- `--skipversion/--no-skipversion` 是否需要跳过版本检查，传递给数据库更新容器的配置，非必填字段，默认不跳过
- `--skipversionupdate/--no-skipversionupdate` 是否需要跳过版本更新，传递给数据库更新容器的配置，非必填字段，默认不跳过
- `receivers` 输入如果数据库更新失败，邮件通知的人，非必填字段，默认为`buhaiqing@hd123.com`
- `timeout` 数据库更新的超时设置，非必填字段，默认45秒
- `threadcount` 指定本次更新的并发数，默认3次

#### 程序逻辑

1. 获取更新数据库的实例信息，逻辑和之前的`toolset+cmdb`一致，生成需要更新的`job`数组
   - 根据不同的产品不同处理方式，`dly、vss`等产品特殊处理
   - 根据`stackid`获取`rdb`信息，根据`toolset`找到符合要求的`rdb_db`
2. 根据输入`threadcount`拆分需要更新的`job`数组
3. 根据本次`job`中的组件状态，如果状态异常，则会先将所有历史`k8s_job`删除
4. 将本次执行跟新的job内容写入`job.yaml`中提交到git仓库，触发`argocd`的更新
5. `argocd`根据配置创建多个`k8s_job`执行任务
6. 监听`argocd`然后`k8s_job`执行完成返回成功，则代表数据库更新完成，可以进行下一批次的更新任务
7. 如果更新任务中出现失败，则会停止后续的任务进度，报错退出

[代码位置, dbUpgrade.py(do_db_upgrade)](http://github.app.hd123.cn:8080/qianfanops/toolsetcore/blob/develop/hdtoolsetcore/kubernetes/dbUpgrade.py#L363) 

### 创建Argocd项目

#### 使用方法

`hdops kubernetes create_argocd_application --subsystem ${image} --stackids None`

- `subsystem` 指代组件名，是必填字段，创建Argocd内项目的组件，有多个用逗号隔开
- `stackids` 指代资源栈，非必填字段，会对组件下指定的资源栈进行创建Argocd内项目

#### 程序逻辑

1. 读取`config.yaml`中的`argocd`相关配置
2. 根据环境变量`{DNET_PRODUCT}`定义argocd中的`project`
3. 根据环境变量`{GIT_xxx}`创建git需要的认证
4. 组合输出中所有的`{stackid}-{subsystem}`，创建argocd内的应用

[代码位置, argocd.py(create_argocd_application)](http://github.app.hd123.cn:8080/qianfanops/toolsetcore/blob/develop/hdtoolsetcore/kubernetes/argocd.py#L366) 

### 部署

#### 使用方法

`hdops kubernetes appinstall --subsystem ${image} --stackids None --checkstatus ${checkstatus}`

- `subsystem` 指代组件名，是必填字段，更新部署的组件
- `stackids` 指代资源栈，非必填字段，会对组件下指定的资源栈进行更新部署
- `checkstatus` 指代是否监听状态，非必填字段，在部署完成后是否等待监听argocd的项目状态

#### 程序逻辑

1. 根据输入的`subsystem`和`stackid`循环处理所有需要部署的组件以及资源栈
2. 读取`config.yaml`中的`version`信息，将其写入到`kustomization.yaml`中
3. 根据渲染逻辑，将`ops.j2`的配置渲染后写入`ops.env`
4. 读取`apollo`内的配置文件，将其写入`dev.env`中
5. 根据开发规范`sc_220510`，记录组件信息到`code_report`
6. 提交Git配置到仓库中
7. 再次遍历输入的`subsystem`和`stackid`，触发同步，输出所有argo的地址
8. 判断是否有环境变量`GATEWAY_NAME`，进行CD部署成功的记录
9. 根据输入`checkstatus`是否为`true`，决定是否进行argocd的监听
   - 如果监听，则会对argocd每10s检查一次状态，检查18次，也就是3分钟
   - 如果argocd状态为`Healthy`则代表部署成功
   - 如果argocd状态为`Degraded`则代表部署失败
   - 如果argocd状态为`Missing`则代表同步触发异常，会重新触发同步

[代码位置, appInstall.py(do_appinstall)](http://github.app.hd123.cn:8080/qianfanops/toolsetcore/blob/develop/hdtoolsetcore/kubernetes/appInstall.py#L32) 

## F&Q

### 修改版本对线上服务的影响

- 修改版本只会将版本改到`config.yaml`中，`argocd`并不会感知到变化
- 更新数据库的时候读取的数据库更新的版本读取的是`config.yaml`内写的配置
- 发版的时候，会将`config.yaml`内的版本信息写入到`kustomization.yaml`中，`argocd`这时候才会感知到版本变化
- 综上: 修改版本其实只是暂存了一个后续部署中使用的版本信息，与之前的`patch`文件修改一样的效果

### 数据库更新异常

- 检查`argocd`的应用状态是否正常
- 如果组件pod异常，根据pod的日志，修复应用至正常，修复过程中使用`deploy_non_db`的job进行部署
- 如果`job`异常，根据`job`的日志，修复数据库更新镜像即可
