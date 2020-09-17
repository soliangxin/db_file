package db_file

import (
	"testing"
)

func TestArguments_Validation(t *testing.T) {
	type fields struct {
		Url            string
		QuerySql       string
		SqlFile        string
		PidFile        string
		Separator      string
		Newline        string
		SaveFile       string
		OverwriteFile  bool
		FromEncoding   string
		ToEncoding     string
		EncodingError  string
		Tag            string
		TagAll         bool
		ColumnName     bool
		CompressFormat string
		CacheNumber    int64
		EmptyVal       string
		LogLevel       string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// 测试正确结果
		{name: "correct",
			fields: fields{
				Url:            "ch://user:pass@localhost:port/dbname",
				QuerySql:       "select * from dual",
				SqlFile:        "",
				PidFile:        "",
				Separator:      ";",
				Newline:        `\n`,
				SaveFile:       "/test/result.txt",
				OverwriteFile:  false,
				FromEncoding:   "",
				ToEncoding:     "",
				EncodingError:  "strict",
				Tag:            "'",
				TagAll:         false,
				ColumnName:     false,
				CompressFormat: "",
				CacheNumber:    1000,
				EmptyVal:       `\N`,
				LogLevel:       "debug",
			},
			wantErr: false,
		},

		// 测试SqlFile有值
		{name: "correct-SqlFile",
			fields: fields{
				Url:            "ch://user:pass@localhost:port/dbname",
				QuerySql:       "",
				SqlFile:        "/test/a.sql",
				PidFile:        "",
				Separator:      ";",
				Newline:        `\n`,
				SaveFile:       "/test/result.txt",
				OverwriteFile:  false,
				FromEncoding:   "",
				ToEncoding:     "",
				EncodingError:  "strict",
				Tag:            "'",
				TagAll:         false,
				ColumnName:     false,
				CompressFormat: "",
				CacheNumber:    1000,
				EmptyVal:       `\N`,
				LogLevel:       "debug",
			},
			wantErr: false,
		},

		// 测试Url 为空
		{name: "error-Url",
			fields: fields{
				Url:            "",
				QuerySql:       "select * from dual",
				SqlFile:        "",
				PidFile:        "",
				Separator:      ";",
				Newline:        `\n`,
				SaveFile:       "/test/result.txt",
				OverwriteFile:  false,
				FromEncoding:   "",
				ToEncoding:     "",
				EncodingError:  "strict",
				Tag:            "'",
				TagAll:         false,
				ColumnName:     false,
				CompressFormat: "",
				CacheNumber:    1000,
				EmptyVal:       `\N`,
				LogLevel:       "debug",
			},
			wantErr: true,
		},

		// 测试QuerySql为空
		{name: "error-QuerySql",
			fields: fields{
				Url:            "ch://user:pass@localhost:port/dbname",
				QuerySql:       "",
				SqlFile:        "",
				PidFile:        "",
				Separator:      ";",
				Newline:        `\n`,
				SaveFile:       "/test/result.txt",
				OverwriteFile:  false,
				FromEncoding:   "",
				ToEncoding:     "",
				EncodingError:  "strict",
				Tag:            "'",
				TagAll:         false,
				ColumnName:     false,
				CompressFormat: "",
				CacheNumber:    1000,
				EmptyVal:       `\N`,
				LogLevel:       "debug",
			},
			wantErr: true,
		},

		// 测试SaveFile为空
		{name: "error-SaveFile",
			fields: fields{
				Url:            "ch://user:pass@localhost:port/dbname",
				QuerySql:       "select * from dual",
				SqlFile:        "",
				PidFile:        "",
				Separator:      ";",
				Newline:        `\n`,
				SaveFile:       "",
				OverwriteFile:  false,
				FromEncoding:   "",
				ToEncoding:     "",
				EncodingError:  "strict",
				Tag:            "'",
				TagAll:         false,
				ColumnName:     false,
				CompressFormat: "",
				CacheNumber:    1000,
				EmptyVal:       `\N`,
				LogLevel:       "debug",
			},
			wantErr: true,
		},

		// 测试EncodingError填值
		{name: "error-EncodingError",
			fields: fields{
				Url:            "ch://user:pass@localhost:port/dbname",
				QuerySql:       "select * from dual",
				SqlFile:        "",
				PidFile:        "",
				Separator:      ";",
				Newline:        `\n`,
				SaveFile:       "/test/result.txt",
				OverwriteFile:  false,
				FromEncoding:   "",
				ToEncoding:     "",
				EncodingError:  "",
				Tag:            "'",
				TagAll:         false,
				ColumnName:     false,
				CompressFormat: "",
				CacheNumber:    1000,
				EmptyVal:       `\N`,
				LogLevel:       "debug",
			},
			wantErr: true,
		},

		// 测试输出文件压缩
		{name: "error-CompressFormat",
			fields: fields{
				Url:            "ch://user:pass@localhost:port/dbname",
				QuerySql:       "select * from dual",
				SqlFile:        "",
				PidFile:        "",
				Separator:      ";",
				Newline:        `\n`,
				SaveFile:       "/test/result.txt",
				OverwriteFile:  false,
				FromEncoding:   "",
				ToEncoding:     "",
				EncodingError:  "strict",
				Tag:            "'",
				TagAll:         false,
				ColumnName:     false,
				CompressFormat: "sz",
				CacheNumber:    1000,
				EmptyVal:       `\N`,
				LogLevel:       "debug",
			},
			wantErr: true,
		},

		// 测试日志级别
		{name: "error-LogLevel",
			fields: fields{
				Url:            "ch://user:pass@localhost:port/dbname",
				QuerySql:       "select * from dual",
				SqlFile:        "",
				PidFile:        "",
				Separator:      ";",
				Newline:        `\n`,
				SaveFile:       "/test/result.txt",
				OverwriteFile:  false,
				FromEncoding:   "",
				ToEncoding:     "",
				EncodingError:  "strict",
				Tag:            "'",
				TagAll:         false,
				ColumnName:     false,
				CompressFormat: "",
				CacheNumber:    1000,
				EmptyVal:       `\N`,
				LogLevel:       "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Arguments{
				Url:            tt.fields.Url,
				QuerySql:       tt.fields.QuerySql,
				SqlFile:        tt.fields.SqlFile,
				PidFile:        tt.fields.PidFile,
				Separator:      tt.fields.Separator,
				Newline:        tt.fields.Newline,
				SaveFile:       tt.fields.SaveFile,
				OverwriteFile:  tt.fields.OverwriteFile,
				FromEncoding:   tt.fields.FromEncoding,
				ToEncoding:     tt.fields.ToEncoding,
				EncodingError:  tt.fields.EncodingError,
				Tag:            tt.fields.Tag,
				TagAll:         tt.fields.TagAll,
				ColumnName:     tt.fields.ColumnName,
				CompressFormat: tt.fields.CompressFormat,
				CacheNumber:    tt.fields.CacheNumber,
				EmptyVal:       tt.fields.EmptyVal,
				LogLevel:       tt.fields.LogLevel,
			}
			if err := a.Validation(); (err != nil) != tt.wantErr {
				t.Errorf("Arguments.Validation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
