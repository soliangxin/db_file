package util

import "testing"

func TestEscapeCode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "current-\\r",
			args: args{str: "\\r"},
			want: "\r",
		},
		{name: "current-\\r",
			args: args{str: "\r"},
			want: "\r",
		},
		{name: "current-\\t",
			args: args{str: "\\t"},
			want: "\t",
		},
		{name: "current-\\t",
			args: args{str: "\t"},
			want: "\t",
		},
		{name: "current-\\n",
			args: args{str: "\\n"},
			want: "\n",
		},
		{name: "current-\\n",
			args: args{str: "\n"},
			want: "\n",
		},
		{name: "current-\\r\\n",
			args: args{str: "\\r\\n"},
			want: "\r\n",
		},
		{name: "current-\\r\\n",
			args: args{str: "\r\n"},
			want: "\r\n",
		},
		{name: "current-ss",
			args: args{str: "ss"},
			want: "ss",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeCode(tt.args.str); got != tt.want {
				t.Errorf("EscapeCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
