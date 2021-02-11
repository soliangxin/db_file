package util

import "strings"

// EscapeCode 命令行特殊字符进行转义
func EscapeCode(str string) string {
	str = strings.Replace(str, "\\r", "\r", -1)      // Mac OS
	str = strings.Replace(str, "\\t", "\t", -1)      // 跳格
	str = strings.Replace(str, "\\n", "\n", -1)      // linux,unix
	str = strings.Replace(str, "\\r\\n", "\r\n", -1) // windows
	return str
}
