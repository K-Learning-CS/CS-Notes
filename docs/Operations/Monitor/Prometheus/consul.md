

### 注册至consul
```bash
# node
 curl -X PUT -d '{"id": "192.168.108.110:9100","name": "node-exporter","address": "192.168.108.110","port": ''9100, "checks": [{"http": "http://192.168.108.110:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register

# es
 curl -X PUT -d '{"id": "192.168.1.130:9114","name": "node-exporter","address": "192.168.1.130","port": ''9114, "checks": [{"http": "http://192.168.1.130:9114/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
# mysql
 curl -X PUT -d '{"id": "192.168.1.81:9104","name": "mysql-exporter","address": "192.168.1.81","port": ''9104, "checks": [{"http": "http://192.168.1.81:9104/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register

# redis
 curl -X PUT -d '{"id": "192.168.108.73:9121","name": "redis-exporter","address": "192.168.108.73","port": ''9121, "checks": [{"http": "http://192.168.108.73:9121/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register


# gitlab
 curl -X PUT -d '{"id": "gitlab-exporter","name": "gitlab-exporter","tags": ["prometheus"],"address": "https://gitlab.rongshujia.net","port": 0,"Meta": {"url": "/-/metrics","params":"biBietRb6YDwNj_2_PAE"}},"check": {"id": "gitlab-exporter-check","name": "GitLab Exporter HTTP Check","hTTP": "https://gitlab.rongshujia.net/-/metrics?token=biBietRb6YDwNj_2_PAE","method": "GET","interval": "10s","timeout": "1s"}}' http://192.168.108.93:32685/v1/agent/service/register

# pve
curl -X PUT -d '{"id": "192.168.108.91:9221","name": "pve-exporter","address": "192.168.108.91","port": '9221', "Meta": {"url": "pve","params":"192.168.1.126:8006"}},"checks": [{"http": "http://192.168.108.91:9221/pve#target=192.168.1.126:8006","interval": "30s"}]}'     http://192.168.108.93:32685/v1/agent/service/register```

### 删除consul注册
```bash
- 需要按列表顺序删除
curl --request PUT http://192.168.108.93:32685/v1/agent/service/deregister/<ID>
# curl --request PUT http://192.168.108.93:32685/v1/agent/service/deregister/192.168.108.110

```


```bash
 curl -X PUT -d '{"id": "192.168.108.112:9100","name": "node-exporter","address": "192.168.108.112","port": ''9100, "checks": [{"http": "http://192.168.108.112:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.111:9100","name": "node-exporter","address": "192.168.108.111","port": ''9100, "checks": [{"http": "http://192.168.108.111:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.121:9100","name": "node-exporter","address": "192.168.1.121","port": ''9100, "checks": [{"http": "http://192.168.1.121:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.122:9100","name": "node-exporter","address": "192.168.1.122","port": ''9100, "checks": [{"http": "http://192.168.1.122:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.123:9100","name": "node-exporter","address": "192.168.1.123","port": ''9100, "checks": [{"http": "http://192.168.1.123:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.124:9100","name": "node-exporter","address": "192.168.1.124","port": ''9100, "checks": [{"http": "http://192.168.1.124:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.130:9100","name": "node-exporter","address": "192.168.1.130","port": ''9100, "checks": [{"http": "http://192.168.1.130:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.131:9100","name": "node-exporter","address": "192.168.1.131","port": ''9100, "checks": [{"http": "http://192.168.1.131:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.133:9100","name": "node-exporter","address": "192.168.1.133","port": ''9100, "checks": [{"http": "http://192.168.1.133:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.134:9100","name": "node-exporter","address": "192.168.1.134","port": ''9100, "checks": [{"http": "http://192.168.1.134:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.135:9100","name": "node-exporter","address": "192.168.1.135","port": ''9100, "checks": [{"http": "http://192.168.1.135:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.150:9100","name": "node-exporter","address": "192.168.1.150","port": ''9100, "checks": [{"http": "http://192.168.1.150:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.211:9100","name": "node-exporter","address": "192.168.1.211","port": ''9100, "checks": [{"http": "http://192.168.1.211:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.215:9100","name": "node-exporter","address": "192.168.1.215","port": ''9100, "checks": [{"http": "http://192.168.1.215:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.52:9100","name": "node-exporter","address": "192.168.1.52","port": ''9100, "checks": [{"http": "http://192.168.1.52:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.66:9100","name": "node-exporter","address": "192.168.1.66","port": ''9100, "checks": [{"http": "http://192.168.1.66:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.69:9100","name": "node-exporter","address": "192.168.1.69","port": ''9100, "checks": [{"http": "http://192.168.1.69:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.81:9100","name": "node-exporter","address": "192.168.1.81","port": ''9100, "checks": [{"http": "http://192.168.1.81:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.82:9100","name": "node-exporter","address": "192.168.1.82","port": ''9100, "checks": [{"http": "http://192.168.1.82:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.83:9100","name": "node-exporter","address": "192.168.1.83","port": ''9100, "checks": [{"http": "http://192.168.1.83:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.89:9100","name": "node-exporter","address": "192.168.1.89","port": ''9100, "checks": [{"http": "http://192.168.1.89:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.200:9100","name": "node-exporter","address": "192.168.108.200","port": ''9100, "checks": [{"http": "http://192.168.108.200:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.202:9100","name": "node-exporter","address": "192.168.108.202","port": ''9100, "checks": [{"http": "http://192.168.108.202:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.40:9100","name": "node-exporter","address": "192.168.108.40","port": ''9100, "checks": [{"http": "http://192.168.108.40:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.41:9100","name": "node-exporter","address": "192.168.108.41","port": ''9100, "checks": [{"http": "http://192.168.108.41:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.42:9100","name": "node-exporter","address": "192.168.108.42","port": ''9100, "checks": [{"http": "http://192.168.108.42:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.43:9100","name": "node-exporter","address": "192.168.108.43","port": ''9100, "checks": [{"http": "http://192.168.108.43:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.61:9100","name": "node-exporter","address": "192.168.108.61","port": ''9100, "checks": [{"http": "http://192.168.108.61:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.72:9100","name": "node-exporter","address": "192.168.108.72","port": ''9100, "checks": [{"http": "http://192.168.108.72:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.73:9100","name": "node-exporter","address": "192.168.108.73","port": ''9100, "checks": [{"http": "http://192.168.108.73:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.74:9100","name": "node-exporter","address": "192.168.108.74","port": ''9100, "checks": [{"http": "http://192.168.108.74:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.90:9100","name": "node-exporter","address": "192.168.108.90","port": ''9100, "checks": [{"http": "http://192.168.108.90:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.108.91:9100","name": "node-exporter","address": "192.168.108.91","port": ''9100, "checks": [{"http": "http://192.168.108.91:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.111.11:9100","name": "node-exporter","address": "192.168.111.11","port": ''9100, "checks": [{"http": "http://192.168.111.11:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.111.12:9100","name": "node-exporter","address": "192.168.111.12","port": ''9100, "checks": [{"http": "http://192.168.111.12:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.111.13:9100","name": "node-exporter","address": "192.168.111.13","port": ''9100, "checks": [{"http": "http://192.168.111.13:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.111.14:9100","name": "node-exporter","address": "192.168.111.14","port": ''9100, "checks": [{"http": "http://192.168.111.14:9100/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
 curl -X PUT -d '{"id": "192.168.1.82:9104","name": "mysql-exporter","address": "192.168.1.82","port": ''9104, "checks": [{"http": "http://192.168.1.82:9104/","interval": "5s"}]}'     http://192.168.108.93:32685/v1/agent/service/register
```