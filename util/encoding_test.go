package util

import (
	"github.com/djimenez/iconv-go"
	"testing"
)

func TestEncoding_Init(t *testing.T) {
	type fields struct {
		isConverter    bool
		encodingErrors string
		converter      *iconv.Converter
	}
	type args struct {
		fromEncoding   string
		toEncoding     string
		encodingErrors string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// 测试正确结果
		{
			name: "current",
			args: args{
				fromEncoding:   "GB2312",
				toEncoding:     "UTF8",
				encodingErrors: "strict",
			},
			wantErr: false,
		},

		// 测试错误结果 fromEncoding
		{
			name: "error-fromEncoding",
			args: args{
				fromEncoding:   "GBX",
				toEncoding:     "UTF8",
				encodingErrors: "strict",
			},
			wantErr: true,
		},

		// 测试错误结果 toEncoding
		{
			name: "error-toEncoding",
			args: args{
				fromEncoding:   "GB2312",
				toEncoding:     "GBX",
				encodingErrors: "strict",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Encoding{
				isConverter:    tt.fields.isConverter,
				encodingErrors: tt.fields.encodingErrors,
				converter:      tt.fields.converter,
			}
			if err := e.Init(tt.args.fromEncoding, tt.args.toEncoding, tt.args.encodingErrors); (err != nil) != tt.wantErr {
				t.Errorf("Encoding.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncoding_ConvertEncodingString(t *testing.T) {
	type fields struct {
		isConverter    bool
		encodingErrors string
		converter      *iconv.Converter
	}
	type args struct {
		str string
	}
	// 生成转换对象
	converter, err := iconv.NewConverter("", "")
	if err != nil {
		t.Error("new converter error, error message: ", err)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  int
	}{
		// 测试不进行字符集转换
		{name: "current-no-convert",
			fields: fields{
				isConverter:    false,
				encodingErrors: "strict",
				converter:      converter,
			},
			args: args{
				str: "test",
			},
			want:  "test",
			want1: 1,
		},

		// 测试进行字符集转换
		{name: "current-convert-strict",
			fields: fields{
				isConverter:    true,
				encodingErrors: "strict",
				converter:      converter,
			},
			args: args{
				str: "test",
			},
			want:  "test",
			want1: 2,
		},

		// 测试字符集转换失败忽略
		{name: "current-convert-ignore",
			fields: fields{
				isConverter:    true,
				encodingErrors: "ignore",
				converter:      converter,
			},
			args: args{
				str: "蘴",
			},
			want:  "蘴",
			want1: 3,
		},

		// 测试字符集转换失败忽略
		{name: "current-convert-skip",
			fields: fields{
				isConverter:    true,
				encodingErrors: "skip",
				converter:      converter,
			},
			args: args{
				str: "蘴",
			},
			want:  "",
			want1: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Encoding{
				isConverter:    tt.fields.isConverter,
				encodingErrors: tt.fields.encodingErrors,
				converter:      tt.fields.converter,
			}
			got, got1 := e.ConvertEncodingString(tt.args.str)
			if got != tt.want {
				t.Errorf("Encoding.ConvertEncodingString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Encoding.ConvertEncodingString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestEncoding_Close(t *testing.T) {
	type fields struct {
		isConverter    bool
		encodingErrors string
		converter      *iconv.Converter
	}
	// 生成转换对象
	converter, err := iconv.NewConverter("", "")
	if err != nil {
		t.Error("new converter error, error message: ", err)
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// 测试正确情况
		{name: "current",
			fields: fields{
				isConverter:    false,
				encodingErrors: "strict",
				converter:      converter,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Encoding{
				isConverter:    tt.fields.isConverter,
				encodingErrors: tt.fields.encodingErrors,
				converter:      tt.fields.converter,
			}
			if err := e.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Encoding.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
