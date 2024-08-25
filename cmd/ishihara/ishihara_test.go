package ishihara

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"reflect"
	"testing"
)

func Test_parseHex(t *testing.T) {
	type args struct {
		hex string
	}
	tests := []struct {
		name    string
		args    args
		want    color.RGBA
		wantErr bool
	}{
		{
			name: "decimal only numbers",
			args: args{
				hex: "#001122",
			},
			want:    color.RGBA{R: 0x00, G: 0x11, B: 0x22},
			wantErr: false,
		},
		{
			name: "with hex numbers",
			args: args{
				hex: "#0A1B2C",
			},
			want:    color.RGBA{R: 0x0A, G: 0x1B, B: 0x2C},
			wantErr: false,
		},
		{
			name: "missing leading #",
			args: args{
				hex: "0A1B2C",
			},
			want:    color.RGBA{},
			wantErr: true,
		},
		{
			name: "hex is too short",
			args: args{
				hex: "#EEE",
			},
			want:    color.RGBA{},
			wantErr: true,
		},
		{
			name: "invalid characters",
			args: args{
				hex: "#BADHEX",
			},
			want:    color.RGBA{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHex(tt.args.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scaleImage(t *testing.T) {
	type args struct {
		destRect image.Rectangle
		src      image.Image
	}
	tests := []struct {
		name string
		args args
		want *image.RGBA
	}{
		{
			name: "2x2 to 4x4",
			args: args{
				destRect: image.Rect(0, 0, 4, 4),
				src:      readTestImage(t, "testdata/2x2.png"),
			},
			want: readTestImage(t, "testdata/4x4.png"),
		},
		{
			name: "4x4 to 2x2",
			args: args{
				destRect: image.Rect(0, 0, 2, 2),
				src:      readTestImage(t, "testdata/4x4.png"),
			},
			want: readTestImage(t, "testdata/2x2.png"),
		},
		{
			name: "scale up orange.png",
			args: args{
				destRect: image.Rect(0, 0, 1450, 1450),
				src:      readTestImage(t, "testdata/orange.png"),
			},
			want: readTestImage(t, "testdata/big_orange.png"),
		},
		{
			name: "scale down orange.png",
			args: args{
				destRect: image.Rect(0, 0, 900, 900),
				src:      readTestImage(t, "testdata/orange.png"),
			},
			want: readTestImage(t, "testdata/small_orange.png"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := scaleImage(tt.args.destRect, tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scaleImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func readTestImage(t *testing.T, path string) *image.RGBA {
	t.Helper()

	f, err := os.Open(path)
	if err != nil {
		t.Fatal("failed to read test image", path, err)
	}

	img, err := png.Decode(f)
	if err != nil {
		t.Fatal("failed to decode test image", path, err)
	}

	imgRGBA, ok := img.(*image.RGBA)
	if !ok {
		t.Fatal("failed to convert image to *image.RGBA", err)
	}

	return imgRGBA
}
