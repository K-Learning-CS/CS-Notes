# Exporter


## 集群概览
![](../../../imgs/QQ20210929-114522@2x.png)

## 一、基础资源

### Node-exporter

~~~bash
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-exporter
  labels:
    name: node-exporter
spec:
  selector:
    matchLabels:
     name: node-exporter
  template:
    metadata:
      labels:
        name: node-exporter
    spec:
      hostPID: true
      hostIPC: true
      hostNetwork: true
      containers:
      - name: node-exporter
        image: prom/node-exporter:latest
        ports:
        - containerPort: 9100
        resources:
          requests:
            cpu: 0.15
        securityContext:
          privileged: true
        args:
        - --path.procfs
        - /host/proc
        - --path.sysfs
        - /host/sys
        - --collector.filesystem.ignored-mount-points
        - '"^/(sys|proc|dev|host|etc)($|/)"'
        volumeMounts:
        - name: dev
          mountPath: /host/dev
        - name: proc
          mountPath: /host/proc
        - name: sys
          mountPath: /host/sys
        - name: rootfs
          mountPath: /rootfs
      tolerations:
      - key: "node-role.kubernetes.io/master"
        operator: "Exists"
        effect: "NoSchedule"
      volumes:
        - name: proc
          hostPath:
            path: /proc
        - name: dev
          hostPath:
            path: /dev
        - name: sys
          hostPath:
            path: /sys
        - name: rootfs
          hostPath:
            path: /
~~~

## 二、集群资源

### kube-state-metrics
~~~bash
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-state-metrics
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kube-state-metrics
  template:
    metadata:
      annotations:
        ack.aliyun.com: c35166695291649498d2d18153b3cbba0
        prometheus.io/path: /metrics
        prometheus.io/port: "8080"
        prometheus.io/scheme: http
        prometheus.io/scrape: "true"
      labels:
        app: kube-state-metrics
    spec:
      serviceAccountName: kube-state-metrics
      containers:
      - name: kube-state-metrics
        image: quay.io/coreos/kube-state-metrics:latest
        ports:
        - containerPort: 8080

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube-state-metrics

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kube-state-metrics
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "resourcequotas", "replicationcontrollers", "limitranges", "persistentvolumeclaims", "persistentvolumes", "namespaces", "endpoints"]
  verbs: ["list", "watch"]
- apiGroups: ["apps"]
  resources: ["daemonsets", "deployments", "replicasets","statefulsets"]
  verbs: ["list", "watch"]
- apiGroups: ["batch"]
  resources: ["cronjobs", "jobs"]
  verbs: ["list", "watch"]
- apiGroups: ["autoscaling"]
  resources: ["horizontalpodautoscalers"]
  verbs: ["list", "watch"]
- apiGroups: ["storage.k8s.io"]
  resources: ["storageclasses"]
  verbs: ["list", "watch"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kube-state-metrics
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kube-state-metrics
subjects:
- kind: ServiceAccount
  name: kube-state-metrics

---
apiVersion: v1
kind: Service
metadata:
  annotations:
    prometheus.io/scrape: 'true'
  name: kube-state-metrics
  labels:
    app: kube-state-metrics
spec:
  ports:
  - name: kube-state-metrics
    port: 8080
    protocol: TCP
  selector:
    app: kube-state-metrics


~~~
## 三、数据库

### mysqld-exporter

```bash
cat > .my.cnf <<EOF
[client]
user='root'
password='66666666'
host='localhost'
port='3306'
EOF

docker run -d \
  --name mysql_exporter \
  -p 9104:9104 \
  --restart always \
  -v /root/export/.my.cnf:/.my.cnf \
  prom/mysqld-exporter
  
```

### elasticsearch-exporter

```bash
# --es.uri=method://user:password@host:port

docker run -d \
 --name elasticsearch_exporter \
 -p 9114:9114 --restart always \
 quay.io/prometheuscommunity/elasticsearch-exporter:latest \
 --es.uri=http://elastic:M2ZjBkN2Z@192.168.1.130:9200 \
 --es.all --es.indices --es.indices_settings \
 --es.indices_mappings --es.shards --es.snapshots --es.timeout=30s
```

### redis-exporter

```bash
docker run -d \
  --name redis_exporter \
  -p 9121:9121 \
  --restart always \
  oliver006/redis_exporter \
  --redis.addr=192.168.108.73:6379 \
  --redis.password=j0MzY3ZGE0Mz0 
```

## 四、中间件

### kafka-exporter

```bash
docker run -d \
  --name kafka_exporter \
  -p 9308:9308 \
  --restart always \
  danielqsj/kafka-exporter \
  --kafka.server=
```

## 五、应用

### pve-exporter

```bash
cat > /root/pve.yml <<EOF
default:
    user: root@pam
    password: Rsj66666666
    # Optional: set to false to skip SSL/TLS verification
    verify_ssl: false
EOF

docker run --name prometheus-pve-exporter -d -p 0.0.0.0:9221:9221 \
 -v /root/pve.yml:/etc/pve.yml prompve/prometheus-pve-exporter
 
# 获取metrics
curl http://<exporter_IP>:9221/pve?target=<pve_IP:port>
```