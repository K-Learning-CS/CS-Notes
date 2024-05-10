1.ssh
```bash
# 跳板机添加管理用户公钥
# 1.在需要登陆跳板机的电脑上生成密钥对
ssh-keygen

# 2.查看公钥
cat ~/.ssh/id_rsa.pub

# 3.将公钥添加至跳板机
vim /root/.ssh/authorized_keys

# 跳板机免密登陆其他云服务器
# 1.在跳板机上生成密钥对
ssh-keygen

# 2.将公钥添加至其他内网机器
ssh-copy-id  -i /root/.ssh/id_rsa.pub root@10.88.70.71
```

2.公共用户
```bash
/usr/sbin/groupadd www
/usr/sbin/useradd -g www www -u 1000
```

3.初始化目录
```bash
# 此目录结构为规定的数据及文件存放目录，服务器上维护的服务必须按照当前规定创建及存放文件或数据

# init.d: 存放所有维护服务的启动脚本
# shell: 存放所有维护的shell脚本
# tools: 存放所有使用到的工具包文件
# pids: 存放所有维护服务的pid文件
# codes: 存放所有的代码项目
# logs: 存放所有维护服务的日志
# data: 存放所有维护服务的数据

mkdir -p /export/{pids,init.d,shell,tools,codes}
mkdir -p /export/logs/{php,nginx}
mkdir -p /export/data/{mysql,redis,elasticsearch,containerd}

chmod 757 /export/pids
```
4.安装目录
```bash
# 所有安装的第三方服务都安装在此目录下，比如redis、mysql
/usr/local/
```
5.命令提示符
```bash
cat >> ~/.bashrc <<'EOF'
IP=$(ifconfig eth0 | awk ' /inet /  { print $2 }')
if test -z "$IP"
then
        IP=$(hostname | awk -F. ' { print $1 } ')
fi
export IP
export PS1="[\u@$IP \w\$]# "
EOF
```

7.监控
```bash

```
8.防火墙
```bash

```
9.系统性能参数优化
```bash
grep -i net.ipv4.icmp_ignore_bogus_error_responses /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.icmp_ignore_bogus_error_responses =\).*#\1 1#g' /etc/sysctl.conf || echo 'net.ipv4.icmp_ignore_bogus_error_responses = 1' >>/etc/sysctl.conf
grep -i net.ipv4.ip_forward /etc/sysctl.conf &>/dev/null &&  sed -i 's#\(net.ipv4.ip_forward =\).*#\1 1#g' /etc/sysctl.conf || echo 'net.ipv4.ip_forward = 1' >> /etc/sysctl.conf
grep -i net.ipv4.conf.default.rp_filter /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.conf.default.rp_filter =\).*#\1 0#g' /etc/sysctl.conf || echo 'net.ipv4.conf.default.rp_filter = 0' >>/etc/sysctl.conf
grep -i net.ipv4.conf.all.rp_filter /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.conf.all.rp_filter =\).*#\1 0#g' /etc/sysctl.conf || echo 'net.ipv4.conf.all.rp_filter = 0' >>/etc/sysctl.conf
grep -i net.ipv4.conf.eth0.rp_filter /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.conf.eth0.rp_filter =\).*#\1 0#g' /etc/sysctl.conf || echo 'net.ipv4.conf.eth0.rp_filter = 0' >>/etc/sysctl.conf
grep -i net.ipv4.conf.default.accept_source_route /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.conf.default.accept_source_route =\).*#\1 0#g' /etc/sysctl.conf || echo 'net.ipv4.conf.default.accept_source_route = 0' >>/etc/sysctl.conf
grep -i kernel.sysrq /etc/sysctl.conf &>/dev/null && sed -i 's#\(kernel.sysrq =\).*#\1 0#g' /etc/sysctl.conf || echo 'kernel.sysrq = 0' >>/etc/sysctl.conf
grep -i kernel.core_uses_pid /etc/sysctl.conf &>/dev/null && sed -i 's#\(kernel.core_uses_pid =\).*#\1 1#g' /etc/sysctl.conf || echo 'kernel.core_uses_pid = 1' >>/etc/sysctl.conf
grep -i kernel.msgmnb /etc/sysctl.conf &>/dev/null && sed -i 's#\(kernel.msgmnb =\).*#\1 65536#g' /etc/sysctl.conf || echo 'kernel.msgmnb = 65536' >>/etc/sysctl.conf
grep -i kernel.msgmax /etc/sysctl.conf &>/dev/null && sed -i 's#\(kernel.msgmax =\).*#\1 68719476736#g' /etc/sysctl.conf || echo 'kernel.msgmax = 68719476736' >>/etc/sysctl.conf
grep -i kernel.shmall /etc/sysctl.conf &>/dev/null && sed -i 's#\(kernel.shmall =\).*#\1 4294967296#g' /etc/sysctl.conf || echo 'kernel.shmall = 4294967296' >>/etc/sysctl.conf
grep -i net.ipv4.tcp_sack /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_sack =\).*#\1 1#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_sack = 1' >>/etc/sysctl.conf
grep -i net.ipv4.tcp_fack /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_fack =\).*#\1 1#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_fack = 1' >>/etc/sysctl.conf
grep -i net.ipv4.tcp_window_scaling /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_window_scaling =\).*#\1 1#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_window_scaling = 1' >>/etc/sysctl.conf
grep -i net.ipv4.tcp_rmem /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_rmem =\).*#\1 8760 256960 4088000#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_rmem = 8760 256960 4088000' >>/etc/sysctl.conf
grep -i net.ipv4.tcp_wmem /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_wmem =\).*#\1 8760 256960 4088000#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_wmem = 8760 256960 4088000' >>/etc/sysctl.conf
grep -i net.ipv4.conf.all.send_redirects /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.conf.all.send_redirects =\).*#\1 0#g' /etc/sysctl.conf || echo 'net.ipv4.conf.all.send_redirects = 0' >>/etc/sysctl.conf
grep -i net.ipv4.conf.all.secure_redirects /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.conf.all.secure_redirects =\).*#\1 0#g' /etc/sysctl.conf || echo 'net.ipv4.conf.all.secure_redirects = 0' >>/etc/sysctl.conf
grep -i net.ipv4.conf.all.accept_redirects /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.conf.all.accept_redirects =\).*#\1 0#g' /etc/sysctl.conf || echo 'net.ipv4.conf.all.accept_redirects = 0' >>/etc/sysctl.conf
grep -i kernel.unknown_nmi_panic /etc/sysctl.conf &>/dev/null && sed -i 's#\(kernel.unknown_nmi_panic =\).*#\1 0#g' /etc/sysctl.conf || echo 'kernel.unknown_nmi_panic = 0' >>/etc/sysctl.conf
grep -i vm.swappiness /etc/sysctl.conf &>/dev/null && sed -i 's#\(vm.swappiness =\).*#\1 10#g' /etc/sysctl.conf || echo 'vm.swappiness = 10' >>/etc/sysctl.conf
grep -i fs.inotify.max_user_watches /etc/sysctl.conf &>/dev/null && sed -i 's#\(fs.inotify.max_user_watches =\).*#\1 10000000#g' /etc/sysctl.conf || echo 'fs.inotify.max_user_watches = 10000000' >>/etc/sysctl.conf
grep -i net.ipv6.conf.all.disable_ipv6 /etc/sysctl.conf &>/dev/null &&  sed -i 's#\(net.ipv6.conf.all.disable_ipv6 =\).*#\1 1#g' /etc/sysctl.conf ||  echo 'net.ipv6.conf.all.disable_ipv6 = 1' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_syncookies /etc/sysctl.conf &>/dev/null &&  sed -i 's#\(net.ipv4.tcp_syncookies =\).*#\1 1#g' /etc/sysctl.conf ||  echo 'net.ipv4.tcp_syncookies = 1' >> /etc/sysctl.conf
grep -i net.core.optmem_max /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.core.optmem_max =\).*#\1 327680#g' /etc/sysctl.conf || echo 'net.core.optmem_max = 327680' >> /etc/sysctl.conf
grep -i net.core.netdev_max_backlog /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.core.netdev_max_backlog =\).*#\1 1048576#g' /etc/sysctl.conf || echo 'net.core.netdev_max_backlog = 1048576' >> /etc/sysctl.conf
grep -i net.core.rmem_default /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.core.rmem_default =\).*#\1 8388608#g' /etc/sysctl.conf || echo 'net.core.rmem_default = 8388608' >> /etc/sysctl.conf
grep -i net.core.wmem_default /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.core.wmem_default =\).*#\1 8388608#g' /etc/sysctl.conf || echo 'net.core.wmem_default = 8388608' >> /etc/sysctl.conf
grep -i net.core.rmem_max /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.core.rmem_max =\).*#\1 16777216#g' /etc/sysctl.conf || echo 'net.core.rmem_max  = 16777216' >> /etc/sysctl.conf
grep -i net.core.wmem_max /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.core.wmem_max =\).*#\1 16777216#g' /etc/sysctl.conf || echo 'net.core.wmem_max = 16777216' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_fin_timeout /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_fin_timeout =\).*#\1 60#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_fin_timeout = 60' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_keepalive_time /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_keepalive_time =\).*#\1 600#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_keepalive_time = 600' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_keepalive_intvl /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_keepalive_intvl =\).*#\1 30#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_keepalive_intvl = 30' >>/etc/sysctl.conf
grep -i net.ipv4.tcp_keepalive_probes /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_keepalive_probes =\).*#\1 3#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_keepalive_probes = 3' >>/etc/sysctl.conf
grep -i net.ipv4.tcp_timestamps /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_timestamps =\).*#\1 1#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_timestamps = 1' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_synack_retries /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_synack_retries =\).*#\1 2#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_synack_retries = 2' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_syn_retries /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_syn_retries =\).*#\1 2#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_syn_retries = 2' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_max_tw_buckets /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_max_tw_buckets =\).*#\1 6000#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_max_tw_buckets = 6000' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_slow_start_after_idle /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_slow_start_after_idle =\).*#\1 0#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_slow_start_after_idle = 0' >>/etc/sysctl.conf
grep -i net.ipv4.route.gc_timeout /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.route.gc_timeout =\).*#\1 100#g' /etc/sysctl.conf || echo 'net.ipv4.route.gc_timeout = 100' >>/etc/sysctl.conf
grep -i net.ipv4.tcp_tw_recycle /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_tw_recycle =\).*#\1 0#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_tw_recycle = 0' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_tw_reuse /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_tw_reuse =\).*#\1 1#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_tw_reuse = 1' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_mem /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_mem =\).*#\1 94500000 915000000 927000000#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_mem = 94500000 915000000 927000000' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_max_orphans /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_max_orphans =\).*#\1 3276800#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_max_orphans = 3276800' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_max_syn_backlog /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_max_syn_backlog =\).*#\1 1048576#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_max_syn_backlog = 1048576' >> /etc/sysctl.conf
grep -i net.ipv4.tcp_moderate_rcvbuf /etc/sysctl.conf &>/dev/null && sed -i 's#\(net.ipv4.tcp_moderate_rcvbuf =\).*#\1 1#g' /etc/sysctl.conf || echo 'net.ipv4.tcp_moderate_rcvbuf = 1' >> /etc/sysctl.conf
grep -i fs.file-max /etc/sysctl.conf &>/dev/null && sed -i 's#\(fs.file-max =\).*#\1 102400#g' /etc/sysctl.conf || echo 'fs.file-max = 102400' >> /etc/sysctl.conf
grep -i fs.nr_open /etc/sysctl.conf &>/dev/null && sed -i 's#\(fs.nr_open =\).*#\1 102400#g' /etc/sysctl.conf || echo 'fs.nr_open  = 102400' >> /etc/sysctl.conf

/usr/sbin/sysctl -p



```

10.常用包安装
```bash
yum install wget tree expect vim net-tools ntp bash-completion ipvsadm ipset jq iptables conntrack sysstat libseccomp -y
```
11.垃圾回收站
```bash

```
12.limit
```bash
grep  '\*.*soft.*nproc.*' /etc/security/limits.conf &> /dev/null && sed -i 's#\(\*.*soft.*nproc\).*#\1   102400#' /etc/security/limits.conf|| echo '*          soft    nproc    102400'>>/etc/security/limits.conf
grep  '\*.*hard.*nproc.*' /etc/security/limits.conf &> /dev/null && sed -i 's#\(\*.*hard.*nproc\).*#\1   102400#' /etc/security/limits.conf|| echo '*          hard    nproc    102400'>>/etc/security/limits.conf
grep  '\*.*soft.*nofile.*' /etc/security/limits.conf &> /dev/null && sed -i 's#\(\*.*soft.*nofile\).*#\1   102400#' /etc/security/limits.conf|| echo '*          soft    nofile    102400'>>/etc/security/limits.conf
grep  '\*.*hard.*nofile.*' /etc/security/limits.conf &> /dev/null && sed -i 's#\(\*.*hard.*nofile\).*#\1   102400#' /etc/security/limits.conf|| echo '*          hard    nofile    102400'>>/etc/security/limits.conf
echo 'ulimit -SHn  102400' >> /etc/profile
grep  '\*.*soft.*nproc.*' /etc/security/limits.d/20-nproc.conf && sed -i 's#\(\*.*soft.*nproc\).*#\1   102400#' /etc/security/limits.d/20-nproc.conf|| echo '*          soft    nproc    102400'>>/etc/security/limits.d/20-nproc.conf
echo 'ulimit -c  unlimited' >>/etc/profile
```

13.selinux
```bash
setenforce 0
sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config
```

14.crontab
```bash
sed -i 's/MAILTO=root/MAILTO=""/g' /etc/crontab
service crond restart &>/dev/null
```

15.history
```bash
echo 'export HISTTIMEFORMAT="%Y-%m-%d_%H:%M:%S "' >>/etc/profile
```