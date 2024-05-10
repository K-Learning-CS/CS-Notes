# Archery

## Why？

### 什么是 Archery ？

*一站式的 SQL 审核查询平台*

[官方网站](https://archerydms.com/)

### 为什么需要 Archery ？

- 数据库统一管理
- sql查询平台
- sql审核平台
- 人员权限细粒度限制

### Archery 有什么功能？

|            | 查询 | 审核 | 执行 | 备份 | 数据字典 | 慢日志 | 会话管理 | 账号管理 | 参数管理 | 数据归档 |
| :--------: | :--: | :--: | :--: | :--: | :------: | :----: | :------: | :------: | :------: | :------: |
|   MySQL    |  √   |  √   |  √   |  √   |    √     |   √    |    √     |    √     |    √     |    √     |
|   MsSQL    |  √   |  ×   |  √   |  ×   |    ×     |   ×    |    ×     |    ×     |    ×     |    ×     |
|   Redis    |  √   |  ×   |  √   |  ×   |    ×     |   ×    |    ×     |    ×     |    ×     |    ×     |
|   PgSQL    |  √   |  ×   |  √   |  ×   |    ×     |   ×    |    ×     |    ×     |    ×     |    ×     |
|   Oracle   |  √   |  √   |  √   |  √   |    ×     |   ×    |    ×     |    ×     |    ×     |    ×     |
|  MongoDB   |  √   |  √   |  √   |  ×   |    ×     |   ×    |    ×     |    ×     |    ×     |    ×     |
|  Phoenix   |  √   |  ×   |  √   |  ×   |    ×     |   ×    |    ×     |    ×     |    ×     |    ×     |
|    ODPS    |  √   |  ×   |  ×   |  ×   |    ×     |   ×    |    ×     |    ×     |    ×     |    ×     |
| ClickHouse |  √   |  ×   |  ×   |  ×   |    ×     |   ×    |    ×     |    ×     |    ×     |    ×     |

## Install

### docker-compose

*使用 docker-compose 部署 archery 1.10*

- 请提前规划数据目录！

#### 1. 安装 docker

```bash
# 安装Docker
yum install -y yum-utils device-mapper-persistent-data lvm2
yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum makecache fast
yum -y install docker-ce
systemctl enable --now docker


# 配置docker加速并修改驱动
cat > /etc/docker/daemon.json <<EOF
{
    "exec-opts": ["native.cgroupdriver=systemd"],
    "registry-mirrors": [
        "http://hub-mirror.c.163.com",
        "https://docker.mirrors.ustc.edu.cn",
        "http://f1361db2.m.daocloud.io",
        "https://registry.docker-cn.com"
    ],
	"insecure-registries" : ["http://192.168.108.2"]
}
EOF
systemctl restart docker
```

#### 2. 安装 docker-compose

```bash
curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
```

#### 3. 下载 archery 1.10 源码，并进入指定目录

```bash
# 下载源码
wget https://github.com/hhyo/Archery/archive/refs/tags/v1.10.0.tar.gz

# 解压
tar xf v1.10.0.tar.gz

# 进入工作目录
cd Archery-1.10.0/src/docker-compose/
```

#### 4. 修改、运行、初始化

```bash
# 将当前目录复制至数据目录
cp -r . /data/

# 进入数据目录
cd /data/

# 启动
docker-compose -f docker-compose.yml up -d

# 表结构初始化
docker exec -ti archery /bin/bash
cd /opt/archery
source /opt/venv4archery/bin/activate
python3 manage.py makemigrations sql  
python3 manage.py migrate 

# 数据初始化
python3 manage.py dbshell<sql/fixtures/auth_group.sql
python3 manage.py dbshell<src/init_sql/mysql_slow_query_review.sql

# 创建管理用户
python3 manage.py createsuperuser
exit

# 重启
docker restart archery

# 日志查看和问题排查
docker logs archery -f --tail=50

```

#### 5. 访问

`http://<server_ip>:9123`

### openldap

*Archery 的配置文件为 .env 文件*

*ldap 配置的具体细节请参考 ldap 中的名称定义，如有不同请修改*

*如果 ldap 用户无法登陆，请根据日志排查错误*

```bash
# 对接 ldap 修改 .env 文件
vi .env
---
# https://docs.djangoproject.com/en/4.0/ref/settings/#csrf-trusted-origins
CSRF_TRUSTED_ORIGINS='http://127.0.0.1:9123,http://192.168.108.110:9123,ldap://192.168.108.91:389'

# https://django-auth-ldap.readthedocs.io/en/latest/
ENABLE_LDAP=True
AUTH_LDAP_ALWAYS_UPDATE_USER=True
AUTH_LDAP_SERVER_URI="ldap://192.168.108.91:389"
AUTH_LDAP_BIND_DN="cn=admin,dc=rsjitcm,dc=com"
AUTH_LDAP_BIND_PASSWORD="YCAEZc2eaR6R5J"
AUTH_LDAP_USER_SEARCH=LDAPSearch('ou=users,dc=rsjitcm,dc=com',ldap.SCOPE_SUBTREE,'(uid=%(user)s)',)
AUTH_LDAP_USER_SEARCH_BASE='ou=users,dc=rsjitcm,dc=com'
```

## Configure

*使用安装时创建的 root 用户登陆*

### 资源组

[系统管理] --> [其他配置管理] --> [全部后台数据] --> [资源组配置] --> [增加资源组管理]

*资源组管理可以给数据库资源分类，例如生产、非生产，或者mysql、mongodb*

### 数据库实例

[实例管理] --> [实例列表] --> [添加实例] --> [资源组] --> [实例标签]

*资源组决定了数据库在哪个分类中*

*实例标签决定了数据库能否进行上线和查询，请按照需求添加标签*

### 用户组管理

[系统管理] --> [其他配置管理] --> [权限组管理]

*用户组即给组赋权限，然后将用户加入其中，便可以获取组的权限*

### 用户管理

[系统管理] --> [其他配置管理] --> [用户管理]

### 审核流程管理

[系统管理] --> [配置项管理] --> [配置项] --> [工单审核流程配置]

*选择多个权限组即审批流程为多级审核，按照选择顺序进行流转，权限组内用户都可审核*

### 系统设置

[系统管理] --> [配置项管理] --> [配置项] --> [系统设置] 

#### goInception配置

*在安装 archery 时默认有安装 goInception，按照安装的配置即可*

#### 脱敏配置
