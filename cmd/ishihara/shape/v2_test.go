package shape

import (
	"math"
	"testing"
)

func TestDistance(t *testing.T) {
	type args struct {
		v1 V2
		v2 V2
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "golden 1",
			args: args{
				v1: V2{X: 4, Y: 3},
				v2: V2{X: 0, Y: 0},
			},
			want: 5.0,
		},
		{
			name: "golden 2",
			args: args{
				v1: V2{X: 3, Y: 2},
				v2: V2{X: 9, Y: 7},
			},
			want: math.Sqrt(61),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := V2Distance(tt.args.v1, tt.args.v2); got != tt.want {
				t.Errorf("Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}
