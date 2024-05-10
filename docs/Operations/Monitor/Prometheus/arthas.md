# Arthas




[官方文档](https://arthas.aliyun.com/zh-cn/)

## 一、服务端

### 下载jar包
```
wget https://github.com/alibaba/arthas/releases/download/arthas-all-3.5.5/arthas-tunnel-server-3.5.5-fatjar.jar
```

### 构建镜像

```
cat > Dockerfile <<'EOF'
FROM openjdk:11.0.13-slim
COPY arthas-tunnel-server-3.5.5-fatjar.jar /opt/
EXPOSE 8080/tcp
WORKDIR /opt
ENTRYPOINT "/bin/sh" "-c" "java -jar arthas-tunnel-server-3.5.5-fatjar.jar"
EOF

docker build -t arthas-tunnel-server:v3.5.5 .
```

### 部署至k8s

```
cat > arthas-tunnel-server.yaml <<'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: arthas-tunnel-server
spec:
  replicas: 1
  selector:
    matchLabels:
      name: arthas-tunnel-server
  template:
    metadata:
      labels:
        name: arthas-tunnel-server
    spec:
      imagePullSecrets:
        - name: "harbor"
      containers:
        - name: arthas-tunnel-server
          image: harbor.qianfan123.com/toolset/arthas-tunnel-server:v3.5.5
          ports:
            - containerPort: 7777
              protocol: TCP
            - containerPort: 8080
              protocol: TCP
          resources:
            limits:
              memory: 1200Mi
              cpu: 1000m
            requests:
              memory: 1000Mi
              cpu: 10m
---
apiVersion: v1
kind: Service
metadata:
  name: arthas-tunnel-server
  labels:
    name: arthas-tunnel-server
spec:
  selector:
    name: arthas-tunnel-server
  ports:
  - name: tunnel-server
    port: 8080
    nodePort: 30832
  - name: arthas
    port: 7777
    nodePort: 30852
  type: NodePort
EOF

kubectl apply -f arthas-tunnel-server.yaml
```

### 通过slb将端口映射至外网

```
30832 --> 30832
30852 --> 7777
```

## 二、agent

### 进入pod

```
kubectl exec -it <pod> -c <container> -n <ns> -- bash
```

### 下载arthas

```
curl -O https://arthas.aliyun.com/arthas-boot.jar
```

### 安装依赖

```
rm -fr /var/cache/apk
mkdir -p mkdir /var/cache/apk
apk update
apk add gcc
```

### 运行

```
java -jar arthas-boot.jar  --tunnel-server 'ws://arthas-tunnel-server:7777/ws' --target-ip 0.0.0.0 --app-name <app name>
```

## 三、访问

### 浏览器打开

```
ip:30832/apps.html
```

### 操作参考官网
