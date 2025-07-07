package dcanvas

import (
	"image"
	"log"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
)

type DCanvas struct {
	img         *image.RGBA
	TypeContext *freetype.Context
	Font        *truetype.Font
	Width       int
	Height      int
	X           int
	Y           int
}

// New creates a new canvas with the given width and height.
func New(x, y, width, height int) *DCanvas {
	img := image.NewRGBA(image.Rect(x, y, width, height))

	return &DCanvas{
		Width:  width,
		Height: height,
		X:      x,
		Y:      y,
		img:    img,
	}
}

// WithFreeTypeContext sets the freetype context for the canvas.
// This is required to draw text on the canvas.
func (c *DCanvas) WithFreeTypeContext(fc *freetype.Context) {
	c.TypeContext = fc
}

// WithFont sets the font for the canvas.
// This is required to draw text on the canvas.
func (c *DCanvas) WithFont(f *truetype.Font) {
	c.Font = f
}

// Image returns the image of the canvas.
// This is useful to save the canvas as an image file.
func (c *DCanvas) Image() *image.RGBA {
	return c.img
}

// Size returns the width and height of the canvas.
// This is useful to position elements on the canvas.
func (c *DCanvas) Size() (int, int) {
	return c.Width, c.Height
}

// FreeType returns the freetype context of the canvas.
// This is useful to draw text on the canvas.
func (c *DCanvas) FreeType() *freetype.Context {
	return c.TypeContext
}

// Add image to the canvas.
// This is useful to add images to the canvas.
func (c *DCanvas) AddImage(dest *image.RGBA, destRect image.Rectangle, src image.Image, srcRect image.Rectangle) {
	draw.ApproxBiLinear.Scale(dest, destRect, src, srcRect, draw.Over, nil)
}

// DrawString draws a string on the canvas at the given position.
func (c *DCanvas) DrawString(posX int, posY int, text string) error {
	pt := freetype.Pt(posX, posY)
	_, err := c.FreeType().DrawString(text, pt)
	if err != nil {
		log.Printf("failed to draw text: %v", err)
		return err
	}

	return nil
}

// DrawMultilineString draws a multiline string on the canvas at the given position.
// This function wraps the text to the next line if the text exceeds the width of the canvas.
// The text is split by space and wrapped to the next line if the text exceeds the width of the canvas.
// The text is drawn from the given position (x, y) on the canvas.
// The font size is the size of the font in pixels.
// The line height is the height of each line in pixels.
func (c *DCanvas) DrawMultilineString(posX int, posY int, text string) error {
	lineWidth := 0
	lineHeight := int(c.FreeType().PointToFixed(80.0) >> 6) // convert font size to pixels
	splitStrings := strings.Split(text, " ")
	x := posX
	y := posY + lineHeight

	for _, splitstr := range splitStrings {
		width, _ := c.MeasureText(80.0, splitstr)
		if lineWidth+int(width) < c.Width {
			// stay on existing row
			pt := freetype.Pt(x+lineWidth, y)
			_, err := c.TypeContext.DrawString(splitstr, pt)
			if err != nil {
				log.Printf("failed to draw text: %v", err)
				return err
			}
			lineWidth += int(width)
			lineWidth += int(c.FreeType().PointToFixed(10.0) >> 6) // add margin
		} else {
			// move to next row
			lineWidth = 0
			y += lineHeight
			if y > c.Height {
				// reached the maximum height, stop drawing
				break
			}
			pt := freetype.Pt(x+lineWidth, y)
			_, err := c.FreeType().DrawString(splitstr, pt)
			if err != nil {
				log.Printf("failed to draw text: %v", err)
				return err
			}
			lineWidth += int(width)
			lineWidth += int(c.FreeType().PointToFixed(10.0) >> 6) // add margin
		}
	}

	return nil
}

// MeasureText measures the width and height of the text with the given font size.
// This function returns the width and height of the text in pixels.
func (c *DCanvas) MeasureText(fontSize float64, text string) (float64, float64) {
	face := truetype.NewFace(c.Font, &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	width := 0
	height := 0

	for _, r := range text {
		advance, ok := face.GlyphAdvance(r)
		if !ok {
			advance = 0
		}
		width += advance.Round()
	}

	metrics := face.Metrics()
	height = (metrics.Ascent + metrics.Descent).Round()

	return float64(width), float64(height)
}
