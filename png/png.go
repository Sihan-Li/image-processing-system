// Package png allows for loading png images and applying
// image flitering effects on them
package png

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// The Image represents a structure for working with PNG images.
type Image struct {
	in  image.Image
	tmp *image.RGBA64
	Out *image.RGBA64
}

//
// Public functions
//

// Load returns a Image that was loaded based on the filePath parameter
func Load(filePath string) (*Image, error) {

	inReader, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer inReader.Close()

	inImg, err := png.Decode(inReader)

	if err != nil {
		return nil, err
	}

	inBounds := inImg.Bounds()

	outImg := image.NewRGBA64(inBounds)

	tmpImg := image.NewRGBA64(inBounds)

	return &Image{inImg, tmpImg, outImg}, nil
}

func (img *Image)InitTmp(){
	bounds := img.Out.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := img.in.At(x, y).RGBA()
			img.tmp.Set(x, y, color.RGBA64{clamp(float64(r)), clamp(float64(g)), clamp(float64(b)), uint16(a)})
		}
	}
}

//func (img *Image) Helper(){
//	img.tmp = img.out
//
//}

// Save saves the image to the given file
func (img *Image) Save(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img.tmp)
	if err != nil {
		return err
	}
	return nil
}

//clamp will clamp the comp parameter to zero if it is less than zero or to 65535 if the comp parameter
// is greater than 65535.
func clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}
