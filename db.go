package db_file

import (
	"bytes"
	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go"         // ClickHouse (clickhouse)
	_ "github.com/MichaelS11/go-cql-driver"         // Cassandra (cql)
	_ "github.com/amsokol/ignite-go-client/sql"     // Apache Ignite (ignite)
	_ "github.com/denisenkom/go-mssqldb"            // Microsoft SQL Server (mssql)
	_ "github.com/go-sql-driver/mysql"              // MySQL (mysql)
	_ "github.com/lib/pq"                           // PostgreSQL (postgres)
	_ "github.com/mattn/go-sqlite3"                 // SQLite3 (sqlite3)
	_ "github.com/prestodb/presto-go-client/presto" // Presto (presto)
	log "github.com/sirupsen/logrus"
	"github.com/soliangxin/db_file/util"
	"github.com/xo/dburl"
	"net/url"
	"reflect"
	_ "sqlflow.org/gohive" // hive
	"strings"
)

// 数据库操作
type DB struct {
	args                 *Arguments     // 命令行参数
	conn                 *sql.DB        // 数据库操作对象
	rows                 *sql.Rows      // 保存查询结果
	columns              []string       // 保存查询出来的列名称
	databaseTypeName     []string       // 保存数据库字段类型名称
	scanType             []reflect.Type // 保存数据扫描时Go的类型
	flashNum             int64          // 刷新到磁盘Num
	echoNum              int64          // 输出Num
	fileNum              int64          // 写入文件记录数
	recordBuffer         bytes.Buffer   // 写入数据时使用, 缓存
	fieldVal             string         // 字段值
	encoding             *util.Encoding // 用于字符集转换
	isEncodingConvert    bool           // 是否进行字符集转换
	noCoverCharNum       int64          // 未转换字符编码记录数
	suCoverCharNum       int64          // 转换成功字符编码记录数
	igCoverCharNum       int64          // 转换错误忽略字符编码记录数
	skCoverCharNum       int64          // 转换错误跳过字符编码记录数
	sTagExcludeFieldType []string       // 排除的字段类型数组
}

// 执行数据库操作初始化
func (d *DB) Init(args *Arguments) error {
	d.args = args // 保存参数
	// 遍历接收到添加标签排除的字段类型, 拆分为数组
	tagExcludeFieldType := strings.Split(args.TagExcludeFieldType, ",")
	for _, t := range tagExcludeFieldType {
		d.sTagExcludeFieldType = append(d.sTagExcludeFieldType, strings.ToUpper(t))
	}
	// 字符集转换不为空, 初始化字符集转换
	if args.FromEncoding != "" && args.ToEncoding != "" {
		d.isEncodingConvert = true
		d.encoding = &util.Encoding{}
		if err := d.encoding.Init(args.FromEncoding, args.ToEncoding, args.EncodingError); err != nil {
			log.Errorf("encoding convert init failed, error message: %s", err)
			return err
		}

	} else { // 不进行字符集转换
		d.isEncodingConvert = false
	}
	return nil
}

// 连接数据库
func (d *DB) Connect() (err error) {
	log.Infof("connection database, url: %s", d.args.Url)

	// 解析连接URL
	u, err := url.Parse(d.args.Url)
	if err != nil {
		log.Errorf("parse database url '%s' failed, error message: %s", d.args.Url, err)
		return err
	}

	// 校验数据库
	switch u.Scheme {
	case "hi", "hive": // hive
		dsn := strings.Split(d.args.Url, "://")[1]
		d.conn, err = sql.Open("hive", dsn)
		if err != nil {
			log.Error("connect hive database error, error message: ", err)
			return err
		}

	default: // 其他数据库
		d.conn, err = dburl.Open(d.args.Url)
		if err != nil {
			log.Error("connect database error, error message: ", err)
			return err
		}
	}

	// 立即连接数据库
	log.Info("ping database connection")
	err = d.conn.Ping()
	if err != nil {
		log.Error("test database connection failed, error message: ", err)
		return err
	}
	log.Info("connection database successful")

	return nil
}

// 转换列名称
func (d *DB) columnsNameFmt(columns []string) string {
	rowLength := len(columns) - 1 // 列的长度
	d.recordBuffer.Reset()        // 重置缓冲区

	// 遍历获取列名称
	for i, value := range columns {
		if d.args.Tag != "" { // 输出字段名称两边添加Tag
			d.recordBuffer.WriteString(d.args.Tag)
			d.recordBuffer.WriteString(value)
			d.recordBuffer.WriteString(d.args.Tag)
		} else {
			d.recordBuffer.WriteString(value)
		}

		// 添加字段分隔符
		if i < rowLength {
			d.recordBuffer.WriteString(d.args.Separator)
		}
	}

	// 添加换行符
	d.recordBuffer.WriteString(d.args.Newline)

	return d.recordBuffer.String()
}

// 转换数据
func (d *DB) dataFmt(data []sql.RawBytes, dataType []string) string {
	rowLength := len(data) - 1 // 列的长度
	d.recordBuffer.Reset()     // 重置缓冲

	// 处理数据
	for i, col := range data {
		// 校验输出值是否为空
		if col == nil {
			d.fieldVal = d.args.EmptyVal // 输出值为空, 填写默认值

		} else {
			d.fieldVal = string(col)
		}

		// 字段两边添加标识符
		if d.args.TagAll == true { // 字段两边添加Tag, 不关心是否是数字类型
			d.recordBuffer.WriteString(d.args.Tag)
			d.recordBuffer.WriteString(d.fieldVal)
			d.recordBuffer.WriteString(d.args.Tag)

		} else { // 根据字段类型判断是否需要增加Tag
			if d.whetherAddTag(dataType[i]) == true { // 为True表示需要增加Tag
				d.recordBuffer.WriteString(d.args.Tag)
				d.recordBuffer.WriteString(d.fieldVal)
				d.recordBuffer.WriteString(d.args.Tag)
			} else {
				d.recordBuffer.WriteString(d.fieldVal)
			}
		}

		// 添加字段分隔符
		if i < rowLength {
			d.recordBuffer.WriteString(d.args.Separator)
		}
	}

	// 写入文件换行符
	d.recordBuffer.WriteString(d.args.Newline)

	return d.recordBuffer.String()
}

// 根据传入的字段类型, 决定是否添加Tag, 返回结果为bool
func (d *DB) whetherAddTag(s string) bool {
	// 验证是否包含在用户输入类型中, 类型与用户输入的一致, 返回false
	for _, field := range d.sTagExcludeFieldType {
		if strings.ToUpper(s) == field {
			return false
		}
	}
	return true
}

// 执行查询
func (d *DB) WriteFile() error {
	log.Infof("execution query: %s", d.args.ExecSql)
	// 执行语句
	rows, err := d.conn.Query(d.args.ExecSql)
	if err != nil {
		log.Errorf("execution query failed, error message: %s", err)
		return err
	}
	defer rows.Close()
	log.Info("execution query complete")

	// 获得列名
	columns, err := rows.Columns()
	if err != nil {
		log.Errorf("get result column name failed, error message: %s", err)
		return err
	}
	d.columns = columns

	// 获取列类型
	d.databaseTypeName = []string{} // 置空
	d.scanType = []reflect.Type{}   // 置空
	columnsType, err := rows.ColumnTypes()
	if err != nil {
		log.Errorf("get result column type failed, error message: %s", err)
	}
	for _, col := range columnsType {
		d.databaseTypeName = append(d.databaseTypeName, col.DatabaseTypeName())
		d.scanType = append(d.scanType, col.ScanType())
	}

	// 保存查询对象
	d.rows = rows

	// 为输出结果值创建一个切片
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// 打开输出文件
	fp, err := OpenFile(d.args.SaveFile, d.args.CompressFormat, d.args.OverwriteFile)
	if err != nil {
		return err
	}

	// 判断是否写入列名称, 将列名称写入文件
	if d.args.ColumnName == true {
		// 转换列名称字符集
		n, result := d.encodingConversion(d.columnsNameFmt(d.columns))
		switch n {
		case 0, 1, 2, 3:
			if _, err = fp.WriteString(result); err != nil {
				log.Errorf("write output file date failed, error message: %s", err)
				return err
			}
		}
	}

	// 获取执行结果, 并将查询结果写入文件
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Errorf("query result rows scan error, error message: %s", err)
			return err
		}

		// 转换字符集, 获取写入数据结果
		n, result := d.encodingConversion(d.dataFmt(values, d.databaseTypeName))
		switch n {
		case 0, 1, 2, 3:
			if _, err = fp.WriteString(result); err != nil {
				log.Errorf("write output file date failed, error message: %s", err)
				return err
			}
		}

		// 累计写入行数
		d.flashNum += 1
		d.echoNum += 1
		d.fileNum += 1
		// 写入指定行后刷新到磁盘
		if d.flashNum >= d.args.CacheNumber {
			d.flashNum = 0 // 重置计数
			if err = fp.Flash(); err != nil {
				log.Errorf("output file write flash failed, error message: %s", err)
				return err
			}
		}
		if d.echoNum >= 500000 {
			d.echoNum = 0 // 重置计数
			log.Infof("write output file number %d", d.fileNum)
		}
	}
	if err := rows.Err(); err != nil {
		log.Errorf("get execution result failed, error message: %s", err)
		return err
	}

	// 刷新到文件
	if err = fp.Flash(); err != nil {
		log.Errorf("output file flash failed, error message: %s", err)
		return err
	}

	// 关闭文件
	if err = fp.Close(); err != nil {
		log.Errorf("output file close failed, error message: %s", err)
		return err
	}

	// 记录写入文件结果
	log.Infof("character encoding conversion results, no conversion: %d, "+
		"successful: %d, ignore: %d, skip: %d", d.noCoverCharNum, d.suCoverCharNum,
		d.igCoverCharNum, d.skCoverCharNum)
	log.Infof("write output file total number: %d", d.fileNum)
	return nil
}

// 转换字符串编码
// 返回数字类型结果如下：
// 0 未进行任何操作
// 1 未进行字符编码转换
// 2 字符编码转换成功
// 3 字符编码转换失败忽略
// 4 字符编码转换失败并跳过
func (d *DB) encodingConversion(data string) (int, string) {
	// 判断是否需要转换, 不需要转换则直接返回原先数据
	if d.isEncodingConvert == false {
		d.noCoverCharNum += 1
		return 1, data
	}
	result, i := d.encoding.ConvertEncodingString(data)
	switch i {
	case 1: // 1 未进行字符编码转换
		d.noCoverCharNum += 1
		return 1, result

	case 2: // 2 字符编码转换成功
		d.suCoverCharNum += 1
		return 2, result

	case 3: // 3 字符编码转换失败忽略
		d.igCoverCharNum += 1
		return 3, result

	case 4: // 4 字符编码转换失败并跳过
		d.skCoverCharNum += 1
		return 4, result
	}
	return 0, ""
}

// GetColumnsName 获取列名称
func (d *DB) GetColumnsName() []string {
	return d.columns
}

// GetColumnsType 获取列数据类型
func (d *DB) GetColumnsType() []string {
	return d.databaseTypeName
}

// Close 关闭数据库连接
func (d *DB) Close() error {
	log.Debug("close database connection")
	if err := d.conn.Close(); err != nil {
		log.Errorf("close database connection failed, error message: %s", err)
		return err
	}
	return nil
}
