# Kustomize
`解决多环境配置管理`

目录
=================

* [Kustomize](#kustomize)
   * [生成配置文件](#生成配置文件)
      * [我有一个nginx服务需要管理](#我有一个nginx服务需要管理)
   * [overlays](#overlays)
      * [我需要test和prod环境](#我需要test和prod环境)
      * [如何区分环境](#如何区分环境)
         * [添加注释](#添加注释)
         * [添加标签](#添加标签)
         * [添加前后缀](#添加前后缀)
         * [为特定资源添加前后缀](#为特定资源添加前后缀)
         * [指定名称空间](#指定名称空间)
         * [指定镜像](#指定镜像)
      * [添建配置或密钥](#添建配置或密钥)
         * [在resource中引用资源](#在resource中引用资源)
         * [使用生成器生成](#使用生成器生成)
         * [添加密钥](#添加密钥)
      * [补丁](#补丁)
      * [CLI](#cli)
         * [构建](#构建)


[官方文档](https://kubectl.docs.kubernetes.io/references/kustomize/resource/)

## 生成配置文件

### 我有一个nginx服务需要管理

- deployment

~~~bash
cat deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      name: nginx
  template:
    metadata:
      labels:
        name: nginx
    spec:
      containers:
      - image: nginx:1.18
        name: nginx
        ports:
        - containerPort: 80
          protocol: TCP
~~~

- service

~~~bash
cat service.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    name: nginx
  name: nginx
spec:
  ports:
  - name: nginx
    port: 80
  selector:
    name: nginx
~~~

- 如何让kustomize识别资源

~~~bash
cat kustomization.yaml
resources:
  - deployment.yaml
  - service.yaml
~~~

- 如何生成yaml

~~~bash
kustomize build
Error: unable to find one of 'kustomization.yaml', 'kustomization.yml' or 'Kustomization' in directory '/tmp'

kustomize build new
~~~

## overlays

### 我需要test和prod环境
~~~bash
mkdir base overlays/{test,prod} -p

cp ../new/* base/

cat > overlays/test/kustomization.yaml <<'EOF'
resources:
- ../../base
EOF

# 规范
cat > overlays/test/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
EOF

cp overlays/test/kustomization.yaml overlays/prod/kustomization.yaml
~~~

### 如何区分环境

#### 添加注释

~~~bash
cat > overlays/test/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: test

EOF


cat > overlays/prod/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: production

EOF
~~~

#### 添加标签

~~~bash
cat > overlays/test/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: test
commonLabels:
  environment: test
EOF


cat > overlays/prod/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: production
commonLabels:
  environment: production
EOF
~~~

#### 添加前后缀

~~~bash
cat > overlays/test/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: test
commonLabels:
  environment: test
namePrefix: test-
nameSuffix: -suffix
EOF


cat > overlays/prod/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: production
commonLabels:
  environment: production
namePrefix: prod-
nameSuffix: -suffix
EOF
~~~

#### 为特定资源添加前后缀

~~~bash
cat > overlays/test/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: test
commonLabels:
  environment: test
namePrefix: test-
nameSuffix: -suffix
namePrefix: plat-dnet-branch-
transformers:
- cm-suffix-transformer.yaml
EOF

cat > overlays/test/cm-suffix-transformer.yaml << 'EOF'
apiVersion: builtin
kind: PrefixSuffixTransformer
metadata:
  name: customsuffixer
suffix: "-customsuffixer"
fieldSpecs:
- kind: ConfigMap
  path: metadata/name
- kind: Secret
  path: metadata/name
EOF

~~~

#### 指定名称空间

~~~bash
cat > overlays/test/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: test
commonLabels:
  environment: test
namePrefix: test-
nameSuffix: -suffix
namespace: kustomize-test
EOF


cat > overlays/prod/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: production
commonLabels:
  environment: production
namePrefix: prod-
nameSuffix: -suffix
namespace: kustomize-prod
EOF
~~~

#### 指定镜像
~~~bash

cat > overlays/test/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: test
commonLabels:
  environment: test
namePrefix: test-
nameSuffix: -suffix
namespace: kustomize-test
images:
- name: nginx
  newName: httpd
  newTag: 'latest'
EOF
~~~


### 添建配置或密钥

#### 在resource中引用资源

~~~bash
# 添加configmap
cat >base/configmap.yaml<< 'EOF'
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-index
data:
  index.html: |
    This is kustomize!
EOF

# 引用configmap
cat > base/kustomization.yaml <<'EOF'
resources:
  - deployment.yaml
  - service.yaml
  - configmap.yaml
EOF

# 挂载configmap
cat > base/deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      name: nginx
  template:
    metadata:
      labels:
        name: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.18
          ports:
            - containerPort: 80
              protocol: TCP
          volumeMounts:
            - mountPath: /usr/share/nginx/html/
              name: config
      volumes:
      - name: config
        configMap:
          defaultMode: 0644
          name: nginx-index
EOF
~~~

#### 使用生成器生成

~~~bash
# 添加configmap
cat >base/index.html<< 'EOF'
This is kustomize!
EOF

# 引用configmap
cat > base/kustomization.yaml <<'EOF'
resources:
  - deployment.yaml
  - service.yaml
configMapGenerator:
- name: nginx-index
  files:
  - index.html
EOF
~~~

- overlays生成器配置

~~~bash
cat > overlays/test/index.html << 'EOF'
This is test!
EOF

cat > overlays/test/kustomization.yaml << 'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: test
commonLabels:
  environment: test
namePrefix: test-
nameSuffix: -suffix
namespace: kustomize-test
configMapGenerator:
- name: nginx-index
  behavior: replace
  files:
  - index.html
EOF


cat > overlays/prod/index.html << 'EOF'
This is prod!
EOF

cat > overlays/prod/kustomization.yaml <<'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: production
commonLabels:
  environment: production
namePrefix: prod-
nameSuffix: -suffix
namespace: kustomize-prod
configMapGenerator:
- name: nginx-index
  behavior: replace
  files:
  - index.html
EOF
# behavior字段仅用于overlays中,有三种可选create|replace|merge
~~~

#### 添加密钥

~~~bash
cat >overlays/test/env.txt<< 'EOF'
env=test
EOF

cat > overlays/test/kustomization.yaml << 'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: test
commonLabels:
  environment: test
namePrefix: test-
nameSuffix: -suffix
namespace: kustomize-test
configMapGenerator:
- name: nginx-index
  behavior: replace
  files:
  - index.html
secretGenerator:
- name: env-file-secret
  envs:
  - env.txt
  type: Opaque
EOF
~~~

### 补丁
~~~bash
cat > overlays/test/kustomization.yaml << 'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../base
commonAnnotations:
  note: test
commonLabels:
  environment: test
namePrefix: test-
nameSuffix: -suffix
namespace: kustomize-test
configMapGenerator:
- name: nginx-index
  behavior: replace
  files:
  - index.html
secretGenerator:
- name: env-file-secret
  envs:
  - env.txt
  type: Opaque
patchesStrategicMerge:
- deployment.yaml
EOF

cat > overlays/test/deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - name: nginx
        resources:
          limits:
            memory: 200Mi
          requests:
            memory: 100Mi
            cpu: 50m
EOF
~~~

### CLI

#### 构建
~~~bash
# 命令
kustomize build <dir>

# 突破目录结构限制
# 当我们引用的资源在应用的目录层级外 如使用同一个harbor仓库登陆secret
# ../../../harbor.yaml

# v3版本
kustomize build --load_restrictor none <dir>

# v4版本
kustomize build --load-restrictor LoadRestrictionsNone <dir>

~~~

