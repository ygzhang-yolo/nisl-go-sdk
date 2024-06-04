package dataset

import (
	"reflect"
	"testing"
)

/**
 * @Author: ygzhang
 * @Date: 2024/1/12 16:31
 * @Func:
 **/

func TestGetDataFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "TestGetDataFromFile",
			args: args{path: "zipfian.data"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDataFromFile(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
