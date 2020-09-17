package db_file

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/soliangxin/db_file/file"
	"io"
	"os"
	"strings"
)

type Write interface {
	NewWriter(w io.Writer) error       // 打开文件写入
	Write(p []byte) (int, error)       // 写入[]byte
	WriteString(s string) (int, error) // 写入字符串
	Reset(w io.Writer)                 // 重置
	Flush() error                      // 刷新数据到磁盘
	Close() error                      // 关闭文件
}

// 打开压缩文件Write
// 文件名称、模式、权限、压缩类型
func NewWrite(w io.Writer, compressType string) (Write, error) {
	var (
		wd  Write
		err error
	)
	// 打开文件
	log.Trace("new compress write type ", compressType)
	// 校验文件压缩类型, 并打开不同的文件类型
	switch strings.ToLower(compressType) {
	case "": // 写入文件不使用压缩
		wd = &file.NoCompression{}

	case "gz", "gzip": // Gzip
		wd = &file.GzipWrite{}

	case "lz", "lz4": // Lz4
		wd = &file.Lz4Write{}

	case "sz", "snappy": // Snappy
		wd = &file.SnappyWrite{}

	case "zst", "zstd": // Zstd
		wd = &file.ZstdWrite{}

	default:
		return nil, errors.New(fmt.Sprintf("write interface not compress type %s", compressType))
	}
	// 初始化 Reader
	err = wd.NewWriter(w)
	if err != nil {
		log.Error("new compress write error, error message: ", err)
		return nil, err
	}
	return wd, nil
}

// 写入输出文件内容
type File struct {
	fileName string   // 文件名
	file     *os.File // 保存打开的文件对象
	write    Write    // 写入文件对象
}

// 打开输出文件
// 参数：输出文件名, 文件压缩格式, 文件存在是否覆盖
// 返回: *File, error
func OpenFile(fileName, compressFormat string, overwriteFile bool) (*File, error) {
	log.Infof("open output file %s", fileName)
	var (
		fp  *os.File // 打开文件对象
		err error    // 保存错误信息
	)

	// 判断是否覆盖输出文件
	if overwriteFile == true { // 打开并清空输出文件
		fp, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0775))
		if err != nil {
			log.Errorf("open output file %s error, error message: %s", fileName, err)
			return nil, err
		}

	} else { // 追加文件模式
		fp, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0775))
	}
	if err != nil {
		log.Errorf("open output file %s error, error message: %s", fileName, err)
		return nil, err
	}

	// 使用Write打开
	log.Debugf("write using compressed format %q", compressFormat)
	write, err := NewWrite(fp, compressFormat)
	if err != nil {
		log.Errorf("write using compressed format %q, error message: %s", compressFormat, err)
	}

	return &File{fileName: fileName, write: write, file: fp}, nil
}

// 写入数据
func (f *File) Write(b []byte) (int, error) {
	n, err := f.write.Write(b)
	if err != nil {
		log.Errorf("data is flushed to disk failed, error message: %s", err)
		return 0, err
	}
	return n, nil
}

// 写入数据
func (f *File) WriteString(s string) (int, error) {
	n, err := f.write.WriteString(s)
	if err != nil {
		log.Errorf("data is flushed to disk failed, error message: %s", err)
		return 0, err
	}
	return n, nil
}


// 刷新数据到磁盘
func (f *File) Flash() error {
	return f.write.Flush()
}

// 关闭文件
func (f *File) Close() error {
	// 关闭写入压缩流
	log.Debugf("close compressed flow")
	err := f.write.Close()
	if err != nil {
		log.Errorf("close compressed flow error, error message: %s", err)
		return err
	}

	// 关闭文件
	log.Infof("close output file %s", f.fileName)
	err = f.file.Close()
	if err != nil {
		log.Errorf("close output file %s failed, error message: %s", f.fileName, err)
		return err
	}
	return nil
}
