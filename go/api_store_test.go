package openapi

import (
	"reflect"
	"testing"
)

func TestRemoveTransaction(t *testing.T) {
	type args struct {
		s     []transaction
		index int
	}
	tests := []struct {
		name string
		args args
		want []transaction
	}{
		{name: "Test remove transaction",
			args: args{
				s: []transaction{
					{
						Id: "1",
					},
				},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveTransaction(tt.args.s, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
