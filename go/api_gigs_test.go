package openapi

import (
	"reflect"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	type args struct {
		a interface{}
		b interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test deep copy",
			args: args{
				a: &Gig{
					Id:                "1",
					Category:          Category{},
					Name:              "Gig 1",
					Description:       []string{"Description 1"},
					Measurableoutcome: []string{"Measurable outcome 1"},
					Tags:              []Tag{},
					Status:            "available",
				},
				b: &Gig{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeepCopy(tt.args.a, tt.args.b)
			if !reflect.DeepEqual(tt.args.a, tt.args.b) {
				t.Errorf("DeepCopy() = %v, want %v", tt.args.b, tt.args.a)
			}
		})
	}
}
