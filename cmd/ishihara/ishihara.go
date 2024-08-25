package ishihara

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/bjatkin/imgdemo/cli"
	"github.com/bjatkin/imgdemo/cmd/ishihara/shape"
)

// ishiharaArgs arg the arguments for the ishihara command
type ishiharaArgs struct {
	primaryColors   []color.RGBA
	secondaryColors []color.RGBA
	maskImagePath   string
	outputImagePath string
}

// Cmd is the ishihara command used to generate ishihara images
var Cmd = &cli.Cmd[ishiharaArgs]{
	Name:        "ishihara",
	Usage:       "ishihara [PRIMARY COLORS] [SECONDARY COLORS] [MASK IMAGE PATH] [OUTPUT IMAGE PATH]",
	Description: "create an ishihara image using the given color pallets and mask image",
	Examples: []cli.Example{
		{
			Description: "create a red green colorblind test image",
			Args:        []string{"3a6a2f,76cd63", "a32222,db5f5f", "mask.png", "red_green.png"},
		},
	},
	ParseArgs: func(args []string) (ishiharaArgs, error) {
		if len(args) != 4 {
			return ishiharaArgs{}, errors.New("expected exactly 4 aguments")
		}

		primary := strings.Split(args[0], ",")
		primaryColors := []color.RGBA{}
		for _, hex := range primary {
			color, err := parseHex(hex)
			if err != nil {
				return ishiharaArgs{}, fmt.Errorf("failed to parse hex color '%s' %w", hex, err)
			}
			primaryColors = append(primaryColors, color)
		}

		secondary := strings.Split(args[1], ",")
		secondaryColors := []color.RGBA{}
		for _, hex := range secondary {
			color, err := parseHex(hex)
			if err != nil {
				return ishiharaArgs{}, fmt.Errorf("failed to parse hex color '%s' %w", hex, err)
			}
			secondaryColors = append(secondaryColors, color)
		}

		if !strings.HasSuffix(args[3], ".png") {
			return ishiharaArgs{}, errors.New("png is the only supported output image format")
		}

		return ishiharaArgs{
			primaryColors:   primaryColors,
			secondaryColors: secondaryColors,
			maskImagePath:   args[2],
			outputImagePath: args[3],
		}, nil
	},
	Fn: func(args ishiharaArgs) error {
		maskFile, err := os.Open(args.maskImagePath)
		if err != nil {
			return fmt.Errorf("failed to read mask image: %v", err)
		}
		defer maskFile.Close()

		mask, _, err := image.Decode(maskFile)
		if err != nil {
			return fmt.Errorf("failed to decode mask image: %v", err)
		}

		imgData := newIshihara(mask)
		img := imgData.render(args.primaryColors, args.secondaryColors)

		outFile, err := os.Create(args.outputImagePath)
		if err != nil {
			return fmt.Errorf("failed to create output image file: %v", err)
		}
		defer outFile.Close()

		return png.Encode(outFile, img)
	},
}

func parseHex(hex string) (color.RGBA, error) {
	if len(hex) != 6 {
		return color.RGBA{}, errors.New("hex color must be in the format #[0-9a-fA-F]{6}")
	}

	rawR := hex[0:2]
	rawG := hex[2:4]
	rawB := hex[4:6]

	r, err := strconv.ParseInt(rawR, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid R channel '%s' must be a valid hex number", rawR)
	}
	g, err := strconv.ParseInt(rawG, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid G channel '%s' must be a valid hex number", rawG)
	}
	b, err := strconv.ParseInt(rawB, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid B channel '%s' must be a valid hex number", rawB)
	}

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b)}, nil
}

type ishihara struct {
	circles []shape.Circle
	mask    *image.RGBA
}

func newIshihara(mask image.Image) ishihara {
	// TODO: create a shape.NewCircle() method
	bounds := shape.Circle{
		Radius: 450,
		Center: shape.V2{X: 512, Y: 512},
	}
	scaledMask := scaleImage(image.Rect(0, 0, 1024, 1024), mask)

	circles := []shape.Circle{}
	// generate a group of circles that align with the scaled mask
	for _, size := range []float64{18, 6, 3} {
		for {
			add, found := bounds.NewSubCircle(size, 10_000, func(c *shape.Circle) bool {
				return !c.Collides(circles, 2) && c.Overlap(scaledMask) > 0.85
			})
			if !found {
				break
			}
			circles = append(circles, *add)
		}
	}

	// now fill the rest of the bounding circle with circles
	for _, size := range []float64{18, 6, 3} {
		for {
			add, found := bounds.NewSubCircle(size, 10_000, func(c *shape.Circle) bool {
				return !c.Collides(circles, 2)
			})
			if !found {
				break
			}
			circles = append(circles, *add)
		}
	}

	return ishihara{
		circles: circles,
		mask:    scaledMask,
	}
}

func (i *ishihara) render(primary, secondary []color.RGBA) image.Image {
	randColor := func(colors []color.RGBA) color.RGBA {
		return colors[rand.Intn(len(colors))]
	}

	img := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	// fill the image with white since it's black by default
	for i := range img.Pix {
		img.Pix[i] = 0xFF
	}

	for _, circle := range i.circles {
		c := randColor(primary)
		if circle.Overlap(i.mask) > 0.85 {
			c = randColor(secondary)
		}

		circle.Render(c, img)
	}

	return img
}

func scaleImage(destRect image.Rectangle, src image.Image) *image.RGBA {
	factorX := float64(src.Bounds().Dx()) / float64(destRect.Dx())
	factorY := float64(src.Bounds().Dy()) / float64(destRect.Dy())

	img := image.NewRGBA(destRect)
	for y := destRect.Min.Y; y < destRect.Max.Y; y++ {
		for x := destRect.Min.X; x < destRect.Max.X; x++ {
			c := src.At(int(float64(x)*factorX), int(float64(y)*factorY))
			img.Set(x, y, c)
		}
	}

	return img
}
