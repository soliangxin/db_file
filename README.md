## db_file 从数据库导出数据到文件
---

### 简介 
##### 从数据库导出数据， 通常我们可以使用数据库提供的命令行进行导出数据，但是跨数据库时，我们常常需要先导出为文件，然后再导入到其他数据库，例如：常见的就有MySQL、PostgreSQL、SQLite3等，db_file就是为了从多种数据库导出数据。


---
### 功能列表
- 支持数据库：Clickhouse、Cassandra、Hive、Apache ignite、SQL Server、MySQL、PostgreSQL、Presto、SQLite3
- 写入文件格式：未压缩文件、Gzip、Lz4、Snappy、Zstd
- 写入文件支持字符集格式转换
- 支持写入覆盖现有文件
- 使用PID文件防止进程重复启动
- 支持输入SQL从命令行或者从文件中读取SQL

---

### 命令行
```
$ db_file 
NAME:
   db_file

USAGE:
   db_file [command options] [arguments...]
   
COMMANDS:
   --cache-num value           the maximum number of records to write to the cache (default: 1000)
   --column-name               output file writes column names (default: false)
   --compress-format value     output file compression, valid values (gzip, lz4, snappy, zstd)
   --empty-val value           the value filled in when the field value is NULL
   --encoding-error value      conversion coding error, valid values (strict, ignore, skip) (default: strict)
   --from-encoding value       input character set encoding
   --level value, -l value     current console log level, valid values (trace, debug, info, warn, error, fatal, panic) (default: info)
   --newline value, -n value   output file newline character (default: \n)
   --overwrite, -o             if the output file exists, whether to overwrite the file (default: false)
   --pid value, -p value       pid file to prevent multiple process starts
   --sep value, -s value       output file separator (default: ;)
   --sql value                 execute the exported SQL statement
   --sql-file value, -f value  get the sql from the file
   --tag value, -t value       add a tag to the output field
   --tag-all                   all fields are added with tags, and the default non-numeric type is added (default: false)
   --tag-exclude value         The database type excluded when adding the tag, with multiple types separated by "," (default: INT,BIGINT)
   --to-encoding value         output character set encoding
   --url value, -u value       database connection url
   --write value, -w value     output file name
   --help, -h                  show help (default: false)
   --version, -v               print the version (default: false)
   
AVAILABLE DRIVERS:
   clickhouse [ch]
   cql [ca, scy, scylla, datastax, cassandra]
   hive [hi, hive]
   ignite [ig, gridgain]
   mssql [ms, sqlserver]
   mysql [my, maria, aurora, mariadb, percona]
   postgres [pg, pgsql, postgresql]
   presto [pr, prs, prestos, prestodb, prestodbs]
   sqlite3 [sq, file, sqlite]

URL EXAMPLE:
   ch://user:pass@localhost:port/dbname
   ca://user:pass@localhost:port/keyspace
   hi://user:pass@localhost:port/dbname
   ig://user:pass@localhost:port/dbname
   ms://user:pass@localhost.com/instance/dbname
   my://user:pass@localhost:port/dbname
   pg://user:pass@localhost:port/dbname
   pr://user:pass@localhost:port/dbname
   sq:/path/to/file.db

VERSION:
   1.0.0
   
```
