package hide

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/bjatkin/imgdemo/bits"
	"github.com/bjatkin/imgdemo/cli"
)

// MagicNumber is the magic number that indicates that there is data hidden in this image
var MagicNumber uint16 = 0x1337

// hideArgs are the arguments for the hide command
type hideArgs struct {
	inputPath  string
	dataPath   string
	outputPath string
}

// Cmd is the hide command that hides data in the least significant bit of the given image.
// it will then write a new 'png' image to the output file path with that hidden data.
var Cmd = &cli.Cmd[hideArgs]{
	Name:        "hide",
	Usage:       "hide [INPUT IMAGE PATH] [DATA FILE PATH] [OUTPUT IMAGE PATH]",
	Description: "hide data inside an image using steganography",
	Examples: []cli.Example{
		{
			Description: "hide data from 'secret.dat' in the in 'img.png'",
			Args:        []string{"src.jpeg", "secret.dat", "img.png"},
		},
		{
			Description: "only png files are supported as output files",
			Args:        []string{"src.jpeg", "secret.dat", "img.jpeg"},
			Error:       errors.New("png is the only supported output image format"),
		},
	},
	ParseArgs: func(args []string) (hideArgs, error) {
		if len(args) != 3 {
			return hideArgs{}, errors.New("expected exactly 3 arguments")
		}

		if !strings.HasSuffix(args[2], ".png") {
			return hideArgs{}, errors.New("png is the only supported output image format")
		}

		return hideArgs{
			inputPath:  args[0],
			dataPath:   args[1],
			outputPath: args[2],
		}, nil
	},
	Fn: func(args hideArgs) error {
		imageFile, err := os.Open(args.inputPath)
		if err != nil {
			return fmt.Errorf("failed to read image file: %w", err)
		}
		defer imageFile.Close()

		img, _, err := image.Decode(imageFile)
		if err != nil {
			return fmt.Errorf("failed to decode png file: %w", err)
		}

		rawData, err := os.ReadFile(args.dataPath)
		if err != nil {
			return fmt.Errorf("failed to read in data to encode: %w", err)
		}

		// copy the incomming image into an NRGBA image to make it easy to work with
		rgbaImg := image.NewNRGBA(img.Bounds())
		draw.Draw(rgbaImg, img.Bounds(), img, image.Pt(0, 0), draw.Src)

		hideData(rawData, rgbaImg)

		fout, err := os.Create(args.outputPath)
		if err != nil {
			return fmt.Errorf("failed to open destination file: %w", err)
		}
		defer fout.Close()

		err = png.Encode(fout, rgbaImg)
		if err != nil {
			return fmt.Errorf("failed to encode png output image: %w", err)
		}

		return nil
	},
}

// hideData takes an image and a set of data, it then hides that data in the lowest
// bit of each bye in the image pixel data
func hideData(data []byte, image *image.NRGBA) {
	magicNumber := bits.FromUint16(MagicNumber)
	messageSize := bits.FromUint16(uint16(len(data)))
	messageData := bits.FromBytes(data)
	bitData := append(magicNumber, messageSize...)
	bitData = append(bitData, messageData...)

	for i, b := range bitData {
		if b {
			image.Pix[i] |= 0x01 // set the lowest bit to 1
		} else {
			image.Pix[i] &= 0xFE // set the lowest bit to 0
		}
	}
}
