package util

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var SignalMap map[string]int

func init() {
	SignalMap = make(map[string]int)
	SignalMap["hangup"] = 1
	SignalMap["interrupt"] = 2
	SignalMap["quit"] = 3
	SignalMap["illegal instruction"] = 4
	SignalMap["trace/breakpoint trap"] = 5
	SignalMap["aborted"] = 6
	SignalMap["bus error"] = 7
	SignalMap["floating point exception"] = 8
	SignalMap["killed"] = 9
	SignalMap["user defined exit 1"] = 10
	SignalMap["segmentation fault"] = 11
	SignalMap["user defined exit 2"] = 12
	SignalMap["broken pipe"] = 13
	SignalMap["alarm clock"] = 14
	SignalMap["terminated"] = 15
}

// 监听并处理Linux系统信号
func ListenSignal() {
	log.Info("receive and process signals...")
	//创建监听退出chan
	c := make(chan os.Signal)

	//监听指定信号 ctrl+c kill
	signal.Notify(c,
		syscall.SIGHUP, // 终端挂起或者控制进程终止        1
		// syscall.SIGINT,  // 键盘的退出键被按下(ctrl + c)   2
		// syscall.SIGQUIT, // 用户发送QUIT字符(Ctrl+/)触发   3
		// syscall.SIGTERM, // 终止信号, kill                15
	)

	// 启动监听
	go func() {
		for s := range c {
			switch s {
			default:
				sigDesc := s.String()
				sig := SignalMap[sigDesc]
				log.Infof("received exit {%d:%s}, ignore this exit", sig, sigDesc)
			}
		}
	}()
}
