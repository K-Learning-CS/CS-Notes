# Lab

## 实验环境

centos

因为实验是32位的所以得补全环境依赖

yum install -y libgcc-4.8.5-44.el7.i686 glibc-devel-2.17-326.el7_9.3.i686



lab下载地址:
http://csapp.cs.cmu.edu/3e/labs.html

Self-Study Handout


### Data Lab

检验对信息表示和处理章节的掌握情况

```shell
1.查看README文档

该实验需要在给定的限制条件(操作数、操作符)内完成题目要求

仅编辑bits.c文件

使用./dlc -e bits.c 检查是否符合限制条件

使用./btest检查逻辑是否符合要求并打分

构建btest: make btest 

修改bits.c文件后需要重新构建  make clean ; make btest
```


