# 数学知识

## 一、指数


$X^AX^B = X^{A+B}$
```math
(X^A) / (X^B) = X^(A-B)
```
```math
(X^A)^B = X^AB
```
```math
(X^N) + (X^N) = 2*(X^N) != X^(2N)
```
```math
2^N + 2^N = 2^(N+1)
```


## 二、对数

仅当 $logxB=A$ 时 $x^A=B$

在 logxB=A 表达式中，x: 底数，B: 真数，A: 以x为底B的对数

在计算机科学中，当底数省略时，默认为2

常数e: 单位时间内增长率为100%的情况下，在单位时间内复合增长率的最大值为e，e是一个无理数，
  约为2.7182804 [详解文章](https://betterexplained.com/articles/an-intuitive-guide-to-exponential-functions-e/)

推导1
```math
设:
C>0
X=logCB
Y=logCA
Z=logAB

可以得出
C^X=B
C^Y=A
A^Z=B
因为
C^Y=A
A^Z=B
所以
(C^Y)^Z=B
因为
C^X=B
所以
(C^Y)^Z=C^X
所以
X=YZ
所以
Z=X/Y

根据设得出以下公式:
logAB = logCB / logCA
```

推导2
```math
设:
X=log A
Y=log B
Z=log AB

可以得出
2^X=A
2^Y=B
2^Z=AB

因为
2^X=A
2^Y=B

所以
2^X * 2^Y = AB
2^(X+Y) = AB

因为
2^Z=AB

所以
2^Z = 2^(X+Y)

所以
Z = X+Y

根据设得出以下公式:
log AB = log A + log B

同理
log A/B = log A - log B
```

推导3
```math
设：
A = 2 ^ N
则
log A = N
则
A^B = (2^N)^B = 2^(NB)
所以
log A^B = NB
因为
log A = N

根据设得出以下公式:
B log A = log A^B
```

简单公式:
```math
log X < X (X>0) 对数小于真数

log 1=0  2^0=1
log 2=1  2^1=2
log 1024=10  2^10=1024
log 1048576=20 2^20=1048576
```

## 三、级数
级数是数学中一个有穷或无穷的序列之和


