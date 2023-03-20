package main

import (
	"image"
	"image/color"
	"image/draw"
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

	tickDur := 10 * time.Millisecond
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
	x2 := w / 2
	y2 := h / 2
	c := color.RGBA{255, 0, 0, 255}
	offset := tick - w/4
	for i := -w / 4; i < w/4; i++ {
		img.Set(x2+i+offset, y2+i, c)
	}
}
