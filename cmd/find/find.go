package find

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/bjatkin/imgdemo/bits"
	"github.com/bjatkin/imgdemo/cli"
	"github.com/bjatkin/imgdemo/cmd/hide"
)

// findArgs are the arguments for the find command
type findArgs struct {
	imagePath string
}

// Cmd is the find command that searches the given image for data hidden using the hide command
var Cmd = &cli.Cmd[findArgs]{
	Name:        "find",
	Usage:       "find [IMAGE PATH]",
	Description: "find data hidden inside an image",
	Examples: []cli.Example{
		{
			Description: "find hidden data from inside 'img.png'",
			Args:        []string{"img.png"},
			Output:      "Here's the hidden data",
		},
		{
			Description: "searching for hidden data in 'img.png' fails",
			Args:        []string{"img.png"},
			Error:       errors.New("magic number does not match"),
		},
	},
	ParseArgs: func(args []string) (findArgs, error) {
		if len(args) != 1 {
			return findArgs{}, errors.New("invalid argument count")
		}

		if !strings.HasSuffix(args[0], ".png") {
			return findArgs{}, errors.New("only png images are supported")
		}

		return findArgs{
			imagePath: args[0],
		}, nil
	},
	Fn: func(args findArgs) error {
		f, err := os.Open(args.imagePath)
		if err != nil {
			return fmt.Errorf("failed to read in an image file: %w", err)
		}
		defer f.Close()

		inImage, err := png.Decode(f)
		if err != nil {
			return fmt.Errorf("failed to decode png file: %w", err)
		}

		img, ok := inImage.(*image.NRGBA)
		if !ok {
			return fmt.Errorf("invalid image format: %T", img)
		}

		got, err := findData(img)
		if err != nil {
			return fmt.Errorf("failed to get hidden data: %w", err)
		}

		fmt.Print(string(got))
		return nil
	},
}

// findData searches an image for data hidden using the hide command
// it first checks for the MagicNumber, if it's found it then pulls the
// data from the lowest bits of the image
func findData(image *image.NRGBA) ([]byte, error) {
	// check for the magic number
	var magic []bool
	for i := 0; i < 16; i++ {
		magic = append(magic, image.Pix[i]&0x01 == 1)
	}

	number, err := bits.ToUint16(magic)
	if err != nil {
		return nil, fmt.Errorf("failed to decode magic number: %w", err)
	}
	if number != hide.MagicNumber {
		return nil, errors.New("magic number does not match")
	}

	// get the size of the message
	var size []bool
	for i := 16; i < 32; i++ {
		size = append(size, image.Pix[i]&0x01 == 1)
	}

	dataLen, err := bits.ToUint16(size)
	if err != nil {
		return nil, errors.New("failed to get data length")
	}

	// now get the hidden data from the lowest bits
	data := []bool{}
	for i := 32; i < int(dataLen*8)+32; i++ {
		data = append(data, image.Pix[i]&0x01 == 1)
	}

	return bits.ToBytes(data)
}
