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

```c
//1
/* 
 * bitXor - x^y using only ~ and & 
 *   Example: bitXor(4, 5) = 1
 *   Legal ops: ~ &
 *   Max ops: 14
 *   Rating: 1
 */
int bitXor(int x, int y) {
  return ~(~x&~y)&~(x&y); // 需要学习布尔代数
}
/* 
 * tmin - return minimum two's complement integer 
 *   Legal ops: ! ~ & ^ | + << >>
 *   Max ops: 4
 *   Rating: 1
 */
int tmin(void) {

  return 0x1<<31; // 简单的移位操作

}
//2
/*
 * isTmax - returns 1 if x is the maximum, two's complement number,
 *     and 0 otherwise 
 *   Legal ops: ! ~ & ^ | +
 *   Max ops: 10
 *   Rating: 1
 */
// maximum 0111... ;int is 32 bit; 0x7fffffff
int isTmax(int x) {
    int i = x+1; // 0111... -> 1000...
    x=x+i; // 1111...
    x=~x; // 0000...
    i=!i; //exclude x=0xffff...
    x=x+i; //exclude x=0xffff...
  return !x;
//  return !(~(1<<31)^x);
}
/* 
 * allOddBits - return 1 if all odd-numbered bits in word set to 1
 *   where bits are numbered from 0 (least significant) to 31 (most significant)
 *   Examples allOddBits(0xFFFFFFFD) = 0, allOddBits(0xAAAAAAAA) = 1
 *   Legal ops: ! ~ & ^ | + << >>
 *   Max ops: 12
 *   Rating: 2
 */
int allOddBits(int x) {
    int i = (0xaa << 24) + (0xaa << 16)+ (0xaa << 8) + 0xaa; // 拼出0xAAAAAAAA
  return !(((x^i)|i)^x);
}
/* 
 * negate - return -x 
 *   Example: negate(1) = -1.
 *   Legal ops: ! ~ & ^ | + << >>
 *   Max ops: 5
 *   Rating: 2
 */
int negate(int x) {
  // 以补码0xFF为例 11111111 的值为-1 因为最高位永远会比其他位之和大1
  // 假设1要转为负数 (1 ^ 0xFF) + 1 即可
  return (~(0x0)^x) + 1; 
}
//3
/* 
 * isAsciiDigit - return 1 if 0x30 <= x <= 0x39 (ASCII codes for characters '0' to '9')
 *   Example: isAsciiDigit(0x35) = 1.
 *            isAsciiDigit(0x3a) = 0.
 *            isAsciiDigit(0x05) = 0.
 *   Legal ops: ! ~ & ^ | + << >>
 *   Max ops: 15
 *   Rating: 3
 */
int isAsciiDigit(int x) {
    int y = (x & ~0xf) ^ 0x30; // 排除不是 0x30 段的
    int z = !((x & 0xe) ^ 0xa); // 排除101开头的 ab
    x = !((x & 0xc) ^ 0xc); // 排除11开头的 cdef
  return !(x+y+z);
}
```

