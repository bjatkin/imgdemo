package bits

import (
	"reflect"
	"strings"
	"testing"
)

func TestFromUint16(t *testing.T) {
	type args struct {
		data uint16
	}
	tests := []struct {
		name string
		args args
		want []bool
	}{
		{
			name: "min",
			args: args{
				data: uint16(0b0000_0000_0000_0000),
			},
			want: []bool{
				false, false, false, false,
				false, false, false, false,
				false, false, false, false,
				false, false, false, false,
			},
		},
		{
			name: "max",
			args: args{
				data: uint16(0b1111_1111_1111_1111),
			},
			want: []bool{
				true, true, true, true,
				true, true, true, true,
				true, true, true, true,
				true, true, true, true,
			},
		},
		{
			name: "mixed",
			args: args{
				data: uint16(0b0111_0110_1011_1011),
			},
			want: []bool{
				false, true, true, true,
				false, true, true, false,
				true, false, true, true,
				true, false, true, true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromUint16(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromUint16() = \n%v, want \n%v", fmtBools(got), fmtBools(tt.want))
			}
		})
	}
}

func TestToUint16(t *testing.T) {
	type args struct {
		bools []bool
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{
			name: "too few bools",
			args: args{
				bools: []bool{true, false},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "max",
			args: args{
				bools: []bool{
					true, true, true, true,
					true, true, true, true,
					true, true, true, true,
					true, true, true, true,
				},
			},
			want:    0b1111_1111_1111_1111,
			wantErr: false,
		},
		{
			name: "min",
			args: args{
				bools: []bool{
					false, false, false, false,
					false, false, false, false,
					false, false, false, false,
					false, false, false, false,
				},
			},
			want:    0b0000_0000_0000_0000,
			wantErr: false,
		},
		{
			name: "mixed",
			args: args{
				bools: []bool{
					false, true, true, false,
					true, false, false, true,
					true, true, false, true,
					true, true, false, false,
				},
			},
			want:    0b0110_1001_1101_1100,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToUint16(tt.args.bools)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToUint16() = \n%08b, want \n%08b", got, tt.want)
			}
		})
	}
}

func TestToFromUint16(t *testing.T) {
	for i := uint16(0); i < 0xFFFF; i++ {
		bools := FromUint16(i)
		got, err := ToUint16(bools)
		if err != nil {
			t.Fatalf("ToFromUint16() failed to convert %d to bool, got error %v", i, err)
		}
		if got != i {
			t.Fatalf("ToFromUint16() failed to convert %d to a bool slice and back", i)
		}
	}
}

func TestFromUint8s(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []bool
	}{
		{
			name: "min",
			args: args{
				data: []byte{0b0000_0000, 0b0000_0000, 0b0000_0000},
			},
			want: []bool{
				false, false, false, false,
				false, false, false, false,
				false, false, false, false,
				false, false, false, false,
				false, false, false, false,
				false, false, false, false,
			},
		},
		{
			name: "max",
			args: args{
				data: []byte{0b1111_1111, 0b1111_1111, 0b1111_1111},
			},
			want: []bool{
				true, true, true, true,
				true, true, true, true,
				true, true, true, true,
				true, true, true, true,
				true, true, true, true,
				true, true, true, true,
			},
		},
		{
			name: "mixed",
			args: args{
				data: []byte{0b1001_0100, 0b1011_0100, 0b1110_0001},
			},
			want: []bool{
				true, false, false, true,
				false, true, false, false,
				true, false, true, true,
				false, true, false, false,
				true, true, true, false,
				false, false, false, true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromBytes(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromUint8s() = \n%v, want \n%v", fmtBools(got), fmtBools(tt.want))
			}
		})
	}
}

func TestToUint8s(t *testing.T) {
	type args struct {
		bools []bool
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "wrong length",
			args: args{
				bools: []bool{true, false},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "single uint8",
			args: args{
				bools: []bool{
					true, true, false, true,
					false, true, false, false,
				},
			},
			want:    []byte{0b1101_0100},
			wantErr: false,
		},
		{
			name: "multiple uint8s",
			args: args{
				bools: []bool{
					true, false, true, false,
					true, true, false, false,
					false, true, true, true,
					true, false, false, true,
					true, true, false, false,
					false, true, false, false,
				},
			},
			want:    []byte{0b1010_1100, 0b0111_1001, 0b1100_0100},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBytes(tt.args.bools)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint8s() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToUint8s() = %08b, want %08b", got, tt.want)
			}
		})
	}
}

func fmtBools(bytes []bool) string {
	s := []string{"["}
	for _, b := range bytes {
		if b {
			s = append(s, "1")
		} else {
			s = append(s, "0")
		}
	}
	s = append(s, "]")
	return strings.Join(s, " ")
}
