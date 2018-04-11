package cfg

// import (
// 	"reflect"
// 	"testing"
// )

// func TestCfg(t *testing.T) {
// 	type args struct {
// 		key string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want interface{}
// 	}{
// 		{"nil", args{"nil"}, nil},
// 		{
// 			"test",
// 			args{
// 				"gzhu.icons.ads",
// 			},
// 			[]map[string]string{
// 				map[string]string{"name": "如何正确使用哪有?", "link": "#"},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := Cfg(tt.args.key); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Cfg() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestStrCfg(t *testing.T) {
// 	type args struct {
// 		key string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{"empty", args{"empty"}, ""},
// 		{"test", args{"gzhu.name"}, "广州大学"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := StrCfg(tt.args.key); got != tt.want {
// 				t.Errorf("StrCfg() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
