package quadtreego

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func HLine(x1, y, x2 int, img image.RGBA, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a veritcal line
func VLine(x, y1, y2 int, img image.RGBA, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func Rect(x1, y1, x2, y2 int, img image.RGBA, col color.Color) {
	HLine(x1, y1, x2, img, col)
	HLine(x1, y2, x2, img, col)
	VLine(x1, y1, y2, img, col)
	VLine(x2, y1, y2, img, col)
}

func drawRect(x1 int, y1 int, x2 int, y2 int,img *image.RGBA) {

	var col color.Color

	col = color.RGBA{0, 0, 255, 255} // blue
	Rect(x1, y1, x2, y2, *img, col)


}
func drawPoint(x, y int,img *image.RGBA)  {
	img.Set(x, y, color.RGBA{255, 0, 0, 255})
}

func render(img *image.RGBA,name string){
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
