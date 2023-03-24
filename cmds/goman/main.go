package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/cmplx"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func main() {
	a := app.New()

	title := "Go, Man"
	w := 640
	h := 480

	win := a.NewWindow(title)

	img := makeImage(w, h)

	tickDur := 100 * time.Millisecond
	ticker := time.NewTicker(tickDur)

	fyneImg := canvas.NewImageFromImage(img)
	win.SetContent(fyneImg)
	win.Resize(fyne.NewSize(float32(w), float32(h)))

	animateImage := func(img *image.RGBA, w, h int, ch <-chan time.Time) {
		tick := 0
		maxTick := w / 2

		for {
			clearImage(img)
			generateImage(img, tick)
			canvas.Refresh(fyneImg)
			tick++
			tick = tick % maxTick
			break
			<-ch
		}
	}
	go animateImage(img, w, h, ticker.C)

	win.ShowAndRun()
}

func widthHeight(img image.Image) (int, int) {
	bounds := img.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y
	return w, h
}

func makeImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{w, h}})
	return img
}

func clearImage(img *image.RGBA) {
	draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
}

func generateImage(img *image.RGBA, tick int) {
	m := NewMandel()
	m.Draw(img)
	/*
		x2 := w / 2
		y2 := h / 2
		c := color.RGBA{255, 0, 0, 255}
		offset := tick - w/4
		for i := -w / 4; i < w/4; i++ {
			img.Set(x2+i+offset, y2+i, c)
		}
	*/
}

func (m *Mandel) Draw(img *image.RGBA) {
	w, h := widthHeight(img)
	fmt.Printf("%s: draw %d\n", time.Now())

	xw := m.Xhi - m.Xlo
	yw := m.Yhi - m.Ylo

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			x := m.Xlo + (xw * float64(i) / float64(w))
			y := m.Ylo + (yw * float64(j) / float64(h))
			mag := m.calcPoint(x, y)
			//			fmt.Printf("x %v y %v mag %v\n", x, y, mag)
			c := color.RGBA{0, 0, 0, 255}
			if mag < 1.0 {
				c.R = 255
			}
			img.Set(i, j, c)
		}
	}
}

type Mandel struct {
	Steps     int
	Xlo, Xhi  float64
	Ylo, Yhi  float64
	Threshold float64
}

func NewMandel() *Mandel {
	return &Mandel{
		Steps:     100,
		Xlo:       -2.5,
		Xhi:       1.0,
		Ylo:       -1.5,
		Yhi:       1.5,
		Threshold: 1000,
	}
}

func (m *Mandel) calcPoint(x, y float64) float64 {
	c := complex(x, y)

	z := complex(0, 0)
	steps := m.Steps
	for steps > 0 {
		//		fmt.Printf("z %v c %v\n", z, c)
		z = z*z + c
		if cmplx.Abs(z) > m.Threshold {
			return m.Threshold
		}
		steps--
	}
	return cmplx.Abs(z)
}
