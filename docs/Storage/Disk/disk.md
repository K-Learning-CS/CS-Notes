    缓存：将数据放入到缓存区 加快数据读取 读-缓存（cache）
    缓冲：将输入放入到缓冲区 加快数据写入 写-缓冲（buffer）
磁盘操作
# 以*区分重要程度
1.查看分区
fdisk -1	****
gdisk -1	***
1sblk   	*****
 
2.创建分区
fdisk (MBR) *****
gdisk(GPT)  ***
 
3.同步分区
partprobe 	*****
 
4.创建文件系统(格式化)
mkfs
mkfs.ext4/xfs
mkfs -t ext4/xfs
 
mkswap
        要求分区类型为82
swapon  开启
swapoff 关闭
 
5.挂载
临时挂载
    mount   	*****
永久挂载
    /etc/fstab  *****
建议使用blkid查看设备uuid,使用uuid挂载
 
卸载
umount 
磁盘管理
基本管理
1.普通分区
拿到一块磁盘->分区->创建文件系统(格式) ->挂载
 
命令
fdisk
lsblk
partprobe
       如果是虚拟机或云主机，分区完之后会立即刷新
       如果是物理主机，分区完之后必须手动刷新
mkfs
mount
    属性
blkid
umount
 
文件
/etc/ fstab
六个字段
进阶管理
RAID
RAID-0:
   读、写性能提升:
   可用空间: N*min(S1,S2,...)
   无容错能力
   最少磁盘数: 2, 2+
 
    RAID0把多块(至少两块)物理硬盘设备通过软件或硬件的方式串联在一起，组成一个
大的卷组，并将数据依次写入到各个物理硬盘中。
    在最理想的情况下，硬盘设备的读写性能会提升数倍，但是若任意一块硬盘发生故障将导
致整个系统的数据都受到破坏。
RAID-1:
   读性能提升、写性能略有下降:
   可用空间: 1*min(S1,S2,...)
   有冗余能力
   最少磁盘数: 2, 2+
 
    RAID1在多块硬盘设备中写入了相同的数据，因此硬盘设备的利用率下降了一半。空间的真实可用率只有50%，由三块硬盘设备组成的RAID1磁盘阵列的可用率只有33%左右，以此类推。
    由于需要把数据同时写入两块以上的硬盘设备，这无疑也在一定程度上增大了系统计算功
能的负载。
RAID-4:
    1101, 0110, 1011  互相校验
 
    RAID4 与 RAID3 的原理大致相同，区别在于条带化的方式不同。 
    RAID4按照块的方式来组织数据，写操作只涉及当前数据盘和校验盘两个盘，多个 I/O 
请求可以同时得到处理，提高了系统性能。
    RAID4 按块存储可以保证单块的完整性，可以避免受到其他磁盘上同条带产生的不利影响。
RAID-5：
   读、写性能提升
   可用空间: (N-1)*min(S1,S2,...) 
   有容错能力: 1块磁盘
   最少磁盘数; 3, 3+
RAID-6:
   读、写性能提升
   可用空间: (N-2)*min(S1,S2,...)
   有容错能力: 2块磁盘
   最少磁盘数: 4, 4+
 
    RAID6：带有两种分布存储的奇偶校验码的独立磁盘结构
    RAID6是对RAID5的扩展，主要是用于要求数据绝对不能出错的场合。
RAID-10:
   读、写性能提升
   可用空间: N*min(S1,52,...)/2
   有容错能力:每组镜像最多只能坏一块: 
   最少磁盘数: 4, 4+
 
    鉴于RAID5因为磁盘设备的成本问题对读写速度和数据的安全性能而有了一定的妥协，但是
企业里更在乎的还是数据本身的价值而非硬盘的价格，因此在生产环境中推荐使用RAID10技术。
    RAID10即RAID0+RAID1的一个组合体。
    RAID10技术需要至少4块硬盘来组建，其中先分别两两制作成RAID1磁盘阵列，以保证数据的安全性；然后再对两个RAID1次哦按阵列实施RAID0技术，进一步提高硬盘设备的读写速度。 
    这样从理论上讲，只要坏的不是同一组中的所有磁盘，那么最多可以损坏50%的硬盘设备而不丢失数据。由于RAID10技术继承了RAID0的高速写速度和RAID1的数据安全性，在不考虑成本的情况下RAID10的性能都超过了RAID5，因此当前成为广泛使用的一种存储技术。
RAID-01： 
    这种架构的安全性低于raid10，而两者由于IO数量一致。读写速度相同，使用的硬盘数量也一致。
 
    所以raid10比raid01是一种更为先进的架构。
RAID 等级	RAID0	RAID1	RAID3	RAID5	RAID6	RAID10
别名	条带	镜像	专用奇偶校验条带	分布奇偶校验条带	双重奇偶校验条带	镜像加条带
容错性	无	有	有	有	有	有
冗余类型	无	有	有	有	有	有
热备份选择	无	有	有	有	有	有
读性能	高	低	高	高	高	高
随机写性能	高	低	低	一般	低	一般
连续写性能	高	低	低	低	低	一般
需要磁盘数	n≥1	2n (n≥1)	n≥3	n≥3	n≥4	2n(n≥2)≥4
可用容量 ​	全部	50%	(n-1)/n	(n-1)/n	(n-2)/n	50%
常用级别: RAID-0, RAID-1, RAID-5, RAID-10，RAID-50
实现方式;
硬件实现方式
软件实现方式
 
CentOS 6上的软件RAID的实现:
结合内核中的md(multi devices)
mdadm:模式化的工具
 
命令的语法格式: mdadm [ mode ] <raiddevice> [ options ] <component -devices>
       支持的RAID级别: LINEAR, RAIDO, RAID1, RAID4, RAID5, RAID6, RAID10;
 
模式: 
    创建: -C
    装配: -A
    监控: -F
    管理: -f, -r, -a
 
<raiddevice>: /dev/md#
<component-devices>:任意块设备
 
-C: 创建模式
    -n #:使用#个块设备来创建此RAID: 
    -1 #:指明要创建的RAID的级别;
    -a {yes|no}; 自动创建目标RAID设备的设备文件; 
    -C CHUNK SIZE: 指明块大小:
    -x #:指明空闲盘的个数:
 
-D: 显示raid的详细信息; 
    mdadm -D /dev/md#
 
管理模式:
      -f: 标记指定磁盘为损坏;
      -a: 添加磁盘
      -r: 移除磁盘
 
观察md的状态: 
    cat /proc/mdstat
 
停止md设备:
    mdadm -S /dev/md#
 
watch命令:
    -n #:刷新间隔，单位是秒:
 
    watch -n# 'COMMAND'
LVM
    逻辑（动态）分区
功能/命令   物理卷管理     卷组管理    逻辑卷管理
  扫描       pvscan      vgscan      lvscan
  建立       pvcreate    vgcreate    lvcreate
  显示       pvdisplay   vgdisplay   lvdisplay
  删除       pvremove    vgremove    lvremove
  扩展                   vgextend    lvextend
  缩小                   vgreduce    lvreduce
# 软RAID和lvm是通过软件来实现，会占用系统资源
总结
磁盘
作用:存放数据
磁盘规划:
 
1.分区
fdisk -1
lsblk
fdisk分区名(/dev/sdb)
 
2.创建文件系统，格式化
mkfs.ext4/xfs /dev/ sdb1
mkfs -t ext4/xfs /dev/sdb1
mkswap
 
3.挂载
mount /dev/sdb1 /dir
永久挂载
/etc/fstab
 
 为什么要用RAID LVM
问题:
1.硬盘坏了
2.硬盘空间不够
3.硬盘读写性能较差