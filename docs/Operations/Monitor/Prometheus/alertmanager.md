# alertmanager


## 一、组件部署

~~~bash
apiVersion: apps/v1
kind: Deployment
metadata:
  name: alertmanager
  labels:
    app: alertmanager
spec:
  replicas: 1
  selector:
    matchLabels:
      app: alertmanager
  template:
    metadata:
      labels:
        app: alertmanager
      annotations:
        prometheus.io/scrape: 'false'
    spec:
      containers:
      - name: alertmanager
        image: prom/alertmanager:latest
        args:
        - "--config.file=/etc/alertmanager/alertmanager.yml"
        - "--log.level=debug"
        ports:
        - containerPort: 9093
          protocol: TCP
          name: alertmanager
        volumeMounts:
        - name: alertmanager-config
          mountPath: /etc/alertmanager
        - name: localtime
          mountPath: /etc/localtime
      volumes:
        - name: alertmanager-config
          configMap:
            name: alertmanager
        - name: localtime
          hostPath:
           path: /usr/share/zoneinfo/Asia/Shanghai

---
apiVersion: v1
kind: Service
metadata:
  labels:
    name: alertmanager
    kubernetes.io/cluster-service: 'true'
  name: alertmanager
spec:
  ports:
  - name: alertmanager
    nodePort: 30066
    port: 9093
    protocol: TCP
    targetPort: 9093
  selector:
    app: alertmanager
  type: NodePort
~~~

## 二、配置
~~~
global:
  resolve_timeout: 1m
route:
  receiver: 'webhook'
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  group_by: [alertname]
  routes:
  - receiver: webhook
    group_wait: 10s
    match:
      team: node
receivers:
- name: 'webhook'
  webhook_configs:
  - url: http://dnet-dingtalk-webhook.dnet-prometheus.svc.cluster.local:8060/dingtalk/webhook/send
    send_resolved: true

~~~

## 三、dingtalk
~~~bash
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dingtalk-webhook
  labels:
    app: dingtalk-webhook
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dingtalk-webhook
  template:
    metadata:
      labels:
        app: dingtalk-webhook
      annotations:
    spec:
      imagePullSecrets:
        - name: "harbor"
      containers:
      - name: dingtalk-webhook
        image: harbor.qianfan123.com/prometheus/prometheus-webhook-dingtalk:1.4.1
        ports:
        - containerPort: 8060
          protocol: TCP
          name: dingtalk
        volumeMounts:
          - name: localtime
            mountPath: /etc/localtime
          - name: config
            mountPath: /etc/prometheus-webhook-dingtalk/config.yml
            subPath: config.yml
          - name: template
            mountPath: /etc/prometheus-webhook-dingtalk/templates/default.tmpl
            subPath: default.tmpl
      volumes:
        - name: localtime
          hostPath:
           path: /usr/share/zoneinfo/Asia/Shanghai
        - name: config
          configMap:
            name: dingtalk-webhook
        - name: template
          configMap:
            name: dingtalk-webhook-template

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: dingtalk-webhook
data:
  config.yml: |
    templates:
    - /etc/prometheus-webhook-dingtalk/templates/*.tmpl
    targets:
      webhook:
        url: https://oapi.dingtalk.com/robot/send?access_token=25a52ce7ad530c1325569fd714c4b3b1cb83504520fa46635c546c5a32871a64
        message:
          text: '{{ template "ding.link.content.erp" . }}'
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: dingtalk-webhook-template
data:
  default.tmpl: |
    {{ define "__text_alert_list1" }}{{ range . }}主机:{{ .Labels.server }}
    级别:{{ .Labels.severity }}
    详情:{{ .Annotations.description }}
    时间:{{ (.StartsAt.Add 28800e9).Format "2006-01-02 15:04:05" }}
    ====================
    {{ end }}{{ end }}
    {{ define "__text_alert_list2" }}{{ range . }}主机:{{ .Labels.server }}
    级别:{{ .Labels.severity }}
    详情:{{ .Annotations.description }}
    时间:{{ (.StartsAt.Add 28800e9).Format "2006-01-02 15:04:05" }}
    恢复:{{ (.EndsAt.Add 28800e9).Format "2006-01-02 15:04:05" }}
    ====================
    {{ end }}{{ end }}
    {{ define "ding.link.content.erp" }}{{ if gt (len .Alerts.Firing) 0 }}[{{ .Alerts.Firing | len }}]【告警通知】
    {{ template "__text_alert_list1" .Alerts.Firing }}{{ end }} {{ if gt (len .Alerts.Resolved) 0 }}[{{ .Alerts.Resolved | len }}]【恢复通知】
    {{ template "__text_alert_list2" .Alerts.Resolved }}{{ end }}{{ end }}
---
apiVersion: v1
kind: Service
metadata:
  name: dingtalk-webhook
  labels:
    name: dingtalk-webhook
spec:
  selector:
    app: dingtalk-webhook
  ports:
  - name: dingtalk
    port: 8060

~~~
