package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/soliangxin/db_file"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
)

// 程序版本
const appVersion = "1.1.0"

// Usage 命令行提示用法
const Usage = `NAME:
   {{.Name}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}{{end}}{{if .Commands}}[command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}

AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Version}}
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
   {{.Version}}
   {{end}}
`

func main() {
	// 设置命令行显示模板
	cli.AppHelpTemplate = Usage

	// 命令行参数保存结构体
	args := &db_file.Arguments{}

	// 初始化命令行参数
	app := &cli.App{
		// 设置命令行参数
		Flags: []cli.Flag{
			// 数据库连接字符串
			&cli.StringFlag{
				Name:        "url",
				Aliases:     []string{"u"},
				Usage:       "database connection url",
				Required:    false,
				Destination: &args.Url,
			},

			// 命令行参数获取SQL
			&cli.StringFlag{
				Name:        "sql",
				Usage:       "execute the exported SQL statement",
				Required:    false,
				Destination: &args.QuerySql,
			},

			// 从文件中获取SQL
			&cli.StringFlag{
				Name:        "sql_file",
				Aliases:     []string{"f"},
				Usage:       "get the sql from the file",
				Required:    false,
				Destination: &args.SqlFile,
			},

			// PID文件
			&cli.StringFlag{
				Name:        "pid",
				Aliases:     []string{"p"},
				Usage:       "pid file to prevent multiple process starts",
				Required:    false,
				Destination: &args.PidFile,
			},

			// 输出文件分隔符
			&cli.StringFlag{
				Name:        "sep",
				Aliases:     []string{"s"},
				Usage:       "output file separator",
				Required:    false,
				Value:       ";",
				Destination: &args.Separator,
			},

			// 输出文件换行符
			&cli.StringFlag{
				Name:        "newline",
				Aliases:     []string{"n"},
				Usage:       "output file newline character",
				Required:    false,
				Value:       "\n",
				Destination: &args.Newline,
			},

			// 输出保存文件名
			&cli.StringFlag{
				Name:        "write",
				Aliases:     []string{"w"},
				Usage:       "output file name",
				Required:    false,
				Destination: &args.SaveFile,
			},

			// 是否覆盖输出文件
			&cli.BoolFlag{
				Name:        "overwrite",
				Aliases:     []string{"o"},
				Usage:       "if the output file exists, whether to overwrite the file",
				Value:       false,
				Destination: &args.OverwriteFile,
			},

			// 输入字符集编码
			&cli.StringFlag{
				Name:        "from_encoding",
				Usage:       "input character set encoding",
				Required:    false,
				Destination: &args.FromEncoding,
			},

			// 输出字符集编码
			&cli.StringFlag{
				Name:        "to_encoding",
				Usage:       "output character set encoding",
				Required:    false,
				Destination: &args.ToEncoding,
			},

			// 字符串编码错误时处理
			&cli.StringFlag{
				Name:        "encoding_error",
				Usage:       "conversion coding error, valid values (strict, ignore, skip)",
				Required:    false,
				Value:       "strict",
				Destination: &args.EncodingError,
			},

			// 输出字段标识符
			&cli.StringFlag{
				Name:        "tag",
				Aliases:     []string{"t"},
				Usage:       "add a tag to the output field",
				Required:    false,
				Value:       "\"",
				Destination: &args.Tag,
			},

			// 字段是否都加标签Tag
			&cli.BoolFlag{
				Name:        "tag_all",
				Usage:       "all fields are added with tags, and the default non-numeric type is added",
				Value:       false,
				Destination: &args.TagAll,
			},

			// 添加Tag排除的字段类型
			&cli.StringFlag{
				Name:        "tag_exclude",
				Usage:       `The database type excluded when adding the tag, with multiple types separated by ","`,
				Value:       "INT,BIGINT",
				Destination: &args.TagExcludeFieldType,
			},

			// 输出文件写入列名
			&cli.BoolFlag{
				Name:        "column_name",
				Usage:       "output file writes column names",
				Required:    false,
				Value:       false,
				Destination: &args.ColumnName,
			},

			// 输出文件压缩
			&cli.StringFlag{
				Name:        "compress_format",
				Usage:       "output file compression, valid values (gzip, lz4, snappy, zstd)",
				Required:    false,
				Destination: &args.CompressFormat,
			},

			// 最大缓存记录条数
			&cli.Int64Flag{
				Name:        "cache_num",
				Usage:       "the maximum number of records to write to the cache",
				Required:    false,
				Value:       1000,
				Destination: &args.CacheNumber,
			},

			// 字段为空时填写的默认值
			&cli.StringFlag{
				Name:        "empty_val",
				Usage:       "the value filled in when the field value is NULL",
				Required:    false,
				Destination: &args.EmptyVal,
			},

			// 日志级别
			&cli.StringFlag{
				Name:        "level",
				Aliases:     []string{"l"},
				Usage:       "current console log level, valid values (trace, debug, info, warn, error, fatal, panic)",
				Required:    false,
				Value:       "info",
				Destination: &args.LogLevel,
			},
		},

		// 验证命令行参数
		Action: func(c *cli.Context) error {
			// 没有收到命令行参数, 打印帮助信息并退出
			if c.NumFlags() < 1 {
				cli.ShowAppHelpAndExit(c, 0)
			}

			// 运行命令, 导出数据
			db_file.Run(args)

			return nil
		},
	}

	// 命令行参数排序
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	// 程序版本
	app.Version = appVersion

	// 运行命令行程序, 如果有错误, 打印错误信息并退出
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化, 设置日志格式
func init() {
	// 设置日志格式
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
}
