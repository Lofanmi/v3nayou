package gzhu

import (
	"fmt"
	"os"
	"testing"

	"github.com/parnurzeal/gorequest"
)

func TestSchedule(t *testing.T) {
	request := getRequest()
	Login(os.Getenv("TEST_SID"), os.Getenv("TEST_PSW"), request)
	type args struct {
		sid     string
		request *gorequest.SuperAgent
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]string
		wantErr bool
	}{
		{"test", args{os.Getenv("TEST_SID"), request}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Schedule(tt.args.sid, tt.args.request)
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
