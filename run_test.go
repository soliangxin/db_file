package db_file

import (
	"os"
	"testing"
)

func TestSetLogLevel(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// 正确案例
		{name: "current-trace",
			args: args{
				level: "trace",
			},
		},
		{name: "current-debug",
			args: args{
				level: "debug",
			},
			wantErr: false,
		},
		{name: "current-info",
			args: args{
				level: "info",
			},
			wantErr: false,
		},
		{name: "current-warn",
			args: args{
				level: "warn",
			},
			wantErr: false,
		},
		{name: "current-error",
			args: args{
				level: "error",
			},
			wantErr: false,
		},
		{name: "current-fatal",
			args: args{
				level: "fatal",
			},
			wantErr: false,
		},
		{name: "current-panic",
			args: args{
				level: "panic",
			},
			wantErr: false,
		},
		// 错误测试案例
		{name: "current-test",
			args: args{
				level: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetLogLevel(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("SetLogLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadSQLFile(t *testing.T) {
	// 写入测试数据
	fileName := "test_file.sql"
	fp, err := os.Create(fileName)
	if err != nil {
		t.Fatal("test write sql file error, error message: ", err)
	}
	_, _ = fp.WriteString("# test sql #\n")
	_, _ = fp.WriteString("- test sql -\n")
	_, _ = fp.WriteString("; test sql ;\n")
	_, _ = fp.WriteString("select *\nfrom student;\n")

	// 测试案例
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// 正确案例
		{name: "current", args: args{
			fileName: fileName,
		}, want: "select *\\nfrom student\\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadSQLFile(tt.args.fileName); got != tt.want {
				t.Errorf("ReadSQLFile() = %v, want %v", got, tt.want)
			}
		})
	}
	// 移除SQL文件
	_ = os.Remove(fileName)
}

func TestRun(t *testing.T) {
	type args struct {
		args *Arguments
	}
	tests := []struct {
		name string
		args args
	}{
		// 测试
		{name: "current",
			args: args{
				args: &Arguments{
					Url:                 "my://test:test!123@10.12.1.236:3307/test",
					LogLevel:            "trace",
					PidFile:             "test.pid",
					QuerySql:            "select * from student",
					SaveFile:            "./result.txt",
					OverwriteFile:       true,
					TagAll:              false,
					Tag:                 `'`,
					FromEncoding:        "UTF8",
					ToEncoding:          "GBK",
					TagExcludeFieldType: "INT",
					EmptyVal:            `\N`,
					Separator:           ";",
					Newline:             "\n",
				},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.args)
		})
	}
	// 移除数据文件
	_ = os.Remove("./result.txt")
}
