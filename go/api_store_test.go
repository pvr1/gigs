package openapi

import (
	"reflect"
	"testing"
)

func _TestRemoveTransaction(t *testing.T) {
	type arguments struct {
		s     []transaction
		index int
	}
	tests := []struct {
		name string
		args arguments
		want []transaction
	}{
		{name: "Test remove transaction",
			args: arguments{
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
