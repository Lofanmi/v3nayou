package gzhu

import (
	"os"
	"testing"

	"github.com/parnurzeal/gorequest"
)

func TestScore(t *testing.T) {
	request := getRequest()
	Login(os.Getenv("TEST_SID"), os.Getenv("TEST_PSW"), request)
	type args struct {
		sid     string
		xn      string
		xq      string
		request *gorequest.SuperAgent
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]string
		wantErr bool
	}{
		{"test", args{os.Getenv("TEST_SID"), "2016-2017", "1", request}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Score(tt.args.sid, tt.args.xn, tt.args.xq, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Score() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
