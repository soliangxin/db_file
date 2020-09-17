package util

import (
	"errors"
	"fmt"
	"github.com/djimenez/iconv-go"
	log "github.com/sirupsen/logrus"
	"strings"
)

// 字符集转换对象
type Encoding struct {
	isConverter    bool             // 转换字符集
	encodingErrors string           // 异常时处理
	converter      *iconv.Converter // 字符串转换对象
}

// 初始化字符转换对象
func (e *Encoding) Init(fromEncoding, toEncoding, encodingErrors string) (err error) {
	// 字符串编码转换
	if fromEncoding != "" && toEncoding != "" {
		log.Infof("character encoding from %s to %s", fromEncoding, toEncoding)
		e.isConverter = true

		// 初始化字符编码转换对象,  初始化失败, 返回错误信息
		e.converter, err = iconv.NewConverter(fromEncoding, toEncoding)
		if err != nil {
			log.Error("new converter error message: ", err)
			return errors.New(fmt.Sprintf(
				"init character encoding from %s to %s conversion failed", fromEncoding, toEncoding))
		}

	} else { // 不进行转换
		e.isConverter = false
	}
	e.encodingErrors = strings.ToLower(encodingErrors)
	return nil
}

// 字符串编码转换-字符串
//  return
// 1 未进行字符编码转换
// 2 字符编码转换成功
// 3 字符编码转换失败忽略
// 4 字符编码转换失败并跳过
func (e *Encoding) ConvertEncodingString(str string) (string, int) {
	if e.isConverter == true {
		/*
			     记录日志导致性能变慢, 输出文件未压缩情况下, 性能差距为每秒8407
				log.Trace("Convert character encoding records: %s", str)
		*/
		output, err := e.converter.ConvertString(str)
		if err != nil {
			log.Warnf("conversion character encoding failed, source record: %s", str)

			// 字符编码错误时处理
			switch e.encodingErrors {
			case "ignore":
				log.Warn("ignoring character encoding errors, not conversion")
				return str, 3

			case "skip":
				log.Warn("skip character encoding errors, no treatment")
				return "", 4

			default:
				log.Panic("character encoding conversion error")
			}
			// 转换字符集未发生错误, 将输出字符串写入outFileBuffer
		} else {
			return output, 2
		}

	}
	// 不需要转换字符编码, 则将原字符编码返回
	return str, 1
}

// 字符串编码转换
func (e *Encoding) ConvertEncoding(b []byte) []byte {
	if e.isConverter == true {
		//logs.L.Trace("Convert character encoding records: %s", string(b))

		out := make([]byte, len(b)*4)
		_, _, err := e.converter.Convert(b, out)

		if err != nil {
			log.Warnf("query result conversion character encoding failed, source record: %s", string(b))

			// 字符编码错误时处理
			switch e.encodingErrors {
			case "ignore":
				log.Warn("ignoring character encoding errors, not conversion")
				return out

			case "skip":
				log.Warn("skip character encoding errors, no treatment")
				return out

			default:
				log.Panic("character encoding conversion error")
			}
			// 转换字符集未发生错误, 将输出字符串写入outFileBuffer
		} else {
			return out
		}

	}
	// 不需要转换字符编码, 则将原字符编码返回
	return b
}

// 关闭文件对象
func (e *Encoding) Close() (err error) {
	log.Info("close character encoding object")
	return e.converter.Close()
}
