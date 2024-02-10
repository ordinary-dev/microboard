package captcha

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const (
	WIDTH     = 120
	HEIGHT    = 50
	FONT_SIZE = 24
)

var (
	preloadedFont *truetype.Font
)

// Generate distorted image with the specified input.
func GenerateCaptcha(text string) (image.Image, error) {
	shape := image.Rectangle{image.Point{0, 0}, image.Point{WIDTH, HEIGHT}}
	img := image.NewRGBA(shape)

	for i := 0; i < HEIGHT; {
		bg := GetRandomBgColor()
		blockHeight := 6 + rand.Intn(8)
		yEnd := i + blockHeight
		if yEnd > HEIGHT {
			yEnd = HEIGHT
		}
		draw.Draw(img, image.Rectangle{image.Point{0, i}, image.Point{WIDTH, yEnd}}, bg, image.Point{}, draw.Src)
		i += blockHeight
	}

	if preloadedFont == nil {
		if err := PreloadFont(); err != nil {
			return nil, err
		}
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(preloadedFont)
	c.SetFontSize(FONT_SIZE)
	c.SetClip(img.Bounds())
	c.SetDst(img)

	for i, ch := range strings.Split(text, "") {
		c.SetSrc(GetRandomTextColor())

		x := (i+1)*15 + rand.Intn(8) - 4
		y := rand.Intn(10) + 30
		pt := freetype.Pt(x, y)

		_, err := c.DrawString(ch, pt)
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

// Load font file into global variable.
// Will be called automatically if `preloadedFont` is nil.
func PreloadFont() error {
	fontBytes, err := os.ReadFile("assets/fonts/Inter-Regular.ttf")
	if err != nil {
		return err
	}

	preloadedFont, err = freetype.ParseFont(fontBytes)
	return err
}

func GetRandomTextColor() *image.Uniform {
	red := rand.Intn(156) + 100
	green := rand.Intn(156) + 100
	blue := rand.Intn(156) + 100
	return image.NewUniform(color.RGBA{uint8(red), uint8(green), uint8(blue), 255})
}

func GetRandomBgColor() *image.Uniform {
	red := rand.Intn(40)
	green := rand.Intn(40)
	blue := rand.Intn(40)
	return image.NewUniform(color.RGBA{uint8(red), uint8(green), uint8(blue), 255})
}
