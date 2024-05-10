# prometheus


## 一、应用部署

[启动参数参考](https://prometheus.io/docs/prometheus/latest/command-line/prometheus/)

~~~bash
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus-server
  labels:
    app: prometheus
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus
  template:
    metadata:
      labels:
        app: prometheus
      annotations:
        prometheus.io/scrape: 'false'
    spec:
      serviceAccountName: prometheus
      containers:
      - name: prometheus
        image: prom/prometheus:latest
        imagePullPolicy: IfNotPresent
        command:
          - prometheus
          - --config.file=/etc/prometheus/prometheus.yml
          - --storage.tsdb.path=/prometheus
          - --storage.tsdb.retention=72h
        ports:
        - containerPort: 9090
          protocol: TCP
        volumeMounts:
        - mountPath: /etc/prometheus
          name: prometheus-config
      volumes:
        - name: prometheus-config
          configMap:
            name: prometheus-config
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: prometheus
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: prometheus
rules:
- apiGroups:
  - ""
  resources:
  - nodes
  - services
  - endpoints
  - pods
  - nodes/proxy
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - "extensions"
  resources:
    - ingresses
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - configmaps
  - nodes/metrics
  verbs:
  - get
- nonResourceURLs:
  - /metrics
  verbs:
  - get
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: prometheus
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: prometheus
subjects:
- kind: ServiceAccount
  name: prometheus
---
apiVersion: v1
kind: Service
metadata:
  name: prometheus
  labels:
    app: prometheus
spec:
  type: NodePort
  ports:
    - port: 9090
      targetPort: 9090
      protocol: TCP
      nodePort: 30090
  selector:
    app: prometheus

~~~


## 二、数据采集存储

~~~bash
global:
  scrape_interval: 60s
  scrape_timeout: 20s
  evaluation_interval: 1m
rule_files:
- /etc/prometheus/rules.yml
alerting:
  alertmanagers:
  - static_configs:
    - targets: ["dnet-alertmanager.dnet-prometheus.svc.cluster.local:9093"]
remote_write:
  - url: "https://prometheus-dev.cn-hangzhou.log.aliyuncs.com/prometheus/prometheus-dev/prometheus-dev-2/api/v1/write?project={{customer}}"
    basic_auth:
      username: ************
      password: ************
scrape_configs:
- job_name: 'kubernetes-node'
  kubernetes_sd_configs:
  - role: node
  relabel_configs:
  - source_labels: [__address__]
    regex: '(.*):10250'
    replacement: '${1}:9100'
    target_label: __address__
    action: replace
  - action: labelmap
    regex: __meta_kubernetes_node_label_(.+)
- job_name: 'kubernetes-node-cadvisor'
  kubernetes_sd_configs:
  - role:  node
  scheme: https
  tls_config:
    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
  bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
  relabel_configs:
  - action: labelmap
    regex: __meta_kubernetes_node_label_(.+)
  - target_label: __address__
    replacement: kubernetes.default.svc:443
  - source_labels: [__meta_kubernetes_node_name]
    regex: (.+)
    target_label: __metrics_path__
    replacement: /api/v1/nodes/${1}/proxy/metrics/cadvisor
- job_name: 'kubernetes-apiserver'
  kubernetes_sd_configs:
  - role: endpoints
  scheme: https
  tls_config:
    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
  bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
  relabel_configs:
  - source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
    action: keep
    regex: default;kubernetes;https
- job_name: 'kubernetes-pods'
  kubernetes_sd_configs:
  - role: pod
  relabel_configs:
  - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
    action: keep
    regex: true
  - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
    action: replace
    target_label: __metrics_path__
    regex: (.+)
  - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
    action: replace
    regex: ([^:]+)(?::\d+)?;(\d+)
    replacement: $1:$2
    target_label: __address__
  - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scheme]
    action: replace
    target_label: __scheme__
    regex: (.+)
  - action: labelmap
    regex: __meta_kubernetes_pod_annotation_(.+)

~~~


## 三、告警规则
~~~bash
groups:
- name: 物理节点状态-监控告警
  rules:
  - alert: 物理节点cpu使用率
    expr: 100-avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by(instance)*100 > 90
    for: 2s
    labels:
      severity: ccritical
    annotations:
      summary: "千帆节点：{{ $labels.instance }}cpu使用率过高"
      description: "千帆节点：{{ $labels.instance }}的cpu使用率超过90%,当前使用率[{{ $value }}]"
  - alert: 物理节点内存使用率
    expr: (node_memory_MemTotal_bytes - (node_memory_MemFree_bytes + node_memory_Buffers_bytes + node_memory_Cached_bytes)) / node_memory_MemTotal_bytes * 100 > 85
    for: 2s
    labels:
      severity: critical
    annotations:
      summary: "千帆节点：{{ $labels.instance }}内存使用率过高"
      description: "千帆节点：{{ $labels.instance }}的内存使用率超过85%,当前使用率[{{ $value }}]"
  - alert: InstanceDown
    expr: up == 0
    for: 2s
    labels:
      severity: critical
    annotations:
      summary: "千帆节点：{{ $labels.instance }}: 服务器宕机"
      description: "千帆节点：{{ $labels.instance }}: 服务器延时超过2分钟"
  - alert: 物理节点磁盘的IO性能
    expr: 100-(avg(irate(node_disk_io_time_seconds_total[1m])) by(instance)* 100) < 60
    for: 2s
    labels:
      severity: critical
    annotations:
      summary: "千帆节点：{{$labels.mountpoint}} 流入磁盘IO使用率过高！"
      description: "千帆节点：{{$labels.mountpoint }} 流入磁盘IO大于60%(目前使用:{{$value}})"
  - alert: 入网流量带宽
    expr: ((sum(rate (node_network_receive_bytes_total{device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[5m])) by (instance)) / 100) > 102400
    for: 2s
    labels:
      severity: critical
    annotations:
      summary: "千帆节点：{{$labels.mountpoint}} 流入网络带宽过高！"
      description: "千帆节点：{{$labels.mountpoint }}流入网络带宽持续5分钟高于100M. RX带宽使用率{{$value}}"
  - alert: 出网流量带宽
    expr: ((sum(rate (node_network_transmit_bytes_total{device!~'tap.*|veth.*|br.*|docker.*|virbr*|lo*'}[5m])) by (instance)) / 100) > 102400
    for: 2s
    labels:
      severity: critical
    annotations:
      summary: "千帆节点：{{$labels.mountpoint}} 流出网络带宽过高！"
      description: "千帆节点：{{$labels.mountpoint }}流出网络带宽持续5分钟高于100M. RX带宽使用率{{$value}}"
  - alert: TCP会话
    expr: node_netstat_Tcp_CurrEstab > 1000
    for: 2s
    labels:
      severity: critical
    annotations:
      summary: "千帆节点：{{$labels.mountpoint}} TCP_ESTABLISHED过高！"
      description: "千帆节点：{{$labels.mountpoint }} TCP_ESTABLISHED大于1000%(目前使用:{{$value}}%)"
  - alert: 磁盘容量
    expr: 100-(node_filesystem_free_bytes{fstype=~"ext4|xfs"}/node_filesystem_size_bytes {fstype=~"ext4|xfs"}*100) > 80
    for: 2s
    labels:
      severity: critical
    annotations:
      summary: "千帆节点：{{$labels.mountpoint}} 磁盘分区使用率过高！"
      description: "千帆节点：{{$labels.mountpoint }} 磁盘分区使用大于80%(目前使用:{{$value}}%)"
- name: 集群告警
  rules:
  - alert: coredns的cpu使用率大于90%
    expr: rate(process_cpu_seconds_total{k8s_app=~"kube-dns"}[2m]) * 100 > 90
    for: 2s
    labels:
      severity: critical
    annotations:
      description: "千帆集群：{{$labels.instance}}的{{$labels.k8s_app}}组件的cpu使用率超过90%"
      value: "{{ $value }}%"
      threshold: "90%"
  - alert: coredns
    expr: process_open_fds{k8s_app=~"kube-dns"}  > 1000
    for: 2s
    labels:
      severity: critical
    annotations:
      description: "千帆集群：插件{{$labels.k8s_app}}({{$labels.instance}}): 打开句柄数超过1000"
      value: "{{ $value }}"
      threshold: "1000"
  - alert: Pod_restarts
    expr: kube_pod_container_status_restarts_total{namespace=~"dnet-integration-test|dnet-branch-test"} > 0
    for: 2s
    labels:
      severity: warnning
    annotations:
      description: "千帆集群：在{{$labels.namespace}}名称空间下发现{{$labels.pod}}这个pod下的容器{{$labels.container}}被重启,这个监控指标是由{{$labels.instance}}采集的"
      value: "{{ $value }}"
      threshold: "0"
  - alert: Pod_waiting
    expr: kube_pod_container_status_waiting_reason{namespace=~"dnet-integration-test|dnet-branch-test"} == 1
    for: 2s
    labels:
      severity: warnning
    annotations:
      description: "千帆集群：空间{{$labels.namespace}}({{$labels.instance}}): 发现{{$labels.pod}}下的{{$labels.container}}启动异常等待中"
      value: "{{ $value }}"
      threshold: "1"
  - alert: Endpoint_ready
    expr: kube_endpoint_address_not_ready{namespace=~"dnet-integration-test|dnet-branch-test"} == 1
    for: 2s
    labels:
      severity: warnning
    annotations:
      description: "千帆集群：空间{{$labels.namespace}}({{$labels.instance}}): 发现{{$labels.endpoint}}不可用"
      value: "{{ $value }}"
      threshold: "1"
  - alert: core_dns_error
    expr: coredns_dns_responses_total{ack_aliyun_com="c35166695291649498d2d18153b3cbba0",rcode=~"FORMERR|SERVFAIL|NOTIMP|REFUSED"} > 0
    for: 2s
    labels:
      severity: warnning
    annotations:
      description: "千帆集群：coredns ({{$labels.instance}}) 解析失败"
      value: "{{ $value }}"
      threshold: "0"

~~~
