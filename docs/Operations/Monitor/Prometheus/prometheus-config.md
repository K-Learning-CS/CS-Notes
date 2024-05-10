
### 总览

- 对于prometheus的配置有以下几大块

```yaml
# 全局刮取配置，如果scrape_configs中没有再定义，则继承此配置
global:
  scrape_interval: 10s
  scrape_timeout: 5s
  evaluation_interval: 10s
---
# 告警规则配置，此为prometheus产生告警的规则，告警产生后将推送给alertmanager
rule_files:
- /path/to/rules.yml
---
# 警报器，另外的服务，为prometheus产生的告警提供接收、处理、发送的服务，此配置为该服务的地址
alerting:
  alertmanagers:
  - static_configs:
    - targets: ["alertmanager.monitor.svc.cluster.local:9093"]
---
# 刮取规则，定义从何处取到数据
scrape_configs:



# 以下为可选配置
---
# 远程写入地址，如有外部存储，则需要配置此配置
remote_write:
---
# 远程读取地址，如有外部存储，则需要配置此配置
remote_read:
---
# 存储配置，定义prometheus的数据存储位置
storage:
---
# 查询配置，定义prometheus的查询配置，如开启追踪等
tracing:
```

### scrape_configs

- 当我们使用exporter产生数据以后，就需要配置 scrape_configs 取到数据，然后存储到prometheus中

1. 使用固定的地址

    ```yaml
    scrape_configs:
    - job_name: 'prometheus'
      static_configs:
      - targets: ['localhost:9090']
    ```

2. 使用服务发现

   ```yaml
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
   ```
   
- 对于服务发现，通常我们使用服务发现的形式，此方式灵活性更高，便于使用

1. k8s中我们通常使用 `kubernetes_sd_configs` 来实现服务发现，动态的发现k8s中需要监控的组件
2. 外部的服务通畅使用 `consul_sd_configs` 来实现服务发现，动态的发现外部需要监控的组件

#### relabel_configs

- relabel_configs 可以在目标被刮取前进行修改，通常用于修改目标的标签，或者修改目标的地址

   ```yaml
   # 一个job通常由四个部分组成，job_name，服务发现方式，relabel_configs，metric_relabel_configs
   
   - job_name: 'kubernetes-service-endpoints'
     kubernetes_sd_configs:
     - role: endpoints
     relabel_configs:
       # 此处的配置为，如果目标的标签中有prometheus.io/scrape=true，则保留此目标
     - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scrape]
       action: keep
       regex: true
       # 此处的配置为，如果目标的标签中有prometheus.io/scheme，则将此标签的值赋值给__scheme__
     - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
       action: replace
       target_label: __scheme__
       regex: (https?)
       # 此处的配置为，如果目标的标签中有prometheus.io/path，则将此标签的值赋值给__metrics_path__
     - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
       action: replace
       target_label: __metrics_path__
       regex: (.+)
       # 此处的配置为，如果目标的标签中有prometheus.io/port，则将此标签的值赋值给__address__，并将__meta_kubernetes_service_annotation_prometheus_io_port标签的值赋值给__metrics_path__
     - source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
       action: replace
       target_label: __address__
       regex: ([^:]+)(?::\d+)?;(\d+)
       replacement: $1:$2
     # 
     metric_relabel_configs:
     - source_labels: [__name__]
       regex: 'node_disk_(.*)_bytes_total'
       target_label: disk
       replacement: '${1}'
   ```
- `__scheme__` `__metrics_path__` `__address__`都是 prometheus 的内置标签
- 通过组合这三个标签即可拼接处需要拉取metrics的地址：`__scheme__` 为协议，`__metrics_path__` 为路径，`__address__` 为地址

```bash
# relabel_configs 中 action 的类型有以下几种：

replace: 替换目标，默认值，匹配到label后，将目标的label替换为匹配到的label
keep: 保留目标，如果没有该label，则丢弃目标
drop: 删除目标，如果有该label，则丢弃目标
labelmap: 将目标的label映射到新的label上，标签重命名，值不变
hashmod: 根据label的值，对目标进行hash，然后根据hash的值，对目标进行分组
labeldrop: 删除目标的label
labelkeep: 保留目标的label
```

#### metric_relabel_configs

- metric_relabel_configs 在刮取目标后对数据进行处理，然后存到 prometheus，通常用来删除不需要的指标

```yaml
     metric_relabel_configs:
     - source_labels: [__name__]
       regex: 'DEMO.*'
       action: drop
```

- 也可以用来对旧标签重命名，例如老版的 exporter 指标和新版的有差别

```yaml
     - source_labels: [__name__]
       regex: 'DEMO(.*)'
       action: replace
       target_label: __name__
       replacement: '${1}'
```

- 查询最大的20个时间序列 `topk(20,count by (__name__,job)({__name__=~".+"}))`
- 如果不需要则可以在收集端进行过滤，减少 prometheus 的压力

#### consul

```yaml
- job_name: 'consul'
  consul_sd_configs:
    - server: 'consul-server:8500'
      services: []  #匹配所有service
  relabel_configs:
    - source_labels: [__meta_consul_service] #service 源标签
      regex: "consul"  #匹配为"consul" 的service
      action: drop       # 执行的动作
    - source_labels: [__meta_consul_service]  # 将service 的label重写为appname
      target_label: job
    - source_labels: [__meta_consul_service_address]
      target_label: instance
      # 以下可以根据我们自定义的metadata来进行替换，将其注册至consul后便可以在此处读取到
    - source_labels: [__meta_consul_service_metadata_url]
      action: replace
      target_label: __metrics_path__
      regex: (.+)
    - source_labels: [__meta_consul_service_metadata_params]
      action: replace
      # 将consul中的metadata中的params替换为__param_target，此处的效果为 ?target=xxx
      # 在prometheus中不可以直接将 ? 作为url，否则将会因为字符集变成%3F，导致无法访问
      target_label: __param_target
      regex: (.+)
```

```bash
# 将exporter注册至consul
curl -X PUT -d '{"id": "192.168.108.91:9221","name": "pve-exporter",
"address": "192.168.108.91","port": '9221',
"Meta": {"url": "pve","params":"192.168.1.126:8006"}},
"checks": [{"http": "http://192.168.108.91:9221/pve#target=192.168.1.126:8006","interval": "30s"}]}' \
http://192.168.108.93:32685/v1/agent/service/register

# 从consul中删除exporter
curl --request PUT http://192.168.108.93:32685/v1/agent/service/deregister/<ID>
# curl --request PUT http://192.168.108.93:32685/v1/agent/service/deregister/192.168.108.110
```

### rule_files

- 用于加载规则文件，可以是单个文件，也可以是目录，目录下的所有文件都会被加载
- 该配置项可以多次配置，用于加载多个规则文件

```yaml
# 匹配以.rules结尾的文件
rule_files:
  - "rules/*.rules"
---
# 指定文件
rule_files:
  - "first.rules"
  - "second.rules"
```

#### recording rules

- recording rules 用于将 prometheus 中的数据进行处理，然后存储到 prometheus 中，效率比外部直接使用表达式查询高，可以节约 prometheus 的资源，提升用户体验
- 通常用于计算一些复杂的指标，例如：计算某个指标的平均值、计算某个指标的百分位数等

```yaml
groups:
# 组名必须唯一
- name: 记录规则
  rules:
  # 规则可以有多条
  - record: node_cpu_seconds_total:avg_rate5m
    expr: avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by(instance)*100
    # 通过labels可以为指标添加标签，可选
    labels:
      severity: warning
```


#### alert rules

- alert rules 用于将 prometheus 中的数据进行处理，然后发送到 alertmanager，通常用于告警

```yaml

groups:
# 组名必须唯一
- name: 物理节点状态-监控告警
  rules:
  # 规则可以有多条
  - alert: 物理节点cpu使用率
    expr: 100-avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by(instance)*100 > 90
    for: 30s
    labels:
      severity: warning
    annotations:
      summary: "节点：{{ $labels.instance }}cpu使用率过高"
      description: "节点：{{ $labels.instance }}的cpu使用率超过90%,当前使用率[{{ $value }}]"
```