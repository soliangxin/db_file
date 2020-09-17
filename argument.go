package db_file

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
)

// 需要的参数
type Arguments struct {
	Url                 string // 连接数据库URL, 需要符合URL格式
	QuerySql            string // 查询SQL
	SqlFile             string // SQL文件
	PidFile             string // PID文件, 放置进程多次启动
	Separator           string // 输出文件分隔符
	Newline             string // 输出文件换行符
	SaveFile            string // 输出保存文件名
	OverwriteFile       bool   // 是否覆盖输出文件
	FromEncoding        string // 输入字符集编码
	ToEncoding          string // 输出字符集编码
	EncodingError       string // 字符串编码错误时处理
	Tag                 string // 输出字段标识符
	TagAll              bool   // 所有字段添加标签, 默认只有字符串类型添加
	ColumnName          bool   // 输出文件写入列名
	CompressFormat      string // 输出文件压缩
	CacheNumber         int64  // 最大缓存记录条数
	EmptyVal            string // 字段为空时填写的默认值
	TagExcludeFieldType string // 添加标签排除的字段类型
	LogLevel            string // 日志级别
	ExecSql             string // 最终执行的SQL
}

// 验证接收的参数是否正确
func (a Arguments) Validation() error {
	// URL 未设置
	if a.Url == "" {
		return errors.New(`required flag "url" not set`)
	}

	// 输出保存文件名未设置
	if a.SaveFile == "" {
		return errors.New(`required flag "write" not set`)
	}

	// QuerySql 或者 SqlFile 需要最少指定一个参数
	if a.QuerySql == "" && a.SqlFile == "" {
		return errors.New(`required flag "sql" or "sql-file" not set`)
	}

	// 字符串编码错误时处理
	if govalidator.IsIn(a.EncodingError, "strict", "ignore", "skip") == false {
		return errors.New(fmt.Sprintf(
			`required flag "encoding-error" valid values (strict, ignore, skip), current value: %q`,
			a.EncodingError),
		)
	}

	// 输出文件压缩
	if govalidator.IsIn(a.CompressFormat, "", "gzip", "lz4", "snappy", "zstd") == false {
		return errors.New(fmt.Sprintf(
			`required flag "compress-format" valid values (gzip, lz4, snappy, zstd), current value: %q`,
			a.CompressFormat),
		)
	}

	// 日志级别
	if govalidator.IsIn(a.LogLevel, "trace", "debug", "info", "warn", "error", "fatal", "panic") == false {
		return errors.New(fmt.Sprintf(
			`required flag "level" valid values (trace, debug, info, warn, error, fatal, panic), current value: %q`,
			a.CompressFormat),
		)
	}

	return nil
}
