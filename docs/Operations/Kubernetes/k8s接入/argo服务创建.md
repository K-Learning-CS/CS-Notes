
## 一、部署 argo-rollouts
~~~bash
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/download/latest/install.yaml
~~~

## 二.部署 argocd


1）部署argocd
~~~bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.3.4/manifests/ha/install.yaml
scp -P12121 47.118.34.109:/usr/local/bin/argocd /usr/local/bin/
~~~

2）修改 web 页面 rootpath
~~~bash
kubectl edit deploy argocd-server -n argocd
...
        - --rootpath
        - /dnet-prd # 视具体情况而定
        - --insecure
...
~~~

3）开放argocd网络
~~~bash
cat > argo-service.yaml <<EOF
apiVersion: v1
kind: Service
metadata:
  namespace: argocd
  labels:
    app.kubernetes.io/component: server
    app.kubernetes.io/name: argocd-server
    app.kubernetes.io/part-of: argocd
  name: argocd-server
spec:
  type: NodePort
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8080
    nodePort: 30002
  - name: https
    port: 443
    protocol: TCP
    targetPort: 8080
  selector:
    app.kubernetes.io/name: argocd-server
EOF
# 在同一内网时 直接在千帆生产集群的proxy中配置prd代理 不在同一内网时需先使用slb代理至外网端口

kubectl apply -f argo-service.yaml
~~~

4）argocd ldap 设置
~~~bash
cat > useradd.yaml <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
  namespace: argocd
  labels:
    app.kubernetes.io/name: argocd-cm
    app.kubernetes.io/part-of: argocd
data:
  kustomize.path.v4.1.3: /usr/local/bin/kustomize
  #kustomize.buildOptions.v4.1.3: --load-restrictor LoadRestrictionsNone
  kustomize.buildOptions.v4.1.3: --load_restrictor none

  #admin.enabled: "false"
  accounts.qianfan: apiKey, login
  accounts.qianfan.enabled: "true"
  url: https://argocd.hd123.com/dnet-prd
  dex.config: |
    connectors:
      - type: ldap
        id: ldap
        name: LDAP
        config:
          host: "ldap.hddomain.cn:389"
          insecureNoSSL: true
          insecureSkipVerify: true
          bindDN: "hddenv@hddomain.cn"
          bindPW: "!hddenv123"
          userSearch:
            baseDN: "OU=HDUsers,dc=hddomain,dc=cn"
            filter: ""
            username: "sAMAccountName"
            idAttr: distinguishedName
            emailAttr: mail
            nameAttr: displayName
          groupSearch:
            baseDN: "OU=HDGroups,dc=hddomain,dc=cn"
            filter: ""
            userAttr: distinguishedName
            groupAttr: member
            nameAttr: name
EOF

kubectl apply -f useradd.yaml
~~~

5）argocd 权限管理
~~~bash
cat > user-rbac.yaml <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-rbac-cm
  namespace: argocd
data:
  #policy.default: role:readonly
  scopes: '[groups,email]'
  policy.csv: |
    #p, role:org-admin, applications, *,/, allow
    p, role:dnet, applications, get, dnet/*, allow
    #p, role:org-admin, clusters, get, *, allow
    #p, role:org-admin, projects, get, *, allow
    #p, role:org-admin, repositories, get, *, allow
    g, buhaiqing@hd123.com, role:admin
    g, kangpeiwen@hd123.com, role:admin
    g, lishuaiqi@hd123.com, role:admin
    g, qianfan, role:admin
    #g, zhangjianlong@hd123.com, role:dnet
    #g, gaochunwei@hd123.com, role:dnet
    #g, liwei@hd123.com, role:dnet
    #g, linzhu@hd123.com, role:dnet
EOF

kubectl apply -f user-rbac.yaml

~~~


6）qianfan用户密码修改
~~~bash
argocd login argocd.hd123.com --grpc-web-root-path /dnet-prd/ --username admin --password $(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)   --insecure
argocd account update-password --account qianfan --new-password a7e01e1512af1685366b --current-password $(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)
~~~

7）禁用admin
~~~bash
cat > useradd.yaml <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
  namespace: argocd
  labels:
    app.kubernetes.io/name: argocd-cm
    app.kubernetes.io/part-of: argocd
data:
  kustomize.path.v4.1.3: /usr/local/bin/kustomize
  #kustomize.buildOptions.v4.1.3: --load-restrictor LoadRestrictionsNone
  kustomize.buildOptions.v4.1.3: --load_restrictor none

  admin.enabled: "false"
  accounts.qianfan: apiKey, login
  accounts.qianfan.enabled: "true"
  url: https://argocd.hd123.com/dnet-prd
  dex.config: |
    connectors:
      - type: ldap
        id: ldap
        name: LDAP
        config:
          host: "ldap.hddomain.cn:389"
          insecureNoSSL: true
          insecureSkipVerify: true
          bindDN: "hddenv@hddomain.cn"
          bindPW: "!hddenv123"
          userSearch:
            baseDN: "OU=HDUsers,dc=hddomain,dc=cn"
            filter: ""
            username: "sAMAccountName"
            idAttr: distinguishedName
            emailAttr: mail
            nameAttr: displayName
          groupSearch:
            baseDN: "OU=HDGroups,dc=hddomain,dc=cn"
            filter: ""
            userAttr: distinguishedName
            groupAttr: member
            nameAttr: name
EOF

kubectl apply -f useradd.yaml
~~~

8）将 argocd 暴露至公网

~~~bash
在先前部署的 ingress 对应的 slb 中将 30002 端口暴露至公网
~~~