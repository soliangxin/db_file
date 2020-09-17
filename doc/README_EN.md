### [中文](https://github.com/soliangxin/db_file/blob/master/README.md)  |  [English](https://github.com/soliangxin/db_file/blob/master/doc/README_EN.md)
## db_file Export data from database to file
---

### Introduction 
##### Exporting data from a database, usually we can use the command line provided by the database to export data, but when cross-database, we often need to export as a file, and then import to other databases, for example: MySQL, PostgreSQL, SQLite3, etc. are common , Db_file is to export data from a variety of databases.

---
### function list
- Support database：Clickhouse、Cassandra、Hive、Apache ignite、SQL Server、MySQL、PostgreSQL、Presto、SQLite3
- Write file format：Uncompressed file、Gzip、Lz4、Snappy、Zstd
- Support character set format conversion for writing files
- Supports writing and overwriting existing files
- Use PID files to prevent processes from starting repeatedly
- Support input SQL from command line or read SQL from file
- Support custom separators between fields
- Support field adding tag
- Support custom line breaks

---

### Command Line
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

---
### Parameter introduction
- cache-num When data is written to the file, when this value is reached, the data is forced to be flushed to the disk
- column-name Add the field names of the columns to the output file
- compress-format When writing data to a file, you can specify the compression format of the file, and the default is not to use compression for writing
- empty-val When the query result is NULL, replace the specified result with this value
- encoding-error When using the character encoding conversion function, the processing rules when encountering conversion errors, strict: the program exits, ignore: ignore the error, and write the original result, skip: skip the record, do not write the file
- from-encoding Specify the original character set encoding when converting the character set
- level Set the log level, the default is "info"
- newline Specify the newline character of the file, the default is "\n"
- overwrite If the specified file exists, choose whether to overwrite the file. Note: The compressed file can only be written in overwrite mode
- pid Use PID files to prevent processes from starting repeatedly
- sep The separator between output file fields, the default is ";"
- sql Get the SQL statement to be executed from the command line
- sql-file Get the SQL statement from the file, exclude the lines beginning with "#", ";", and "-", and use the first place separated by ";" for each line. When the parameters --sql and --sql-file are specified at the same time, the parameters specified by --sql are preferred
- tag Tags added on both sides of the field, not added by default
- tag-all Whether to specify all fields to add Tag
- tag-exclude When tag is added to the field, the field type specified by the parameter does not add Tag, and the type is database type. When the parameter --tag-all is True, all fields are added and this parameter is not used
- to-encoding Specify the output character set encoding when converting the character set encoding
- url To connect to the database URL, please refer to the command line URL EXAMPLE
- write Specify the file name to write the output file
---
