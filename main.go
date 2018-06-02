package main

import "os"
import (
	"crypto/md5"
	"image/color"
	"image/png"
	"image"
)

type UserData struct {
name string
hashValue [16]byte
}

type Identicon struct {
	pixelColor color.RGBA
	pictureMatrix [25]byte
	pixelMatrix [25]bool
	pixelWidth int
	pixelHeight int
}
func createMatrix(hashMatrix [16]byte) [25]byte {
	var out [25]byte
	j := 0
	for i := 0; i < len(hashMatrix) - 1; i += 3 {
			out[j] = hashMatrix[i]
			out[j+1] = hashMatrix[i+1]
			out[j+2] = hashMatrix[i+2]
			out[j+3] = hashMatrix[i+1]
			out[j+4] = hashMatrix[i]
			j += 5
	}
	return out
}

func (i *Identicon) createPixelMatrix(){
	var out [25]bool
	for i,data := range i.pictureMatrix {
		if data % 2 == 1 {
			out[i] = false
		} else {
			out[i] = true
		}
	}
	 i.pixelMatrix = out
}

type Rect struct {
	x0 int
	y0 int
	width int
	height int
}

func (r *Rect) drawRect(c color.RGBA, img *image.RGBA64) {
	for i:= r.x0; i <  r.x0 + r.width; i++ {
		for j:= r.y0; j < r.y0 + r.height; j++ {
			img.Set(i, j, c)
		}
	}
}
func (i *Identicon) drawIdenticon(img *image.RGBA64) {
	i.createPixelMatrix()
	r := Rect{0,0, i.pixelWidth, i.pixelHeight}

	for k:= 0; k < 5; k++ {
		for l:= 0; l < 5; l++ {
			if i.pixelMatrix[5*k+l] {
				r.drawRect(i.pixelColor, img)
			}
			r.x0 += i.pixelWidth
		}
		r.x0 = 0
		r.y0 += i.pixelHeight
	}
}
func main() {
	inputData := UserData{name: os.Args[1]}
	inputData.hashValue = md5.Sum([]byte(inputData.name))

	identicon := Identicon{pixelColor: color.RGBA{inputData.hashValue[1], inputData.hashValue[2], inputData.hashValue[3], 255}}

	identicon.pictureMatrix = createMatrix(inputData.hashValue)

	var w, h int = 280, 240

	var img = image.NewRGBA64(image.Rect(0,0,w,h))

	rectWidth := w / 5
	rectHeight := h / 5

	identicon.pixelWidth = rectWidth
	identicon.pixelHeight = rectHeight

	identicon.drawIdenticon(img)

	f, _ := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}