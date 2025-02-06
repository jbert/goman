package main

import (
	"image"
	"image/draw"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/jbert/goman"
)

type tappableRasterImage struct {
	widget.BaseWidget
	img      *canvas.Image
	onTapped func(x, y float64)
}

func NewTappableRasterImage(img *canvas.Image) *tappableRasterImage {
	w := &tappableRasterImage{BaseWidget: widget.BaseWidget{}, img: img}
	w.ExtendBaseWidget(w)
	return w
}

func (tri *tappableRasterImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(tri.img)
}

func (tri *tappableRasterImage) SetTapped(f func(x, y float64)) {
	tri.onTapped = f
}

func (tri *tappableRasterImage) Tapped(e *fyne.PointEvent) {
	//	size := tri.Size()
	log.Printf("Tapped!: %+v", e)
	if tri.onTapped != nil {
		xProportion := float64(e.AbsolutePosition.X) / float64(tri.img.Size().Width)
		yProportion := float64(e.AbsolutePosition.Y) / float64(tri.img.Size().Height)
		tri.onTapped(xProportion, yProportion)
	}
}

func main() {
	a := app.New()

	title := "Go, Man"
	w := 640
	h := 480

	imgs := make([]*image.RGBA, 2)
	imgs[0] = makeImage(w, h)
	imgs[1] = makeImage(w, h)
	currentImg := 0

	m := goman.NewMandel(w, h)
	uiControls := makeUIControls(m)

	rasterImages := make([]*tappableRasterImage, 2)

	rasterImages[0] = NewTappableRasterImage(canvas.NewImageFromImage(imgs[0]))
	rasterImages[0].img.SetMinSize(fyne.Size{Width: float32(w), Height: float32(h)})
	rasterImages[0].Show()
	rasterImages[0].SetTapped(m.OnTap)

	rasterImages[1] = NewTappableRasterImage(canvas.NewImageFromImage(imgs[1]))
	rasterImages[1].img.SetMinSize(fyne.Size{Width: float32(w), Height: float32(h)})
	rasterImages[1].Hide()
	rasterImages[1].SetTapped(m.OnTap)

	ui := container.New(layout.NewHBoxLayout(), rasterImages[0], rasterImages[1], uiControls)

	tickDur := 50 * time.Millisecond
	ticker := time.NewTicker(tickDur)

	win := a.NewWindow(title)
	win.SetContent(ui)
	win.Resize(fyne.NewSize(float32(w), float32(h)))

	tickMax := 100
	animTick := 0
	animateImage := func(ch <-chan time.Time) {
		for {
			lastImg := currentImg

			currentImg = (currentImg + 1) % 2
			m.UpdateMagMap(animTick, tickMax)
			clearImage(imgs[currentImg])
			m.Draw(animTick, tickMax, imgs[currentImg])

			rasterImages[lastImg].Hide()
			rasterImages[currentImg].Show()
			ui.Refresh()
			animTick = (animTick + 1) % tickMax
			<-ch
		}
	}
	go animateImage(ticker.C)

	win.ShowAndRun()
}

func makeUIControls(m *goman.Mandel) fyne.CanvasObject {
	items := make([]fyne.CanvasObject, 0)

	items = append(items, widget.NewLabel("Settings"))
	items = append(items, makeOptionsForm(m))

	items = append(items, widget.NewButton("+", m.ZoomIn))
	items = append(items, widget.NewButton("-", m.ZoomOut))

	return container.New(layout.NewVBoxLayout(), items...)
}

func makeOptionsForm(m *goman.Mandel) fyne.CanvasObject {

	formData := binding.BindStruct(m)
	keys := formData.Keys()
	items := make([]*widget.FormItem, 0)

	for _, k := range keys {
		data, err := formData.GetItem(k)
		if err != nil {
			items = append(items, widget.NewFormItem(k, widget.NewLabel(err.Error())))
		} else {
			items = append(items, widget.NewFormItem(k, createBoundItem(data)))
		}
	}

	form := widget.NewForm(items...)

	return form
}

func createBoundItem(v binding.DataItem) fyne.CanvasObject {
	switch val := v.(type) {
	case binding.Bool:
		return widget.NewCheckWithData("", val)
	case binding.Float:
		return widget.NewEntryWithData(binding.FloatToString(val))
		//		s := widget.NewSliderWithData(0, 1, val)
		//		s.Step = 0.01
		//		return s
	case binding.Int:
		return widget.NewEntryWithData(binding.IntToString(val))
	case binding.String:
		return widget.NewEntryWithData(val)
	default:
		return widget.NewLabel("")
	}
}

func makeImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{w, h}})
	return img
}

func clearImage(img *image.RGBA) {
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)
}
