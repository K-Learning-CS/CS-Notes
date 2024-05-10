# Mydumper


### 安装

```bash
# centos 7
# 如果访问不了 github 会导致失败
release=$(curl -Ls -o /dev/null -w %{url_effective} https://github.com/mydumper/mydumper/releases/latest | cut -d'/' -f8)
yum install https://github.com/mydumper/mydumper/releases/download/${release}/mydumper-${release:1}.el7.x86_64.rpm

yum install -y epel-release
yum install -y libzstd-devel
```

### 命令详解
```bash
连接选项
-h, --host 需要连接的主机
-u, --user 具有必要权限的用户名
-p, --password 用户密码
-a, --ask-password 提示用户密码
-P, --port 连接的TCP/IP端口
-S, --socket 用于连接的UNIX域套接字文件
--protocol 用于连接的协议 (tcp, socket)
-C, --compress-protocol 在MySQL连接上使用压缩
--ssl 使用SSL连接
--ssl-mode 与服务器连接的期望安全状态：DISABLED, PREFERRED, REQUIRED, VERIFY_CA, VERIFY_IDENTITY
--key 密钥文件的路径名称
--cert 证书文件的路径名称
--ca 证书颁发机构文件的路径名称
--capath 包含PEM格式的受信任SSL CA证书的目录的路径名称
--cipher 用于SSL加密的允许的密码列表
--tls-version 服务器允许的加密连接的协议

过滤选项
-x, --regex 'db.table'匹配的正则表达式
-B, --database 需要导出的数据库
-i, --ignore-engines 需要忽略的存储引擎的逗号分隔列表
--where 仅导出选定的记录
-U, --updated-since 使用Update_time仅导出在过去U天更新的表
--partition-regex 按分区名称过滤的正则表达式
-O, --omit-from-file 包含每行一个需要跳过的database.table条目的文件（在应用regex选项之前跳过）
-T, --tables-list 需要导出的逗号分隔的表列表（不排除regex选项）。表名必须包含数据库名。例如：test.t1,test.t2

锁定选项
-z, --tidb-snapshot 用于TiDB的快照
-k, --no-locks 不执行临时共享读锁。警告：这将导致不一致的备份
--use-savepoints 使用保存点以减少元数据锁定问题，需要SUPER权限
--no-backup-locks 不使用Percona备份锁
--lock-all-tables 对所有的表使用LOCK TABLE，而不是FTWRL
--less-locking 在InnoDB表上最小化锁定时间
--trx-consistency-only 仅事务一致性

PMM选项
--pmm-path 默认值将是/usr/local/percona/pmm2/collectors/textfile-collector/high-resolution
--pmm-resolution 默认将是高分辨率

执行选项
--exec-threads 使用--exec的线程数量
--exec 使用文件作为参数执行的命令
--exec-per-thread 设置将通过STDIN接收并在STDOUT中写入输出文件的命令
--exec-per-thread-extension 当使用--exec-per-thread时，为STDOUT文件设置扩展名

如果发现长时间运行的查询 
--long-query-retries 重试检查长时间运行的查询，默认0（不重试）
--long-query-retry-interval 在重试长时间查询检查之前等待的时间（以秒为单位），默认60
-l, --long-query-guard 设置长时间查询计时器（以秒为单位），默认60
-K, --kill-long-queries 终止长时间运行的查询（而不是中止）

作业选项
--max-rows 在表被估算后，限制每个块的行数，默认1000000。已弃用，使用--rows代替。将在未来的版本中移除
--char-deep 定义当主键为字符串时，使用的字符数量
--char-chunk 该选项定义将表拆分为多少个部分。默认情况下，工具使用可用的线程数。
-r 或 --rows 该选项将表拆分为指定行数的块。指定行数的格式为 MIN:START_AT:MAX。MAX 的值可以设置为 0，表示没有限制。如果查询时间小于 1 秒，则块大小将加倍；如果查询时间大于 2 秒，则块大小将减半。
--split-partitions 将分区转储到单独的文件中。此选项会覆盖分区表的 --rows 选项。

校验选项 
-M 或 --checksum-all 转储所有元素的校验和。
--data-checksums 将数据的表校验和与数据一起转储。
--schema-checksums 转储模式表和视图创建的校验和。
--routine-checksums 转储触发器、函数和例程的校验和。

对象选项 
-m 或 --no-schemas 不要将表的模式与数据和触发器一起转储。
-Y 或 --all-tablespaces 转储所有表空间。
-d 或 --no-data 不要转储表数据。
-G 或 --triggers 转储触发器。默认情况下，不会转储触发器。
-E 或 --events 转储事件。默认情况下，不会转储事件。
-R 或 --routines 转储存储过程和函数。默认情况下，不会转储存储过程和函数。
--views-as-tables 将视图导出为表格形式。
-W 或 --no-views 不要转储视图。

语句选项 
--load-data 使用LOAD DATA语句和.dat文件来替代创建INSERT INTO语句。
--csv 自动启用--load-data，并设置变量以以CSV格式导出。
--fields-terminated-by 定义字段之间的分隔字符。
--fields-enclosed-by 定义用于封闭字段的字符。默认值为`"`。
--fields-escaped-by 用于在LOAD DATA语句中转义字符的单个字符。默认为''。
--lines-starting-by 在每行的开头添加字符串。当使用--load-data时，它会添加到LOAD DATA语句中。当使用INSERT INTO语句时也会生效。
--lines-terminated-by 在每行的末尾添加字符串。当使用--load-data时，它会添加到LOAD DATA语句中。当使用INSERT INTO语句时也会生效。
--statement-terminated-by 除非您知道自己在做什么，否则可能永远不会使用此选项。
-N 或 --insert-ignore 使用INSERT IGNORE转储行。
--replace 使用REPLACE转储行。
--complete-insert 使用包含列名的完整INSERT语句。
--hex-blob 使用十六进制表示法转储二进制列。
--skip-definer 从CREATE语句中删除DEFINER。默认情况下，不修改语句。
-s 或 --statement-size 尝试以字节为单位设置INSERT语句的大小，默认为1000000字节。
--tz-utc 在转储顶部添加SET TIME_ZONE='+00:00'，以允许在具有不同时区数据的服务器之间进行转储，或者在不同时区的服务器之间移动数据，默认为启用。使用--skip-tz-utc禁用此选项。
--skip-tz-utc 不在备份文件中添加SET TIMEZONE。
--set-names 设置名称，使用时需谨慎，默认为二进制。

额外选项 
-F 或 --chunk-filesize 将表拆分为指定输出文件大小的块。该值以MB为单位。
--exit-if-broken-table-found 如果发现损坏的表，则退出。
--success-on-1146 如果表不存在，则不增加错误计数，并将其视为警告而不是关键错误。
-e 或 --build-empty-files 即使表中没有可用数据，也构建转储文件。
--no-check-generated-fields 不执行与生成字段相关的查询。如果存在生成列，则可能导致恢复问题。
--order-by-primary 按主键或唯一键对数据进行排序（如果没有主键）。
-c 或 --compress 使用压缩输出文件，可选值为 /usr/bin/gzip 和 /usr/bin/zstd。选项 GZIP 和 ZSTD。默认为GZIP。

守护进程选项 
-D 或 --daemon 启用守护进程模式。
-I 或 --snapshot-interval 每个转储快照之间的间隔（以分钟为单位），需要使用--daemon，默认为60分钟。
-X 或 --snapshot-count 快照数量，默认为2个。

应用程序选项 
-? 或 --help 显示帮助选项。
-o 或 --outputdir 指定输出文件的目录。
--stream 在文件写入完成后通过标准输出流进行流式传输。自v0.12.7-1起，接受NO_DELETE、NO_STREAM_AND_NO_DELETE和TRADITIONAL等值。如果没有给出参数，则使用TRADITIONAL作为默认值。
-L 或 --logfile 指定日志文件的名称，默认情况下使用标准输出。
--disk-limits 设置磁盘空间限制，如果确定没有足够的磁盘空间，会暂停和恢复操作。接受格式为'<resume>:<pause>'的值，单位为MB。例如：100:500表示当只有100MB可用时暂停，当有500MB可用时恢复。
-t 或 --threads 指定要使用的线程数，默认为4。
-V 或 --version 显示程序版本并退出。
--identifier-quote-character 设置用于INSERT语句的标识符引用字符，仅适用于mydumper并用于myloader的语句拆分。使用SQL_MODE来更改CREATE TABLE语句。可能的取值为：BACKTICK和DOUBLE_QUOTE。默认值为BACKTICK。
-v 或 --verbose 输出详细程度，0 = 静默，1 = 错误，2 = 警告，3 = 信息，默认为2。
--defaults-file 使用特定的默认文件。默认值为/etc/mydumper.cnf。
--defaults-extra-file 使用附加的默认文件。此文件在--defaults-file之后加载，替换先前定义的值。
--fifodir 指定在需要时创建FIFO文件的目录。默认为与备份相同的目录。
```

### 备份

```bash
# --set-names 默认为 binary，会导致 json 导入失败
# --exec-threads 指定启动的线程数
mydumper -u root -p 66666666  -h 127.0.0.1 --regex '^(?!(mysql|information_schema|performance_schema|sys))' --set-names utf8mb4 -G -R -E  --exec-threads 2 -o backup/
```

### 恢复
```bash
# -o 为覆盖导入
myloader --regex '^(?!(mysql|information_schema|performance_schema|sys))'  -o -d backup/
```