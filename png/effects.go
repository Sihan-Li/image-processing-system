// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image"
	"image/color"
	"math"
)

func (img *Image) Grayscale() {
	bounds := img.Out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.tmp.At(x, y).RGBA()
			greyC := clamp(float64(r+g+b) / 3)
			img.Out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
	img.tmp = img.Out
	img.Out = image.NewRGBA64(bounds)
}

func (img *Image)Sharpen() [][]float64{
	kernelReflect := [][] float64{
		{0,-1,0},
		{-1,5,-1},
		{0,-1,0},
	}
	return kernelReflect
}

func (img *Image)Edge_detection() [][]float64{
	kernelReflect := [][] float64{
		{-1,-1,-1},
		{-1,8,-1},
		{-1,-1,-1},
	}
	return kernelReflect
}

func (img *Image)Blur() [][]float64{
	kernelReflect := [][] float64{
		{float64(1)/float64(9),float64(1)/float64(9),float64(1)/float64(9)},
		{float64(1)/float64(9),float64(1)/float64(9),float64(1)/float64(9)},
		{float64(1)/float64(9),float64(1)/float64(9),float64(1)/float64(9)},
	}
	return kernelReflect
}

// Convolution a image to a given kernel
func (img *Image) ConvertImage(kernelReflect [][]float64) {
	bounds := img.Out.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			apparent := uint32(0)
			red := float64(0.0)
			green := float64(0.0)
			blue := float64(0.0)
			for m := 0; m < 3; m++ {
				for n := 0; n < 3; n++ {
					xx := x + m - 1
					yy := y + n -1
					if xx >= 0 && xx < bounds.Max.X && yy >= 0 && yy < bounds.Max.Y {
						r, g, b, a := img.tmp.At(xx, yy).RGBA()
						red += float64(r) * kernelReflect[m][n]
						green += float64(g) * kernelReflect[m][n]
						blue += float64(b) * kernelReflect[m][n]
						apparent = a
					}
				}
			}
			img.Out.Set(x, y, color.RGBA64{clamp(float64(red)), clamp(float64(green)), clamp(float64(blue)), clamp(float64(apparent))})
		}
	}
	img.tmp = img.Out
	img.Out = image.NewRGBA64(bounds)
}

func (img *Image) DivideImage(kernelReflect [][]float64,numberOfThreads int, lengthOfNormalPart int) {
	bounds := img.Out.Bounds()
	start := bounds.Min.X
	width := bounds.Max.X
	end := width

	flag := make(chan string,numberOfThreads+1)
	for i:= 0; i < numberOfThreads+1; i++ {
		end = int(math.Min(float64(start+lengthOfNormalPart), float64(width)))
		tempStart := start
		tempEnd := end
		go func() {
			for x := tempStart; x < tempEnd; x++ {
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					apparent := uint32(0)
					red := float64(0.0)
					green := float64(0.0)
					blue := float64(0.0)

					for m := 0; m < 3; m++ {
						for n := 0; n < 3; n++ {
							xx := x + m - 1
							yy := y + n - 1
							if xx >= 0 && xx < bounds.Max.X && yy >= 0 && yy < bounds.Max.Y {
								r, g, b, a := img.tmp.At(xx, yy).RGBA()
								red += float64(r) * kernelReflect[m][n]
								green += float64(g) * kernelReflect[m][n]
								blue += float64(b) * kernelReflect[m][n]
								apparent = a
							}
						}
					}
					img.Out.Set(x, y, color.RGBA64{clamp(red), clamp(green), clamp(blue), clamp(float64(apparent))})
				}
			}
			flag <- "task is finished"
		}()
		start = end
	}

	for i:=0;i<numberOfThreads+1;i++{
		<-flag
	}
	img.tmp = img.Out
	img.Out = image.NewRGBA64(bounds)
}

func (img *Image) DivideGrayscale(numberOfThreads int,lengthOfNormalPart int) {
	bounds := img.Out.Bounds()
	start := bounds.Min.X
	width := bounds.Max.X
	end := width

	flag := make(chan string, numberOfThreads+1)
	for i:= 0; i < numberOfThreads+1; i++ {
		end = int(math.Min(float64(start+lengthOfNormalPart), float64(width)))
		tempStart := start
		tempEnd := end

		go func() {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := tempStart; x < tempEnd; x++ {
					r, g, b, a := img.tmp.At(x, y).RGBA()
					greyC := clamp(float64(r+g+b) / 3)
					img.Out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
				}
			}
			flag <- "task finished"
		}()
		start = end
	}
	for i:=0;i<numberOfThreads+1;i++{
		<-flag
	}

	img.tmp = img.Out
	img.Out = image.NewRGBA64(bounds)
}

