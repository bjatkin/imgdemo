package shape

import (
	"image"
	"image/color"
)

// Circle represents a cirle that can be drawn on a canvas
type Circle struct {
	Center V2
	Radius float64
	bounds image.Rectangle
}

// NewCircle creates a new circle
func NewCircle(center V2, radius float64) *Circle {
	return &Circle{
		Center: center,
		Radius: radius,
		bounds: image.Rect(
			int(center.X-radius-1.5),
			int(center.Y-radius-1.5),
			int(center.X+radius+1.5),
			int(center.Y+radius+1.5),
		),
	}
}

// forEachPixel runs the given fn for each pixel that is considered to be
// within the radius of the circle
func (c *Circle) forEachPixel(fn func(x, y int)) {
	for y := c.bounds.Min.Y; y < c.bounds.Max.Y; y++ {
		for x := c.bounds.Min.X; x < c.bounds.Max.X; x++ {
			check := NewV2(float64(x), float64(y))
			dist := V2Distance(check, c.Center)
			if dist-1 > c.Radius {
				continue
			}
			fn(x, y)
		}
	}

}

// Render draws the circle onto the destination image using the rgba color. It will use subpixel
// sampling for anti-aliasing
func (c *Circle) Render(rgba color.RGBA, dest *image.RGBA) {
	c.forEachPixel(func(x, y int) {
		density := sampleSubpixel(c, NewV2(float64(x), float64(y)), image.Rect(0, 0, 10, 10))

		if density > 0 {
			bg := dest.RGBAAt(int(x), int(y))
			c := lerpColor(bg, rgba, density)
			dest.Set(int(x), int(y), c)
		}
	})
}

// lerpColor linearly interpolats between c0, and c1 by t where t is between 0 and 1
func lerpColor(c0, c1 color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c0.R) + t*(float64(c1.R)-float64(c0.R))),
		G: uint8(float64(c0.G) + t*(float64(c1.G)-float64(c0.G))),
		B: uint8(float64(c0.B) + t*(float64(c1.B)-float64(c0.B))),
		A: 255,
	}
}

// sampleSubpixel loops through the subpixel and returns the portion of the pixel that is within the radius of the circle
func sampleSubpixel(circle *Circle, pixel V2, subPixels image.Rectangle) float64 {
	delta := NewV2(
		1.0/float64(subPixels.Dx()),
		1.0/float64(subPixels.Dy()),
	)
	covered := 0
	for y := subPixels.Min.Y; y < subPixels.Max.Y; y++ {
		for x := subPixels.Min.X; x < subPixels.Max.X; x++ {
			check := NewV2(pixel.X+delta.X*float64(x), pixel.Y+delta.Y*float64(y))
			dist := V2Distance(check, circle.Center)
			if dist <= circle.Radius {
				covered++
			}
		}
	}

	return float64(covered) / (float64(subPixels.Dx()) * float64(subPixels.Dy()))
}

// Overlap returns the portion of the circle that overlaps with the black pixels in the mask
func (c *Circle) Overlap(mask image.Image) float64 {
	total, covered := 0, 0
	c.forEachPixel(func(x, y int) {
		total++

		mask := mask.At(x, y)
		rgba := color.RGBAModel.Convert(mask).(color.RGBA)
		if rgba.R > 0 || rgba.G > 0 || rgba.B > 0 {
			return
		}

		covered++
	})

	return float64(covered) / float64(total)
}

// Collides returns true if the target circle overlaps with any of the circles provided
func (c *Circle) Collides(circles []Circle, padding float64) bool {
	for _, circle := range circles {
		if CircleCollide(c, &circle, padding) {
			return true
		}
	}
	return false
}

// NewSubCircle generates a new circle with the given radius that passes the success function. It will randomly generate
// sub circles within the radius of the circle until it fails maxTries times at which point it will return (nil, false)
func (c *Circle) NewSubCircle(radius float64, maxTries int, success func(*Circle) bool) (*Circle, bool) {
	for i := 0; ; i++ {
		add := NewCircle(
			polarToCartesian(c, randomPolar(c.Radius-radius)),
			radius,
		)

		if success(add) {
			return add, true
		}

		if i > maxTries {
			return nil, false
		}
	}
}

// CircleCollide returns true if c0 and c1 overlap, the padding is applied only once
func CircleCollide(c0, c1 *Circle, padding float64) bool {
	check := c0.Radius + c1.Radius + padding
	return V2Distance(c0.Center, c1.Center) < check
}
