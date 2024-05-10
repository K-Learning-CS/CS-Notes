 ### 一、安装

#### Argo Workflows

[官方文档](https://argoproj.github.io/argo-workflows/installation/)

##### 1.安装

```yaml
1.找到对应版本并安装
https://github.com/argoproj/argo-workflows/releases

2.暴露 service
apiVersion: v1
kind: Service
metadata:
  name: argo-server
  namespace: argo
spec:
  ports:
  - name: web
    nodePort: 30003
    port: 2746
    protocol: TCP
    targetPort: 2746
  selector:
    app: argo-server
  type: NodePort
```

##### 2.使用argocd dex登陆

[官方文档](https://argoproj.github.io/argo-workflows/argo-server-sso-argocd/)

~~~bash
1.创建通信密钥并在 argo-workflow 与 argo-cd 名称空间中部署
# 密钥内容为引号中字符使用base64加密后的结果 可以自己定义
# 如果自定义请将以下所有资源中相关处修改为变更后的内容
apiVersion: v1
kind: Secret
metadata:
  name: argo-workflows-sso
data:
  # client-id is 'argo-workflows-sso'
  client-id: YXJnby13b3JrZmxvd3Mtc3Nv
  # client-secret is 'MY-SECRET-STRING-CAN-BE-UUID'
  client-secret: TVktU0VDUkVULVNUUklORy1DQU4tQkUtVVVJRA==


2.向 argo-cd 的组件 argocd-dex-server 中添加 env
kubectl  -n argocd edit deployments.apps argocd-dex-server

apiVersion: apps/v1
kind: Deployment
metadata:
  name: argocd-dex-server
spec:
  template:
    spec:
      containers:
        - name: dex
        # 以下为添加的内容
          env:
            - name: ARGO_WORKFLOWS_SSO_CLIENT_SECRET
              valueFrom:
                secretKeyRef:
                  name: argo-workflows-sso
                  key: client-secret

3.向 dex 认证配置中加入 argo-workflow
kubectl  -n argocd edit configmaps argocd-cm

apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
data:
  dex.config: |
  # 以下为添加的内容
    staticClients:
      - id: argo-workflows-sso
        name: Argo Workflow
        redirectURIs:
          - https://argo-workflows.mydomain.com/oauth2/callback # argo-workflows地址
        secretEnv: ARGO_WORKFLOWS_SSO_CLIENT_SECRET

3.向 argo-server 添加参数 --auth-mode=sso
kubectl  -n argo edit deployments.apps argo-server

apiVersion: apps/v1
kind: Deployment
metadata:
  name: argo-server
spec:
  template:
    spec:
      containers:
        - name: argo-server
          args:
            - server
            - --auth-mode=sso

4.修改配置
kubectl  -n argo edit configmaps workflow-controller-configmap

apiVersion: v1
kind: ConfigMap
metadata:
  name: workflow-controller-configmap
data:
  sso: |
    issuer: https://argo-cd.mydomain.com/api/dex # argocd dex地址
    clientId: # 最开始定义的密钥
      name: argo-workflows-sso
      key: client-id
    clientSecret:
      name: argo-workflows-sso
      key: client-secret
    redirectUrl: https://argo-workflows.mydomain.com/oauth2/callback # argo-workflows地址
~~~

#### Pipelines

##### 安装

```bash
1.安装pipeline
https://github.com/argoproj-labs/argo-dataflow/blob/main/docs/QUICK_START.md

kubectl apply -f https://raw.githubusercontent.com/argoproj-labs/argo-dataflow/main/config/quick-start.yaml

2.添加 argo 权限
kubectl edit clusterrole argo-server-cluster-role

...
- apiGroups:
  - dataflow.argoproj.io
  resources:
  - pipelines
  - steps
  verbs:
  - create
  - get
  - list
  - watch
  - update
  - patch
  - delete
```

#### argo-events

##### 安装

```
https://argoproj.github.io/argo-events/installation/

kubectl create namespace argo-events
kubectl apply -f https://raw.githubusercontent.com/argoproj/argo-events/stable/manifests/install.yaml
kubectl apply -f https://raw.githubusercontent.com/argoproj/argo-events/stable/manifests/install-validating-webhook.yaml
kubectl apply -n argo-events -f https://raw.githubusercontent.com/argoproj/argo-events/stable/examples/eventbus/native.yaml
```

### 二、使用

#### rdb-upgrade


##### WorkflowTemplate
- 模版用于 workflow 引用，在引用时只需要传递需要改变的参数即可
```yaml
apiVersion: argoproj.io/v1alpha1
kind: WorkflowTemplate
metadata:
  name: rdb-upgrade-template
spec:
  entrypoint: main
  arguments:
    parameters:
      - name: image
        value: harbor.qianfan123.com/baas/mpas-rdb-upgrade:1.2.4-SNAPSHOT
      - name: datasource
        value: rm-bp1496f14i5b9xaod.mysql.rds.aliyuncs.com:3306/mpas
      - name: user
        value: baas
      - name: passwd
        value: cdiWappZtnd6GMMu
      - name: arg
        value: ""
      - name: skip
        value: false
  templates:
    - name: main
      inputs:
        parameters:
          - name: image
          - name: datasource
          - name: user
          - name: passwd
          - name: arg
          - name: skip
      steps:
        - - name: rdb-upgrade
            template: rdb-upgrade-template
            arguments:
              parameters:
                - name: image
                  value: "{{inputs.parameters.image}}"
                - name: datasource
                  value: "{{inputs.parameters.datasource}}"
                - name: user
                  value: "{{inputs.parameters.user}}"
                - name: passwd
                  value: "{{inputs.parameters.passwd}}"
                - name: skip
                  value: "{{inputs.parameters.skip}}"
            when: "{{inputs.parameters.skip}} == false"
          - name: rdb-upgrade-skip
            template: rdb-upgrade-skip-template
            arguments:
              parameters:
                - name: image
                  value: "{{inputs.parameters.image}}"
                - name: datasource
                  value: "{{inputs.parameters.datasource}}"
                - name: user
                  value: "{{inputs.parameters.user}}"
                - name: passwd
                  value: "{{inputs.parameters.passwd}}"
                - name: arg
                  value: "{{inputs.parameters.arg}}"
                - name: skip
                  value: "{{inputs.parameters.skip}}"
            when: "{{inputs.parameters.skip}} == true"
    - name: rdb-upgrade-template
      inputs:
        parameters:
          - name: image
          - name: datasource
          - name: user
          - name: passwd
      container:
        image: "{{inputs.parameters.image}}"
        command: ["java","-jar","upgrade.jar"]
        args: ["all","-d","jdbc:mysql://{{inputs.parameters.datasource}}","-u","{{inputs.parameters.user}}","-p","{{inputs.parameters.passwd}}"]
    - name: rdb-upgrade-skip-template
      inputs:
        parameters:
          - name: image
          - name: datasource
          - name: user
          - name: passwd
          - name: arg
      container:
        image: "{{inputs.parameters.image}}"
        command: ["java","-jar","upgrade.jar"]
        args: ["all","-d","jdbc:mysql://{{inputs.parameters.datasource}}","-u","{{inputs.parameters.user}}","-p","{{inputs.parameters.passwd}}","{{inputs.parameters.arg}}"]
```
##### Workflow
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: mpas-rdb-upgrade-
spec:
  entrypoint: main
  templates:
    - name: main
      steps:
        - - name: rdb-upgrade
            templateRef:
              name: rdb-upgrade-template
              template: main
            arguments:
              parameters:
                - name: image
                  value: harbor.qianfan123.com/baas/mpas-rdb-upgrade:1.2.4-SNAPSHOT
                - name: datasource
                  value: rm-bp1496f14i5b9xaod.mysql.rds.aliyuncs.com:3306/mpas
                - name: user
                  value: baas
                - name: passwd
                  value: cdiWappZtnd6GMMu
                - name: arg
                  value: --skip-version-check
                - name: skip
                  value: false
```

#### argocd sync
##### WorkflowTemplate

```yaml
apiVersion: argoproj.io/v1alpha1
kind: WorkflowTemplate
metadata:
  name: argocd-sync-template
spec:
  entrypoint: main
  templates:
  - name: main
    inputs:
      parameters:
      - name: argocd-version
        value: v2.3.4
      - name: application-name
        value: daojia-int-qw-mpas-service
      - name: flags
        value: --insecure
      - name: argocd-server-address
        value: argocd.hd123.com
      - name: argocd-rootpath
        value: /dnet-int/
      - name: argocd-credentials-secret
        value: argocd-secret
    script:
      image: argoproj/argocd:{{inputs.parameters.argocd-version}}
      command: [bash]
      env:
        - name: ARGOCD_USERNAME
          valueFrom:
            secretKeyRef:
              name: "{{inputs.parameters.argocd-credentials-secret}}"
              key: username
              optional: true
        - name: ARGOCD_PASSWORD
          valueFrom:
            secretKeyRef:
              name: "{{inputs.parameters.argocd-credentials-secret}}"
              key: password
              optional: true
        - name: ARGOCD_SERVER
          value: "{{inputs.parameters.argocd-server-address}}"
      source: |
        #!/bin/bash
        set -euo pipefail
        argocd login "$ARGOCD_SERVER" --grpc-web-root-path "{{inputs.parameters.argocd-rootpath}}" --username=$ARGOCD_USERNAME --password=$ARGOCD_PASSWORD {{inputs.parameters.flags}}
        echo "Running as ArgoCD User:"
        argocd account get-user-info {{inputs.parameters.flags}}
        argocd app sync {{inputs.parameters.application-name}} {{inputs.parameters.flags}}
        argocd app wait {{inputs.parameters.application-name}} --health {{inputs.parameters.flags}}

```
##### Workflow
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: mpas-rdb-upgrade-
spec:
  entrypoint: main
  templates:
    - name: main
      steps:
        - - name: argocd-sync
            templateRef:
              name: argocd-sync-template
              template: main
            arguments:
              parameters:
                - name: argocd-version
                  value: v2.3.4
                - name: application-name
                  value: daojia-int-qw-mpas-service
                - name: flags
                  value: --insecure
                - name: argocd-server-address
                  value: argocd.hd123.com
                - name: argocd-rootpath
                  value: /dnet-int/
                - name: argocd-credentials-secret
                  value: argocd-secret
```
#### rdb-upgrade & argocd sync
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: mpas-rdb-upgrade-
spec:
  entrypoint: main
  templates:
    - name: main
      steps:
        - - name: rdb-upgrade
            templateRef:
              name: rdb-upgrade-template
              template: main
            arguments:
              parameters:
                - name: image
                  value: harbor.qianfan123.com/baas/mpas-rdb-upgrade:1.2.4-SNAPSHOT
                - name: datasource
                  value: rm-bp1496f14i5b9xaod.mysql.rds.aliyuncs.com:3306/mpas
                - name: user
                  value: baas
                - name: passwd
                  value: cdiWappZtnd6GMMu
                - name: arg
                  value: --skip-version-check
                - name: skip
                  value: false
        - - name: argocd-sync
            templateRef:
              name: argocd-sync-template
              template: main
            arguments:
              parameters:
                - name: argocd-version
                  value: v2.3.4
                - name: application-name
                  value: daojia-int-qw-mpas-service
                - name: flags
                  value: --insecure
                - name: argocd-server-address
                  value: argocd.hd123.com
                - name: argocd-rootpath
                  value: /dnet-int/
                - name: argocd-credentials-secret
                  value: argocd-secret
```

#### argo events
- 使用 webhook 触发 Workflow 部署
##### EventSource 事件源
```yaml
apiVersion: argoproj.io/v1alpha1
kind: EventSource
metadata:
  name: upgrade
  namespace: argo-events
spec:
  service:
    ports:
      - port: 12000
        targetPort: 12000
  webhook:
    # event-source can run multiple HTTP servers. Simply define a unique port to start a new HTTP server
    upgrade:
      # port to run HTTP server on
      port: "12000"
      # endpoint to listen to
      endpoint: /upgrade
      # HTTP request method to allow. In this case, only POST requests are accepted
      method: POST
```
##### Sensor 传感器
- 监控事件源并触发触发器
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Sensor
metadata:
  name: upgrade
  namespace: argo-events
spec:
  template:
    serviceAccountName: operate-workflow-sa
  dependencies:
    - name: test
      eventSourceName: upgrade
      eventName: upgrade
  triggers: # 触发器
    - template:
        name: webhook-workflow-trigger
        k8s:
          operation: create
          source:
            resource:
              apiVersion: argoproj.io/v1alpha1
              kind: Workflow
              metadata:
                generateName: mpas-rdb-upgrade-
              spec:
                entrypoint: main
                templates:
                  - name: main
                    steps:
                      - - name: rdb-upgrade
                          templateRef:
                            name: rdb-upgrade-template
                            template: main
                          arguments:
                            parameters:
                              - name: image
                                value: harbor.qianfan123.com/baas/mpas-rdb-upgrade:1.2.4-SNAPSHOT
                              - name: datasource
                                value: rm-bp1496f14i5b9xaod.mysql.rds.aliyuncs.com:3306/mpas
                              - name: user
                                value: baas
                              - name: passwd
                                value: cdiWappZtnd6GMMu
                              - name: arg
                                value: --skip-version-check
                              - name: skip
                                value: false
                      - - name: argocd-sync
                          templateRef:
                            name: argocd-sync-template
                            template: main
                          arguments:
                            parameters:
                              - name: argocd-version
                                value: v2.3.4
                              - name: application-name
                                value: daojia-int-qw-mpas-service
                              - name: flags
                                value: --insecure
                              - name: argocd-server-address
                                value: argocd.hd123.com
                              - name: argocd-rootpath
                                value: /dnet-int/
                              - name: argocd-credentials-secret
                                value: argocd-secret
```

```json
curl -d '{}' -H "Content-Type: application/json" -X POST http://192.168.55.159:12000/upgrade
```
##### 使用webhook传递参数
```yaml
apiVersion: argoproj.io/v1alpha1
kind: Sensor
metadata:
  name: upgrade
  namespace: argo-events
spec:
  template:
    serviceAccountName: operate-workflow-sa
  dependencies:
    - name: test
      eventSourceName: upgrade
      eventName: upgrade
  triggers:
    - template:
        name: webhook-workflow-trigger
        k8s:
          operation: create
          source:
            resource:
              apiVersion: argoproj.io/v1alpha1
              kind: Workflow
              metadata:
                generateName: mpas-rdb-upgrade-
              spec:
                entrypoint: main
                arguments:
                  parameters:
                  - name: image
                    value: harbor.qianfan123.com/baas/mpas-rdb-upgrade:1.2.4-SNAPSHOT
                  - name: datasource
                    value: rm-bp1496f14i5b9xaod.mysql.rds.aliyuncs.com:3306/mpas
                  - name: user
                    value: baas
                  - name: passwd
                    value: cdiWappZtnd6GMMu
                  - name: arg
                    value: --skip-version-check
                  - name: skip
                    value: false
                  - name: argocd-version
                    value: v2.3.4
                  - name: application-name
                    value: daojia-int-qw-mpas-service
                  - name: flags
                    value: --insecure
                  - name: argocd-server-address
                    value: argocd.hd123.com
                  - name: argocd-rootpath
                    value: /dnet-int/
                  - name: argocd-credentials-secret
                    value: argocd-secret 
                templates:
                  - name: main
                    inputs:
                      parameters:
                      - name: image
                      - name: datasource
                      - name: user
                      - name: passwd
                      - name: arg
                      - name: skip
                      - name: argocd-version
                      - name: application-name
                      - name: flags
                      - name: argocd-server-address
                      - name: argocd-rootpath
                      - name: argocd-credentials-secret
                    steps:
                      - - name: rdb-upgrade
                          templateRef:
                            name: rdb-upgrade-template
                            template: main
                            clusterScope: true
                          arguments:
                            parameters:
                              - name: image
                                value: "{{inputs.parameters.image}}"
                              - name: datasource
                                value: "{{inputs.parameters.datasource}}"
                              - name: user
                                value: "{{inputs.parameters.user}}"
                              - name: passwd
                                value: "{{inputs.parameters.passwd}}"
                              - name: arg
                                value: "{{inputs.parameters.arg}}"
                              - name: skip
                                value: "{{inputs.parameters.skip}}"
                      - - name: argocd-sync
                          templateRef:
                            name: argocd-sync-template
                            template: main
                            clusterScope: true
                          arguments:
                            parameters:
                              - name: argocd-version
                                value: "{{inputs.parameters.argocd-version}}"
                              - name: application-name
                                value: "{{inputs.parameters.application-name}}"
                              - name: flags
                                value: "{{inputs.parameters.flags}}"
                              - name: argocd-server-address
                                value: "{{inputs.parameters.argocd-server-address}}"
                              - name: argocd-rootpath
                                value: "{{inputs.parameters.argocd-rootpath}}"
                              - name: argocd-credentials-secret
                                value: "{{inputs.parameters.argocd-credentials-secret}}"
          parameters:
            - src:
                dependencyName: test
                dataKey: body.skip
              dest: spec.arguments.parameters.5.value
```
- 传递参数
```json
curl -d '{"skip":"true"}' -H "Content-Type: application/json" -X POST http://192.168.55.159:12000/upgrade
```



