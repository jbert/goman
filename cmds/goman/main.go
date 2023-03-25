package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/cmplx"
	"sync"
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

	m := NewMandel(w, h)
	uiControls := makeUIControls(m)

	ui := container.New(layout.NewHBoxLayout(), canvasImage, uiControls)

	tickDur := 100 * time.Millisecond
	ticker := time.NewTicker(tickDur)

	win := a.NewWindow(title)
	win.SetContent(ui)
	win.Resize(fyne.NewSize(float32(w), float32(h)))

	tickMax := 100
	animTick := 0
	animateImage := func(img *image.RGBA, w, h int, ch <-chan time.Time) {
		for {
			m.UpdateMagMap(animTick, tickMax)
			clearImage(img)
			m.Draw(animTick, tickMax, img)
			ui.Refresh()
			animTick = (animTick + 1) % tickMax
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

	thresholdEntry := widget.NewEntryWithData(binding.FloatToString(binding.BindFloat(&(m.Threshold))))
	items = append(items, container.New(layout.NewHBoxLayout(), widget.NewLabel("threshold"), thresholdEntry))

	stepsEntry := widget.NewEntryWithData(binding.IntToString(binding.BindInt(&(m.Steps))))
	items = append(items, container.New(layout.NewHBoxLayout(), widget.NewLabel("steps"), stepsEntry))

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

func magnitudeToColour(animTick, tickMax int, mag float64) color.Color {
	animF := float64(animTick) / float64(tickMax)
	mag = math.Log(mag / animF)
	scale := 1.0
	f := math.Min(mag, scale) // 0 -> scale
	umax := 255.0 / scale

	cr := uint8(128 * animTick)
	cb := uint8(10 * (1.0 - animTick))
	y := uint8(umax * f)
	return color.YCbCr{y, cb, cr}
}

func (m *Mandel) Draw(animTick, tickMax int, img *image.RGBA) {
	w, h := widthHeight(img)
	fmt.Printf("%s: draw\n", time.Now())

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			c := magnitudeToColour(animTick, tickMax, m.MagMap[j][i])
			img.Set(i, j, c)
		}
	}
}

func (m *Mandel) UpdateMagMap(animTick, tickMax int) {
	//	start := time.Now()
	//	defer func() {
	//		fmt.Printf("Update took %s\n", time.Since(start))
	//	}()
	h := len(m.MagMap)
	w := len(m.MagMap[0])

	xw := m.Xhi - m.Xlo
	yw := m.Yhi - m.Ylo

	wg := sync.WaitGroup{}
	wg.Add(h)
	for j := 0; j < h; j++ {
		j := j
		go func() {
			y := m.Ylo + (yw * float64(j) / float64(h))
			for i := 0; i < w; i++ {
				x := m.Xlo + (xw * float64(i) / float64(w))
				m.MagMap[j][i] = m.calcPoint(x, y)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

type Mandel struct {
	MagMap [][]float64

	Steps     int
	Xlo, Xhi  float64
	Ylo, Yhi  float64
	Threshold float64
}

func NewMandel(w, h int) *Mandel {
	m := &Mandel{
		MagMap: make([][]float64, h),

		Steps:     100,
		Xlo:       -2.5,
		Xhi:       1.0,
		Ylo:       -1.5,
		Yhi:       1.5,
		Threshold: 1000,
	}
	for j := 0; j < h; j++ {
		m.MagMap[j] = make([]float64, w)
	}
	return m
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
