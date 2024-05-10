# Adapter

`将prometheus的指标转化为k8s接口`



[官网](https://github.com/kubernetes-sigs/prometheus-adapter)


## 一、示例应用
~~~bash
apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpserver
spec:
  replicas: 1
  selector:
    matchLabels:
      app: httpserver
  template:
    metadata:
      labels:
        app: httpserver
    spec:
      containers:
      - name: httpserver
        image: imroc.tencentcloudcr.com/test/httpserver:v1
        imagePullPolicy: Always
        ports:
          - containerPort: 80
            protocol: TCP
---

apiVersion: v1
kind: Service
metadata:
  name: httpserver
  labels:
    app: httpserver
  annotations:
    prometheus.io/scrape: "true"
spec:
  type: ClusterIP
  ports:
  - port: 80
    protocol: TCP
    name: http
  selector:
    app: httpserver
~~~

- 这里会取到如下指标

~~~bash
httpserver_requests_total{app="httpserver", instance="172.16.161.217:80", job="kubernetes-service-endpoints", kind="Pod", kubernetes_name="httpserver", kubernetes_namespace="default", name="httpserver-8454cb5fc-lv54t", status="200"}  66
~~~



## 二、安装adapter
- 安装参考[官方文档](https://github.com/kubernetes-sigs/prometheus-adapter/tree/master/deploy)

### RBAC

~~~bash
kind: ServiceAccount
apiVersion: v1
metadata:
  name: custom-metrics-apiserver
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: custom-metrics:system:auth-delegator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
- kind: ServiceAccount
  name: custom-metrics-apiserver
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: custom-metrics-auth-reader
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: extension-apiserver-authentication-reader
subjects:
- kind: ServiceAccount
  name: custom-metrics-apiserver
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: custom-metrics-resource-reader
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: custom-metrics-resource-reader
subjects:
- kind: ServiceAccount
  name: custom-metrics-apiserver
  namespace: default
---
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  name: v1beta1.custom.metrics.k8s.io
spec:
  service:
    name: custom-metrics-apiserver
    namespace: default
  group: custom.metrics.k8s.io
  version: v1beta1
  insecureSkipTLSVerify: true
  groupPriorityMinimum: 100
  versionPriority: 100
---
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  name: v1beta2.custom.metrics.k8s.io
spec:
  service:
    name: custom-metrics-apiserver
    namespace: default
  group: custom.metrics.k8s.io
  version: v1beta2
  insecureSkipTLSVerify: true
  groupPriorityMinimum: 100
  versionPriority: 200
---
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  name: v1beta1.external.metrics.k8s.io
spec:
  service:
    name: custom-metrics-apiserver
    namespace: default
  group: external.metrics.k8s.io
  version: v1beta1
  insecureSkipTLSVerify: true
  groupPriorityMinimum: 100
  versionPriority: 100
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: custom-metrics-server-resources
rules:
- apiGroups:
  - custom.metrics.k8s.io
  - external.metrics.k8s.io
  resources: ["*"]
  verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: custom-metrics-resource-reader
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - nodes
  - nodes/stats
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: hpa-controller-custom-metrics
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: custom-metrics-server-resources
subjects:
- kind: ServiceAccount
  name: horizontal-pod-autoscaler
  namespace: kube-system
~~~

### Deployment

- 1.`- --prometheus-url=` 填写prometheus地址

- 2.`- --metrics-relist-interval=` 间隔时间必须不大于prometheus采集指标的间隔

- 3.`TLS secret` 的生成参考[gencerts.sh](https://github.com/prometheus-operator/kube-prometheus/blob/62fff622e9900fade8aecbd02bc9c557b736ef85/experimental/custom-metrics-api/gencerts.sh)
~~~bash
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: custom-metrics-apiserver
  name: custom-metrics-apiserver
spec:
  replicas: 1
  selector:
    matchLabels:
      app: custom-metrics-apiserver
  template:
    metadata:
      labels:
        app: custom-metrics-apiserver
      name: custom-metrics-apiserver
    spec:
      serviceAccountName: custom-metrics-apiserver
      containers:
      - name: custom-metrics-apiserver
        image: registry.cn-zhangjiakou.aliyuncs.com/kpw/kubernetes:v0.9.0
        args:
        - --secure-port=6443
        - --tls-cert-file=/var/run/serving-cert/serving.crt
        - --tls-private-key-file=/var/run/serving-cert/serving.key
        - --logtostderr=true
        - --prometheus-url=http://dnet-prometheus.dnet-prometheus.svc.cluster.local:9090/
        - --metrics-relist-interval=1m
        - --v=10
        - --config=/etc/adapter/config.yaml
        ports:
        - containerPort: 6443
        volumeMounts:
        - mountPath: /var/run/serving-cert
          name: volume-serving-cert
          readOnly: true
        - mountPath: /etc/adapter/
          name: config
          readOnly: true
        - mountPath: /tmp
          name: tmp-vol
      volumes:
      - name: volume-serving-cert
        secret:
          secretName: cm-adapter-serving-certs
      - name: config
        configMap:
          name: adapter-config
      - name: tmp-vol
        emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: custom-metrics-apiserver
spec:
  ports:
  - port: 443
    targetPort: 6443
  selector:
    app: custom-metrics-apiserver
---
apiVersion: v1
kind: Secret
metadata:
  name: cm-adapter-serving-certs
type: Opaque
data:
  serving.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURlVENDQW1HZ0F3SUJBZ0lVQ2V0a3BVMWJXMy9ieHVEYlF1cmN6NWlFdjdJd0RRWUpLb1pJaHZjTkFRRUwKQlFBd0RURUxNQWtHQTFVRUF3d0NZMkV3SGhjTk1qRXhNREV4TURVMU5qQXdXaGNOTWpZeE1ERXdNRFUxTmpBdwpXakFqTVNFd0h3WURWUVFERXhoamRYTjBiMjB0YldWMGNtbGpjeTFoY0dselpYSjJaWEl3Z2dFaU1BMEdDU3FHClNJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUURDRndNaTUrRmxiQytlODE2TkJVVUNjcis3cStjNDNSMWoKTS9mUTFYRExKSW80YmgyYVZlc0FydzVnMlFqOUxXa3d1U2hwaGZ3YURjT1BtRzkra09BUTFnZ2JpV0VQRmFseQo3QmVjeTBXSXhjK0k1ZVlxREpVZzkzcWI1SHNibjZhMTlOQWI0LzI1YnE2a2FGd1R5a2I0TFVacDVXM2tKRUdWClBpMW43aHl0cGZiaGhic29OekFaM0dXd3RBQkhEQmc4ZVhHRXRNRkdiRFdHRUI2Y2ZaQ0xySTd0U29HUFF0RVUKSTFyek9KVW5WWnIzQUpLaFYyUi92VjdOWCtmTGtLUHhzK3o5L2pqVlVNNEtuSHVvNTB3VDMwbnpGTjZ1S1k1ZQpOR2VqU25MNGRPaGhJTW1ZemRaenMvdTVYWStGR0RLZ1NRYzhmL01SbVpnZXh6b0ZHdmlWQWdNQkFBR2pnYm93CmdiY3dEZ1lEVlIwUEFRSC9CQVFEQWdXZ01Bd0dBMVVkRXdFQi93UUNNQUF3SFFZRFZSME9CQllFRkdDYXM1eEQKY0xUT3dVRTZtd0JTZkVRVnY4MjRNQjhHQTFVZEl3UVlNQmFBRk9hZ1gyS2N3SVlwcWJPWDc1N1gyZVNpN3FlWApNRmNHQTFVZEVRUlFNRTZDSTJOMWMzUnZiUzF0WlhSeWFXTnpMV0Z3YVhObGNuWmxjaTV0YjI1cGRHOXlhVzVuCmdpZGpkWE4wYjIwdGJXVjBjbWxqY3kxaGNHbHpaWEoyWlhJdWJXOXVhWFJ2Y21sdVp5NXpkbU13RFFZSktvWkkKaHZjTkFRRUxCUUFEZ2dFQkFGSmF3c0FoZy83UTNBbUR5RzArOWpzbXE4d1FUZ2ZmTURuR1VqWGN5UGY1USsyYwpERW5nd01LS1FmelByWjFpN3UvUkkvQVQ4YWk1MjJPU3FFU2Y1Zk5ta0Iwa1ZsclFGeDc1clVXekdoaUF3R09KClRPNkpnNHlqd2dCT0grbGlwRTBIbGdCeDZ1WEZZV1FaOUFTbHhvcW40R0VaZXQxNVdCd3VrMmRVNWs5NFBLSEgKTjhrYVpnQVd4cjhsMTFETHcxdXNXVzEyYlBlTm5LNlNrOE9Kcm8rZ3kwRUpZbEFYdGRIVTVOdXg2cVcyQmVlQQpVdXNhYTEyNUVNVlVUVnVjd3B3OTZRcGk4Mm9TbWZhanlLdW5tcDVCUUh3TUVqdEtCOUFRaHZvRC9DN2lWSVVZCk5VZlVYRUVObGRDWnZqMmNSaFVacHVaRGx1dE43YnVLQUQrbW54TT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  serving.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBd2hjREl1ZmhaV3d2bnZOZWpRVkZBbksvdTZ2bk9OMGRZelAzME5Wd3l5U0tPRzRkCm1sWHJBSzhPWU5rSS9TMXBNTGtvYVlYOEdnM0RqNWh2ZnBEZ0VOWUlHNGxoRHhXcGN1d1huTXRGaU1YUGlPWG0KS2d5VklQZDZtK1I3RzUrbXRmVFFHK1A5dVc2dXBHaGNFOHBHK0MxR2FlVnQ1Q1JCbFQ0dForNGNyYVgyNFlXNwpLRGN3R2R4bHNMUUFSd3dZUEhseGhMVEJSbXcxaGhBZW5IMlFpNnlPN1VxQmowTFJGQ05hOHppVkoxV2E5d0NTCm9WZGtmNzFlelYvbnk1Q2o4YlBzL2Y0NDFWRE9DcHg3cU9kTUU5OUo4eFRlcmltT1hqUm5vMHB5K0hUb1lTREoKbU0zV2M3UDd1VjJQaFJneW9Fa0hQSC96RVptWUhzYzZCUnI0bFFJREFRQUJBb0lCQVFDY1cyV1BiVFpMT29oeQppS1NXL3JQRmNTTzgwSk9KWDdnWS92aVpLQm1oeldIOGE5azFTQm4xaHhFU1BFWGRrQU81MkxBUnNucVJrcDBFCnhVeXNyWkdVZnBneGRzN1dGQ0ZhRDVCR0pBdDBUOGNOQmdnUnYra3prYXNZZzB3WnlOZklwZHd4Vzg0KzRFZVkKOHVtYWw4M3NpS3k5Q3JNb28zeWgrbUVoNU5UOW5kb0JlQ1oxRzFQVy9qaThNbTNOZVJodFlSTHRjYUlYU2ptYgpYOFN0UFBrMWEwdFVhMlhjVEJNWDlTd0ljZEdXRmlvSjMwSWRJVk92d2kxZFVneDFPaHRMSmNJbmdvSW8rUTFsCitOeHdMTFRHUGVoTDRZS0V4VUpVQlErTlI1TDlFZTNJR2hPMjRvalVoNklTa0pOYzlTR093d0JPbk5nRERCZEQKMGM4blV2ZlJBb0dCQU80eFBMbGlGL0hxYVl6RThSWmpHRVhNUjJBbVpwd3hEc1JsTVJhamdlcS9LaUY0WFJZegpzdnJtcVltSzVCNWpza2R1Mm9uSW1zVzMwcElLYlplTzMyUThyV293a1NTM2l1MVdLWkRNb3VxZVg1Q0xudlkvCk1IK1ZMejJkUEpxcGY3WTdmN0c4UkN0SmhVVmw3eC94UklKOXJYN1k2MS8rZUE3aUpLekVDSVEzQW9HQkFOQ1oKc3M4cVNtbVp0TUVYZHBvMjdtdS9PWDZyTDVsdW1neFhsRGtKd0lUQnd4R2U4YXpnZUI1MjAvZklnTGJpeWNscQp2U1R0NWwvbFQ4UlRzb05uMnFmUXkyZ2xvaGozdVJXdHNEYnhqc2JoOXVXUERxZVFibGllQ3Z3ZmJYR1JYUHg5ClRCV0NqVE13ZHQ2eVRRY1JTeHpkeWpqWE1ObVkwck94bjNRWHhOdVRBb0dCQU5jeGdxUVZ2SDVpQXJRY0paZk4KTlZ1eDMvWTlHMDBYZ1Rqc0Z6cFZ4SVVaNm0xTXVnVFo5bVI1U2tncVJFZzBXQmZ6VGR0WGNvVVl1MVFYdWNWSQpYZ2pJVFAvNEd0bHFQVWlKSklwZVp2M0MwYUhja21QMDJOTWJMQS9sWTZCemJCOXVoOEpDemUreHY0YmdQZmJFCjJkbHV3L1VxOHhQSjZodkFNZFFvVTIxbkFvR0FOYkZxTWlMYmxvVG0zdERRU1crY1BRV3lvZVVrVW1VQ3AreWYKRFhOeUozbk1ZU3U5WDFkRDgrdDRNZzVjK3pZeTVISmlEekJoSFF2a1ZVK0o0b01INkN3NVB5eDRwZDZWdUh2RgpvTTdhaGx6QmRXTTJUWEZDeGZLZ056ZExyM0RRTTNsNDdReDJsZGVDc1YzSnIra0dvWDZCUDlJOEU3WmZmYnRaCnBNTTllNXNDZ1lFQXduQjN0c3RQSW9HenJDSndzY0t4QllNRVpTSlhqNVRvNk92SjZwNytNNDJ1dm9qZG5paW0KWWY0NGswSzlBWmlvTmVRUmhWV0xzQlU4c2ExQ1RFeW9vTXErL0twVXNTSHM5MW1INFdjSU9yTkFmRU5uNTdLNApFbVp6SWVkMnROWVZqTVVIaXRvZ2JCOWM4MjRpZVVDZStGU3EzcEE1L21JWWptZzcvRDNQNEw0PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
~~~

### Configmap
- 配置文件是整个软件的核心，配置参考[官方文档](https://github.com/kubernetes-sigs/prometheus-adapter/tree/master/docs)

~~~bash
apiVersion: v1
kind: ConfigMap
metadata:
  name: adapter-config
data:
  config.yaml: |
    rules:
      - seriesQuery: 'httpserver_requests_total{kubernetes_namespace!="",name!=""}' # 过滤prometheus指标
        resources:
          overrides: # 使用prometheus指标的标签来对应k8s资源
            kubernetes_namespace: {resource: "namespace"}
            name: {resource: "pod"}
        name:
          matches: "httpserver_requests_total" # prometheus指标
          as: "httpserver_requests_qps" # k8s接口名称
        metricsQuery: 'sum(rate(<<.Series>>{<<.LabelMatchers>>}[2m])) by (<<.GroupBy>>)'
~~~

- 1.`metricsQuery` 是adapter去查询prometheus所使用的表达式，中间的关键字使用go模板定义

- 2.注意查看日志，可能会因为标签值不匹配导致接口请求返回没有数据，这时需要修改prometheus收集时的标签转换，获取所需标签

## 三、HPA

~~~bash
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: httpserver
spec:
  minReplicas: 1
  maxReplicas: 10
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: httpserver
  metrics:
  - type: Pods
    pods:
      metric:
        name: httpserver_requests_qps
      target:
        averageValue: '5'
        type: AverageValue
~~~

- 这里取到的值非常小时会出现`m`单位，如下

~~~bash
[root@iZbp11hp13k83iuw7pu1diZ ~]# kubectl get hpa
NAME         REFERENCE               TARGETS    MINPODS   MAXPODS   REPLICAS   AGE
httpserver   Deployment/httpserver    16m/5        1        10         1       22h
~~~

- 这是k8s风格的计算单位，意为千分之，这里由于我的刮取策略为每分钟刮去一次数据，这个指标反应的正是被刮取的次数，所以根据上面的计算表达式2/120得到0.0166无限循环小数，转化为千分之十六，即16m

