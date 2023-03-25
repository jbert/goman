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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	title := "Go, Man"
	w := 640
	h := 480

	img := makeImage(w, h)
	canvasImage := canvas.NewImageFromImage(img)
	// Unclear why setting FillMode doesn't achieve this
	canvasImage.SetMinSize(fyne.Size{float32(w), float32(h)})
	//	canvasImage.FillMode = canvas.ImageFillOriginal

	m := NewMandel()
	uiControls := makeUIControls(m)

	ui := container.New(layout.NewHBoxLayout(), canvasImage, uiControls)

	tickDur := 100 * time.Millisecond
	ticker := time.NewTicker(tickDur)

	win := a.NewWindow(title)
	win.SetContent(ui)
	win.Resize(fyne.NewSize(float32(w), float32(h)))

	animateImage := func(img *image.RGBA, w, h int, ch <-chan time.Time) {
		tick := 0
		maxTick := w / 2

		for {
			clearImage(img)
			generateImage(m, img, tick, maxTick)
			canvasImage.Refresh()
			ui.Refresh()
			tick++
			tick = tick % maxTick
			<-ch
		}
	}
	go animateImage(img, w, h, ticker.C)

	win.ShowAndRun()
}

func makeUIControls(m *Mandel) fyne.CanvasObject {
	var items []fyne.CanvasObject

	items = append(items, widget.NewLabel("Settings"))

	xloEntry := widget.NewEntryWithData(binding.FloatToString(binding.BindFloat(&(m.Xlo))))
	items = append(items, container.New(layout.NewHBoxLayout(), widget.NewLabel("xlo"), xloEntry))

	xhiEntry := widget.NewEntryWithData(binding.FloatToString(binding.BindFloat(&(m.Xhi))))
	items = append(items, container.New(layout.NewHBoxLayout(), widget.NewLabel("xhi"), xhiEntry))

	yloEntry := widget.NewEntryWithData(binding.FloatToString(binding.BindFloat(&(m.Ylo))))
	items = append(items, container.New(layout.NewHBoxLayout(), widget.NewLabel("ylo"), yloEntry))

	yhiEntry := widget.NewEntryWithData(binding.FloatToString(binding.BindFloat(&(m.Yhi))))
	items = append(items, container.New(layout.NewHBoxLayout(), widget.NewLabel("yhi"), yhiEntry))

	return container.New(layout.NewVBoxLayout(), items...)
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

func generateImage(m *Mandel, img *image.RGBA, tick, maxTick int) {
	fmt.Printf("%s: generate [%d/%d]\n", time.Now(), tick, maxTick)
	offset := float64(tick) / float64(maxTick)
	m.Xlo += offset
	m.Xhi += offset
	m.Draw(img)
}

func (m *Mandel) Draw(img *image.RGBA) {
	w, h := widthHeight(img)
	fmt.Printf("%s: draw\n", time.Now())

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
