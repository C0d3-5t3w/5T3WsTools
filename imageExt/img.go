// Package image extends the functionality of the Go standard image library
// providing additional utility functions for image manipulation and processing.
package imageExt

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

// Resize scales an image to the specified width and height using
// nearest neighbor algorithm (fast but lower quality)
func Resize(img image.Image, width, height int) *image.RGBA {
	bounds := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := bounds.Min.X + x*(bounds.Dx())/width
			srcY := bounds.Min.Y + y*(bounds.Dy())/height
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}

	return dst
}

// ResizeBilinear scales an image to the specified dimensions using bilinear interpolation
// (better quality than nearest neighbor)
func ResizeBilinear(img image.Image, width, height int) *image.RGBA {
	bounds := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	xRatio := float64(bounds.Dx()) / float64(width)
	yRatio := float64(bounds.Dy()) / float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := float64(x) * xRatio
			srcY := float64(y) * yRatio

			x1, y1 := int(srcX), int(srcY)
			x2, y2 := int(math.Min(srcX+1, float64(bounds.Max.X-1))), int(math.Min(srcY+1, float64(bounds.Max.Y-1)))

			xWeight := srcX - float64(x1)
			yWeight := srcY - float64(y1)

			c1 := img.At(x1+bounds.Min.X, y1+bounds.Min.Y)
			c2 := img.At(x2+bounds.Min.X, y1+bounds.Min.Y)
			c3 := img.At(x1+bounds.Min.X, y2+bounds.Min.Y)
			c4 := img.At(x2+bounds.Min.X, y2+bounds.Min.Y)

			r1, g1, b1, a1 := c1.RGBA()
			r2, g2, b2, a2 := c2.RGBA()
			r3, g3, b3, a3 := c3.RGBA()
			r4, g4, b4, a4 := c4.RGBA()

			// Bilinear interpolation
			r := uint8(float64(r1>>8)*(1-xWeight)*(1-yWeight) + float64(r2>>8)*(xWeight)*(1-yWeight) +
				float64(r3>>8)*(1-xWeight)*(yWeight) + float64(r4>>8)*(xWeight)*(yWeight))
			g := uint8(float64(g1>>8)*(1-xWeight)*(1-yWeight) + float64(g2>>8)*(xWeight)*(1-yWeight) +
				float64(g3>>8)*(1-xWeight)*(yWeight) + float64(g4>>8)*(xWeight)*(yWeight))
			b := uint8(float64(b1>>8)*(1-xWeight)*(1-yWeight) + float64(b2>>8)*(xWeight)*(1-yWeight) +
				float64(b3>>8)*(1-xWeight)*(yWeight) + float64(b4>>8)*(xWeight)*(yWeight))
			a := uint8(float64(a1>>8)*(1-xWeight)*(1-yWeight) + float64(a2>>8)*(xWeight)*(1-yWeight) +
				float64(a3>>8)*(1-xWeight)*(yWeight) + float64(a4>>8)*(xWeight)*(yWeight))

			dst.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	return dst
}

// Crop returns a cropped subset of the image
func Crop(img image.Image, rect image.Rectangle) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(dst, dst.Bounds(), img, rect.Min, draw.Src)
	return dst
}

// FlipHorizontal returns a horizontally flipped version of the image
func FlipHorizontal(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(bounds.Max.X-x+bounds.Min.X-1, y, img.At(x, y))
		}
	}

	return dst
}

// FlipVertical returns a vertically flipped version of the image
func FlipVertical(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, bounds.Max.Y-y+bounds.Min.Y-1, img.At(x, y))
		}
	}

	return dst
}

// Grayscale converts an image to grayscale
func Grayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	dst := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, y, img.At(x, y))
		}
	}

	return dst
}

// AdjustBrightness changes the brightness of an image by the given percentage
// percentage ranges from -100 to 100
func AdjustBrightness(img image.Image, percentage float64) *image.RGBA {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	factor := 1.0 + percentage/100.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r = uint32(math.Min(math.Max(float64(r>>8)*factor, 0), 255))
			g = uint32(math.Min(math.Max(float64(g>>8)*factor, 0), 255))
			b = uint32(math.Min(math.Max(float64(b>>8)*factor, 0), 255))

			dst.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a >> 8)})
		}
	}

	return dst
}

// SaveJPEG saves an image to a file in JPEG format with the given quality
func SaveJPEG(img image.Image, filename string, quality int) error {
	if quality < 1 || quality > 100 {
		return errors.New("quality must be between 1 and 100")
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return jpeg.Encode(f, img, &jpeg.Options{Quality: quality})
}

// SavePNG saves an image to a file in PNG format
func SavePNG(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// LoadImage loads an image from a file (supports JPEG and PNG)
func LoadImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

// GetAverageColor returns the average color of an image
func GetAverageColor(img image.Image) color.RGBA {
	bounds := img.Bounds()
	var r, g, b, a uint64
	pixelCount := uint64(bounds.Dx() * bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pr, pg, pb, pa := img.At(x, y).RGBA()
			r += uint64(pr >> 8)
			g += uint64(pg >> 8)
			b += uint64(pb >> 8)
			a += uint64(pa >> 8)
		}
	}

	return color.RGBA{
		uint8(r / pixelCount),
		uint8(g / pixelCount),
		uint8(b / pixelCount),
		uint8(a / pixelCount),
	}
}
