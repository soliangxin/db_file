package db_file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/soliangxin/db_file/util"
	"io"
	"os"
	"regexp"
	"strings"
)

// SetLogLevel 设置日志级别
func SetLogLevel(level string) error {
	switch strings.ToLower(level) {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		return errors.New(fmt.Sprintf("log level valid values (trace, debug, "+
			"info, warn, error, fatal, panic), current value: %q", level))
	}
	return nil
}

// ReadSQLFile 从文件中获取查询SQL
func ReadSQLFile(fileName string) string {
	// 打开SQL文件
	log.Infof("get the query SQL from the file %s", fileName)
	log.Debugf("open sql file %s", fileName)
	fp, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("open sql file %s failed, error message: %s", fileName, err)
	}
	defer fp.Close()

	// 用于校验字符串首个字符, 不以 #/;/-开头
	re := regexp.MustCompile("^[#|;|-]")

	// 保存读取结果
	var buffer bytes.Buffer

	// 读取SQL文件内容
	log.Debugf("read sql file %s", fileName)
	buf := bufio.NewReader(fp)
	for {
		line, err := buf.ReadString('\n')
		// 读取文件内容错误
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("read sql file %s failed, error message: %s", fileName, err)
			}
		}
		// 处理文件中SQL
		line = strings.TrimSpace(line)

		// 每行不以 #/;/-开头
		if re.MatchString(line) == false {
			line1 := strings.Split(line, ";")[0]
			buffer.WriteString(strings.Split(line1, "#")[0])
			buffer.WriteString(`\n`) // 添加换行符
		}
	}
	sql := buffer.String()
	log.Debugf("get file sql is: %s", sql)

	return sql
}

// 数据库执行操作, 导出数据到文件
func exportToFile(args *Arguments) error {
	// 初始化
	db := &DB{}
	if err := db.Init(args); err != nil {
		return err
	}
	// 连接数据库
	if err := db.Connect(); err != nil {
		return err
	}
	// 执行导出数据
	if err := db.WriteFile(); err != nil {
		return err
	}
	// 关闭数据库连接
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}

// Run 接收命令行参数后, 处理参数
func Run(args *Arguments) {
	// 设置日志格式
	if err := SetLogLevel(args.LogLevel); err != nil {
		log.Fatalf("setting log level failed, error message: %s", err)
	}

	// PID文件不为空, 设置PID文件, 锁定PID文件, 防止进程重复启动
	if args.PidFile != "" {
		p := util.PidFile{
			PidFileName: args.PidFile,
		}
		p.Lockfile()               // 锁定PID文件
		defer util.RemovePidFile() // 移除PID文件
	}

	// 获取数据库执行的SQL语句
	if args.QuerySql != "" { // 从命令行出去SQL语句
		args.ExecSql = args.QuerySql

	} else { // 从文件中读取SQL语句
		args.ExecSql = ReadSQLFile(args.SqlFile)
	}

	// 执行数据库操作, 将数据导出到文件
	if err := exportToFile(args); err != nil {
		log.Fatalf("exec database dat export to file error, error message: %s", err)
	}
}
