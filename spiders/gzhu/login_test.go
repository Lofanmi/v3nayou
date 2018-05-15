package gzhu

import (
	"os"
	"testing"

	"github.com/parnurzeal/gorequest"
)

func TestLogin(t *testing.T) {
	req := getRequest()
	type args struct {
		sid     string
		psw     string
		request *gorequest.SuperAgent
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"name_1", args{os.Getenv("TEST_SID"), os.Getenv("TEST_PSW"), req}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := Login(tt.args.sid, tt.args.psw, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
