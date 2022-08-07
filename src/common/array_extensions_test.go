package common

import (
	"reflect"
	"testing"
)

func TestToEnumerableIds(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    []int
		wantErr bool
	}{
		{name: "Happy#1", args: "1,2,3", want: []int{1, 2, 3}, wantErr: false},
		{name: "Happy#2", args: "1,2,,2,3", want: []int{1, 2, 3}, wantErr: false},
		{name: "Happy#3", args: "1,2,,2,", want: []int{1, 2}, wantErr: false},
		{name: "Error#1", args: "a,2,,2,3", want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToEnumerableIds(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToEnumerableIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToEnumerableIds() got = %v, want %v", got, tt.want)
			}
		})
	}
}
