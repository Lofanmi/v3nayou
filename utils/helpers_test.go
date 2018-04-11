package utils

import "testing"

func TestStrCut(t *testing.T) {
	type args struct {
		s     string
		begin string
		end   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test_1", args{"zh支持中文截取cn", "支持", "截取"}, "中文"},
		{"test_2", args{"zh支持中文截取cn", "我", "文"}, ""},
		{"test_3", args{"zh支持中文截取cn", "支", "我"}, ""},
		{"test_4", args{"zh支持中文截取cn", "取", "支"}, ""},
		{"test_5", args{`img="test"`, `"`, `"`}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrCut(tt.args.s, tt.args.begin, tt.args.end); got != tt.want {
				t.Errorf("StrCut() = %v, want %v", got, tt.want)
			}
		})
	}
}
