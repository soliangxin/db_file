package util

import (
	"errors"
	"github.com/nightlyone/lockfile"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type PidFile struct {
	PidFileName string // Pid 文件名称
}

// 获取Pid文件所有者
func (p *PidFile) getFileOwner(lock lockfile.Lockfile) *os.Process {
	owner, err := lock.GetOwner()
	if err != nil {
		log.Fatal("get pid file owner failed, error message: ", err)
	}
	return owner
}

// 锁定进程Pid文件
func (p *PidFile) Lockfile() {
	// 初始化, 获取Pid文件
	var err error
	p.PidFileName, err = filepath.Abs(p.PidFileName)
	if err != nil {
		log.Fatalf("gets pid file %s abs path failed", p.PidFileName)
	}
	log.Infof("pid file %s", p.PidFileName)
	pidDir := filepath.Dir(p.PidFileName)

	// 创建Pid文件
	err = os.MkdirAll(pidDir, os.FileMode(0775))
	if err != nil {
		log.Fatalf("create pid dir %s failed", pidDir)
	}

	// 初始化Pid文件
	lock, err := lockfile.New(p.PidFileName)
	if err != nil {
		log.Fatalf("init lock pid file %s failed", p.PidFileName)

	} else {
		// 锁定文件
		err = lock.TryLock()
		if err != nil {
			owner := p.getFileOwner(lock)
			log.Errorf("daemon lock pid file %s failed", p.PidFileName)
			log.Errorf("process is already running..., pid: %d", owner.Pid)
			SetPidFile("")
		}
	}

	// 获取文件所有者, 等待进程退出
	owner := p.getFileOwner(lock)
	err = owner.Release()
	if err != nil {
		log.Errorf("daemon release pid file %s failed", p.PidFileName)
	}

	// 设置Pid文件
	SetPidFile(p.PidFileName)
}

var (
	errNotConfigured = errors.New("pidfile not configured")
	isPid            = false // 是否使用PID文件
	pidfile          *string
)

// 设置Pid文件
func SetPidFile(p string) {
	pidfile = &p
	isPid = true
}

// 返回是否已设置Pid文件
func IsConfiguredPidFile() bool {
	return isPid
}

// 删除Pid文件
func RemovePidFile() error {
	log.Debugf("remove pid file %s", *pidfile)
	if *pidfile == "" {
		return errNotConfigured
	}
	// 删除PID文件
	if err := os.Remove(*pidfile); err != nil {
		log.Errorf("remove pid file %s failed, error message: %s", *pidfile, err)
		return err
	}
	return nil
}
