## 一、前期准备

### 1）服务基础信息

```bash
1.应用名称（镜像名）
- 例如镜像为 harbor.qianfan123.com/baas/fms-service:1.25.0-SNAPSHOT
- 则应用名称为 fms-service

2.环境
- 此环境为 toolset develop 分支中 对应的 settings-xxx.yaml 中的环境信息
- 例如 integration_test 环境、daojia_int 环境

3.stackid
- stackid 存放在产品对应的 cmdb 中

4.镜像全称（只有创建base目录时填写）
- 镜像名称包括仓库名称，如 harbor.qianfan123.com/baas/fms-service:1.25.0-SNAPSHOT
```

### 2）apollo相关信息

```bash
# apollo相关信息为 apollo创建项目后工单备注的信息
1.apollo服务名称
- 例如 9999.fms-service.baas

2.apollo环境
- apollo的环境只有 INT BRA UAT PRD 四种
- 根据对应的服务基础信息填写

3.apollo token
- 见工单备注
```

### 3）其他信息

```bash
1.argocd URL后缀（k8s集群地址）
- 通常为 产品代码-环境  如 /dnet-int /dnet-prd
- 具体信息以部署集群为准 上述只做参考

2.k8s 名称空间
- 此信息为部署到集群后服务所在名称空间
- 通常为 产品代码-环境全称 如 dnet-integration-test，具体以部署目的为指向，可将多个产品部署至同一名称空间
```

### 4）配置分离

```bash
# 配置分为两部分：运维配置、开发配置
1.运维配置
- 主要以账号密码类配置为主，以及一些开发不关心的配置
- 存放在 git 中

2.开发配置
- 开发关心或修改的配置
- 存放在 apollo 中
```

## 二、生成

### 1） git

```bash
1.创建新分支
- 从 toolset 的 k8s_init 分支创建新分支

2.旧分支
- 将 toolset 的 k8s_init 分支中的 init 目录复制至当前分支根目录
```

### 2）填入信息

```bash
1.打开init.txt

2.按照前期准备填入信息
"应用名称 环境 stackid image" "apollo服务名称 apollo环境 apollo token" "argocdURL后缀 k8s名称空间"
- 例如: # 注意引号
cat >> init.txt <<'EOF'
"fms-service daojia_int daojia_int harbor.qianfan123.com/baas/fms-service:1.25.0-SNAPSHOT" "9999.fms-service.baas INT bc5f1869b70391bba94866e04c18b62483be25d3" "dnet-int baas-integration-test"
EOF

3.将运维配置放入./config/应用名_环境名.env
- 例如./config/fms-service_daojia_int.env
```

### 3）生成

```bash
./start.sh
```

### 4）校验

```bash
1.查看结构
tree ../应用名称
- 例如 tree ../fms-service

2.查看 deployment 的名称和镜像配置
cat ../应用名称/base/deployment.yaml
- 例如 cat ../fms-service/base/deployment.yaml

3.查看 overlay 中 config.yaml 的配置
cat ../应用名称/overlays/环境/stackid/config.yaml
- 例如 cat ../fms-service/overlays/daojia_int/daojia_int/config.yaml

4.检查配置
cat ../应用名称/overlays/环境/stackid/ops.j2
- 例如 cat ../fms-service/overlays/daojia_int/daojia_int/ops.j2
```

## 三、细节
~~~bash
init
├── base
│   ├── deployment.yaml
│   ├── kustomization.yaml
│   ├── rollout-transform.yaml
│   └── service.yaml
├── config
│   └── test_int.env
├── init.
├── init.sh
├── init.txt
├── overlays
│   ├── cm-suffix-transformer.yaml
│   ├── config.yaml
│   ├── deployment.yaml
│   └── kustomization.yaml
└── start.sh
~~~

### 1）基础文件

~~~bash
# cat init/base/deployment.yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: app_name
spec:
  replicas: 1
  selector:
    matchLabels:
      name: app_name
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
  template:
    metadata:
      labels:
        name: app_name
    spec:
      imagePullSecrets:
        - name: "harbor.qianfan123.com"
        - name: "harborka.qianfan123.com"
      hostAliases:
      - ip: "47.97.75.60"
        hostnames:
        - "kafka-aliyun-test-cluster.com"
      containers:
        - name: app_name
          image: "image_name"
          imagePullPolicy: Always
          lifecycle:
            preStop:
              exec:
                command: ["/bin/sh","-c","curl -X POST  127.0.0.1:8443/actuator/shutdown"]
          volumeMounts:
            - mountPath: /opt/heading/tomcat/logs
              name: logs-storage
            - name: timezone
              mountPath: /etc/localtime
              readOnly: true
          ports:
            - containerPort: 8080
              protocol: TCP
            - containerPort: 8443
              protocol: TCP
          envFrom:
          - configMapRef:
              name: app_name-env
          - secretRef:
              name: app_name-secret
          livenessProbe:
            httpGet:
              path: /actuator/health
              port: 8443
              scheme: HTTP
            failureThreshold: 3
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 1
          startupProbe:
            httpGet:
              path: /actuator/health
              port: 8443
              scheme: HTTP
            failureThreshold: 30
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 1
          readinessProbe:
            httpGet:
              path: /actuator/health
              port: 8443
              scheme: HTTP
            failureThreshold: 2
            periodSeconds: 6
            successThreshold: 3
            timeoutSeconds: 1
          resources:
            limits:
              cpu: 800m
              memory: 1100Mi
            requests:
              cpu: 100m
              memory: 900Mi
        - name: filebeat
          image: harbor.qianfan123.com/elk/filebeat:7.13.0
          args: [
            "-c", "/etc/filebeat/filebeat-kafka.yml",
            "-e",
          ]
          ports:
            - containerPort: 5678
              protocol: TCP
          securityContext:
            runAsUser: 0
            privileged: true
          livenessProbe:
            tcpSocket:
              port: 5678
            failureThreshold: 3
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 1
          resources:
            limits:
              cpu: 200m
              memory: 180Mi
            requests:
              cpu: 50m
              memory: 180Mi
          volumeMounts:
          - name: config
            mountPath: /etc/filebeat/
          - name: logs-storage
            mountPath: /opt/heading/tomcat/logs
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
      - name: timezone
        hostPath:
          path: /usr/share/zoneinfo/Asia/Shanghai
# 此文件为 k8s deployment 控制器的配置，只需要将 app_name 和 image_name 替换为对应值即可
~~~


~~~bash
# cat init/base/service.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/path: "/actuator/prometheus"
    prometheus.io/port: "8443"
  name: app_name
spec:
  selector:
    name: app_name
  ports:
    - name: http
      port: 8080
      targetPort: 8080
# 此文件为 k8s service 的配置，只需要将 app_name 替换为对应值即可
~~~

~~~bash
# cat init/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
metadata:

configurations:
- rollout-transform.yaml

resources:
- deployment.yaml
- service.yaml
- ../../associated-resources/harbor.yaml
- ../../associated-resources/harborka.yam

# 此文件为 kustomize 的引用文件，不需要做改变
~~~

~~~bash
# cat init/overlays/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  app: app_name
  version: v1


patchesStrategicMerge:
- deployment.yaml

configMapGenerator:
- name: app_name-env
  behavior: create
  envs:
  - dev.env

secretGenerator:
- name: app_name-secret
  envs:
  - ops.env
  type: Opaque
- name: filebeat-secret
  files:
  - ../../../../associated-resources/int/filebeat-kafka.yml
  type: Opaque
resources:
- ../../../base

transformers:
- cm-suffix-transformer.yaml
# 此文件为 overlays 中对应环境的引用文件，只需要将 app_name 替换为对应值即可
~~~

~~~bash
# cat init/overlays/deployment.yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: app_name
spec:
  replicas: 1
  strategy:
    canary:
      steps:
      - setWeight: 50
      - pause: {duration: 3}
# 此文件为 overlays 中对应环境的 k8s deployment 控制器的配置，只需要将 app_name 替换为对应值即可
~~~

~~~bash
# cat init/overlays/config.yaml
apiVersion: v0.1

# apollo相关配置， addr有默认值可以不填写
apollo:
  addr: "http://apollo-portal.hd123.com"
  app_id: "apollo_app_id"
  env: "apollo_env"
  token: "apollo_token"
  cluster: "default"
  namespace: "application"

# argocd相关配置， addr有默认值可以不填写， server如果是本地所在集群则使用默认值可以不填
argocd:
  addr: "http://argocd.hd123.com"
  cluster: "k8s_cluster"
  namespace: "k8s_namespace"
  server: "https://kubernetes.default.svc"
  kustomize_version: "v4.1.3"
# 此文件为 toolset 读取的配置文件，分为 apollo 与 argocd 配置两部分
~~~

### 2）脚本

~~~bash
cat init/init.sh
#!/bin/bash

base=($1)
apollo=($2)
other=($3)

function create_base() {
    mkdir -p ../${app}/{base,overlays}
    cp ./base/* ../${app}/base/
    sed -i "s#app_name#${app}#g" ../${app}/base/deployment.yaml
    sed -i "s#app_name#${app}#g" ../${app}/base/service.yaml
    sed -i "s#image_name#${image_name}#g" ../${app}/base/deployment.yaml
    echo "  ${app}/base 创建成功"
}

function create_overlays() {
    mkdir -p ../${app}/overlays/${env_name}/${stackid}
    cp ./overlays/* ../${app}/overlays/${env_name}/${stackid}
    sed -i "s#apollo_app_id#${apollo_app_id}#g" ../${app}/overlays/${env_name}/${stackid}/config.yaml
    sed -i "s#apollo_env#${apollo_env}#g" ../${app}/overlays/${env_name}/${stackid}/config.yaml
    sed -i "s#apollo_token#${apollo_token}#g" ../${app}/overlays/${env_name}/${stackid}/config.yaml
    sed -i "s#k8s_cluster#${k8s_cluster}#g" ../${app}/overlays/${env_name}/${stackid}/config.yaml
    sed -i "s#k8s_namespace#${k8s_namespace}#g" ../${app}/overlays/${env_name}/${stackid}/config.yaml
    sed -i "s#app_name#${app}#g" ../${app}/overlays/${env_name}/${stackid}/deployment.yaml
    sed -i "s#app_name#${app}#g" ../${app}/overlays/${env_name}/${stackid}/kustomization.yaml
    cp ./config/${app}_${env_name}.env ../${app}/overlays/${env_name}/${stackid}/ops.j2
    echo "  ${app}/overlays/${env_name}/${stackid} 创建成功"
}

function delete_apollo() {
    sed -i 2,10d ../${app}/overlays/${env_name}/${stackid}/config.yaml
}

if [ $(echo ${#base[*]}) -ge 3 ];then
    app=${base[0]}
    env_name=${base[1]}
    stackid=${base[2]}
    image_name=${base[3]}
    base_path=$(find ../${app}/base/ -type d 2> /dev/null | xargs echo )
    [ _${base_path} == _"" ] && create_base
fi

if [ $(echo ${#apollo[*]}) -eq 3 -a $(echo ${#other[*]}) -eq 2 ];then
    apollo_app_id=${apollo[0]}
    apollo_env=${apollo[1]}
    apollo_token=${apollo[2]}
    k8s_cluster=${other[0]}
    k8s_namespace=${other[1]}
    overlay_path=$(find ../${app}/overlays/${env_name}/${stackid} -type d 2> /dev/null | xargs echo )
    [ _${overlay_path} == _"" ] && create_overlays
elif [ $(echo ${#apollo[*]}) -ne 3 -a $(echo ${#other[*]}) -eq 2 ];then
    k8s_cluster=${other[0]}
    k8s_namespace=${other[1]}
    overlay_path=$(find ../${app}/overlays/${env_name}/${stackid} -type d 2> /dev/null | xargs echo )
    [ _${overlay_path} == _"" ] && create_overlays && delete_apollo
fi

# 此脚本问生成脚本，由 start.sh 调用
~~~

~~~bash
# cat init/start.sh
#!/bin/bash

sed "s#^#./init.sh #"  init.txt | bash
~~~