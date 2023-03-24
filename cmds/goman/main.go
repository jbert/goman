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
			clearImage(img, w, h)
			generateImage(img, w, h, tick)
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

func makeImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{w, h}})
	return img
}

func clearImage(img *image.RGBA, w, h int) {
	draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
}

func generateImage(img *image.RGBA, w, h, tick int) {
	drawMandelbrot(img, w, h, tick)
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

func drawMandelbrot(img *image.RGBA, w, h, tick int) {
	fmt.Printf("%s: draw %d\n", time.Now(), tick)
	steps := 100
	xlo := -2.5
	xhi := 1.0
	xw := xhi - xlo
	ylo := -1.5
	yhi := 1.5
	yw := yhi - ylo

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			x := xlo + (xw * float64(i) / float64(w))
			y := ylo + (yw * float64(j) / float64(h))
			mag := mandelCalc(x, y, steps)
			//			fmt.Printf("x %v y %v mag %v\n", x, y, mag)
			c := color.RGBA{0, 0, 0, 255}
			if mag < 1.0 {
				c.R = 255
			}
			img.Set(i, j, c)
		}
	}
}

func mandelCalc(x, y float64, steps int) float64 {
	threshold := 1000.0

	c := complex(x, y)

	z := complex(0, 0)
	for steps > 0 {
		//		fmt.Printf("z %v c %v\n", z, c)
		z = z*z + c
		if cmplx.Abs(z) > threshold {
			return threshold
		}
		steps--
	}
	return cmplx.Abs(z)
}
