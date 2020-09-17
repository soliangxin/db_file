package util

import (
	"github.com/nightlyone/lockfile"
	"os"
	"testing"
)

func TestPidFile_getFileOwner(t *testing.T) {
	type fields struct {
		PidFileName string
	}
	type args struct {
		lock lockfile.Lockfile
	}
	// 测试对象
	fileName := "/go/src/tools/db_file/util/a.pid"
	fp, err := os.Create(fileName)
	if err != nil {
		t.Error("create test file error, error message: ", err)
	}
	fp.Close()
	lock, err := lockfile.New(fileName)
	if err != nil {
		t.Error("new lockfile error, error message: ", err)
	}
	if err := lock.TryLock(); err != nil {
		t.Error("try lock error, error message: ", err)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *os.Process
	}{
		{
			name: "current",
			fields: fields{
				PidFileName: fileName,
			},
			args: args{
				lock: lock,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PidFile{
				PidFileName: tt.fields.PidFileName,
			}
			if got := p.getFileOwner(tt.args.lock); got == nil {
				t.Errorf("PidFile.getFileOwner() = %#v, want %v", got, tt.want)
			}
		})
	}
}

func TestPidFile_Lockfile(t *testing.T) {
	type fields struct {
		PidFileName string
	}
	fileName := "./a.pid"
	tests := []struct {
		name   string
		fields fields
	}{
		// 测试正确结果
		{
			name:   "current",
			fields: fields{PidFileName: fileName},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PidFile{
				PidFileName: tt.fields.PidFileName,
			}
			p.Lockfile()
		})
	}
}

func TestSetPidFile(t *testing.T) {
	type args struct {
		p string
	}
	fileName := "a.pid"
	tests := []struct {
		name string
		args args
	}{
		// 测试正确结果
		{
			name: "current",
			args: args{p: fileName},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetPidFile(tt.args.p)
		})
	}
}

func TestIsConfiguredPidFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "current",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsConfiguredPidFile(); got != tt.want {
				t.Errorf("IsConfiguredPidFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemovePidFile(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "current",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemovePidFile(); (err != nil) != tt.wantErr {
				t.Errorf("RemovePidFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
