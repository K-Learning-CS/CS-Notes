### 时间序列

- 一个时间序列是一个以时间为索引的数据流，它是一个由时间戳和浮点数值组成的序列

```promql
# 以node_cpu_seconds_total为例，它由以下几个部分组成

|-----metrics name----｜---------------labels-----------------|-timestamp-| |--value--|
node_cpu_seconds_total{cpu="0",instance="master-0",mode="idle"}@1689325440   3912085.06
|--------------------------------key--------------------------------------| |--value--|

# 存储的时候则是以key-value的形式存储
```

对于prometheus时间序列而言，它有四种数据类型：

1. 瞬时向量
2. 范围向量
3. 标量
4. 字符串

### 瞬时向量

- 在prometheus中，没从exporter拉取一次数据，便会更新一次数据，这个数据，就是瞬时向量
- 当然，旧的数据也是瞬时向量，只不过它的时间戳比较旧而已，只要是单独一个数据，就是瞬时向量
- 绘制图表展示数据时，只能对瞬时向量进行展示，因为瞬时向量一个点是一个数据，而范围向量一个点是一组数据，无法展示

```promql
node_cpu_seconds_total
```

### 范围向量

- 范围向量是一组时间序列，它们的时间戳都在一个范围内，也可以理解为一组瞬时向量
- 无法直接展示，但是可以用内置函数处理成瞬时向量，然后展示

```promql
node_cpu_seconds_total[5m]
```

### 标量

- 标量是一个单独的浮点数值，它没有时间戳，也没有标签，也就是说，它不是一个时间序列
- 只能通过内置函数scalar()将范围向量转换为标量

### 字符串

- 一个简单的字符串。字符串用单双引号或反引号来指定。



```promql


```promql
# rate( ... )[30m:1m]：这是一个子查询，用于计算数据的速率
# 以一分钟每次的频率查询近三十分钟的数据，也就是说会得到三十个数据点
# 在这个三十个数据点中，每次查询的内容为计算过去五分钟指标的变化速率
rate(node_cpu_seconds_total{cpu="0",instance="master-0",mode="idle"}[5m])[30m:1m]
```



### predict_linear

- 用于预测未来的值

```promql
# predict_linear(v range-vector, t scalar) 

predict_linear(node_filesystem_free_bytes[24h],240*3600)
# 以过去24小时的数据为基础，预测10天后的数据量
```

### rate

- 用于计算速率

```promql
# rate(v range-vector)

rate(node_cpu_seconds_total{mode="idle"}[5m])
# 计算过去5分钟的CPU空闲率
```

### sort

- 用于排序

```promql
# sort(v instant-vector)

sort(node_cpu_seconds_total{mode="idle"})
# 对CPU空闲率进行排序
```

### topk

- 用于取出最大的K个值

```promql
# topk(k int, v instant-vector)

topk(5, node_cpu_seconds_total{mode="idle"})
# 取出CPU空闲率最高的5个节点
```

### bottomk

- 用于取出最小的K个值

```promql
# bottomk(k int, v instant-vector)

bottomk(5, node_cpu_seconds_total{mode="idle"})
# 取出CPU空闲率最低的5个节点
```

### quantile

- 用于计算分位数

```promql
# quantile(φ float, v instant-vector)

quantile(0.9, node_cpu_seconds_total{mode="idle"})
# 计算CPU空闲率的90%分位数
```

### histogram_quantile

- 用于计算直方图的分位数

```promql
# histogram_quantile(φ float, b instant-vector)

histogram_quantile(0.9, node_cpu_seconds_total{mode="idle"})
# 计算CPU空闲率的90%分位数
```

### absent

- 用于判断指标是否存在

```promql
# absent(v instant-vector)

absent(node_cpu_seconds_total{mode="idle"})
# 判断CPU空闲率是否存在
```

### changes

- 用于计算指标值发生变化的次数

```promql
# changes(v range-vector)

changes(node_cpu_seconds_total{mode="idle"}[5m])
# 计算过去5分钟CPU空闲率发生变化的次数
```

### delta

- 用于计算指标值的增量

```promql
# delta(v range-vector)

delta(node_cpu_seconds_total{mode="idle"}[5m])
# 计算过去5分钟CPU空闲率的增量
```

### deriv

- 用于计算指标值的导数

```promql
# deriv(v range-vector)

deriv(node_cpu_seconds_total{mode="idle"}[5m])
# 计算过去5分钟CPU空闲率的导数
```

### idelta

- 用于计算指标值的增量

```promql
# idelta(v range-vector)

idelta(node_cpu_seconds_total{mode="idle"}[5m])
# 计算过去5分钟CPU空闲率的增量
```

