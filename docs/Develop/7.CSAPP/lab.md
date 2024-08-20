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
/* 
 * conditional - same as x ? y : z 
 *   Example: conditional(2,4,5) = 4
 *   Legal ops: ! ~ & ^ | + << >>
 *   Max ops: 16
 *   Rating: 3
 */
int conditional(int x, int y, int z) {
    // 根据x来控制 y z 的表达 当 x 为真时 ~(~!x+1) = 0xffffffff  ~!x+1 = 0x00000000
    // x 为假时则相反
    y=~(~!x+1)&y;
    z=(~!x+1)&z;
  return y+z;
}
/* 
 * isLessOrEqual - if x <= y  then return 1, else return 0 
 *   Example: isLessOrEqual(4,5) = 1.
 *   Legal ops: ! ~ & ^ | + << >>
 *   Max ops: 24
 *   Rating: 3
 */
int isLessOrEqual(int x, int y) {
//    错误思路
//    int i = x ^ y;
//    int j;
//    int z=0x1<<31;
//    i = i & (~z);
//    i |= (i >> 1);
//    i |= (i >> 2);
//    i |= (i >> 4);
//    i |= (i >> 8);
//    i |= (i >> 16);
//    j = i & (~i >> 1);
//    y = (y & z) & (~(x & z));
//    x = x & j;
//
//
//  return !(x+y);
/* (y >=0 && x <0) || ((x * y >= 0) && (y + (-x) >= 0)) */
    int signX = (x >> 31) & 1;
    int signY = (y >> 31) & 1;
    int signXSubY = ((y + ~x + 1) >> 31) & 1;
    return (signX & ~signY) | (!(signX ^ signY) & !signXSubY);
}
//4
/* 
 * logicalNeg - implement the ! operator, using all of 
 *              the legal operators except !
 *   Examples: logicalNeg(3) = 0, logicalNeg(0) = 1
 *   Legal ops: ~ & ^ | + << >>
 *   Max ops: 12
 *   Rating: 4 
 */
int logicalNeg(int x) {
// 解1
    // 用x的最高位往下补充，如果x有值，那么最低位一定为1
    //    x |= (x >> 1);
    //    x |= (x >> 2);
    //    x |= (x >> 4);
    //    x |= (x >> 8);
    //    x |= (x >> 16);
    // 如果x最低位为1则返回0，为零则返回1
    //  return (x&0x1)^0x1;
// 解2
    // if x < 0 ;sign = 1
    int sign = (x >> 31) & 1;
    // TMAX = 0x7fffffff
    int TMAX = ~(1 << 31);
    // if x < 0 return 0;elif x!=0; (x + TMAX) >> 31) = 1
    return (sign ^ 1) & ((((x + TMAX) >> 31) & 1) ^ 1);
}
/* howManyBits - return the minimum number of bits required to represent x in
 *             two's complement
 *  Examples: howManyBits(12) = 5
 *            howManyBits(298) = 10
 *            howManyBits(-5) = 4
 *            howManyBits(0)  = 1
 *            howManyBits(-1) = 1
 *            howManyBits(0x80000000) = 32
 *  Legal ops: ! ~ & ^ | + << >>
 *  Max ops: 90
 *  Rating: 4
 */
int howManyBits(int x) {
//    // 确定正负
//    int sign = x >> 31;
//    //根据正负对x进行处理 将x转换为正数
//    int unsignX =  ~sign & (~x + 1) | x;
//    // 从0x70000000开始位移并计数  计算log2x
//    int n = 0;
//
//  return sign + n;
    int sign, b16, b8, b4, b2, b1, b0;
    sign = x >> 31;
    x = (sign & ~x) | (~sign & x); // 如果x是负数，取反，正数不变
    // x >> 16 将 x 右移 16 位，丢弃掉低 16 位，保留高 16 位。如果高 16 位中有 1，结果将非零
    // !! 将结果转换为0 或者 1
    // << 4 将b16根据是否有1赋值为10000或者00000
    b16 = !!(x >> 16) << 4;
    // 根据b16的值选择是否丢弃掉低16位
    x = x >> b16;
    b8 = !!(x >> 8) << 3;
    x = x >> b8;
    b4 = !!(x >> 4) << 2;
    x = x >> b4;
    b2 = !!(x >> 2) << 1;
    x = x >> b2;
    b1 = !!(x >> 1);
    x = x >> b1;
    b0 = x;

    return b16 + b8 + b4 + b2 + b1 + b0 + 1;
}
```

