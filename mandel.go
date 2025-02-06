package goman

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/cmplx"
	"sync"
)

type Mandel struct {
	magMap    [][]float64
	zoomScale float64
	Threshold float64
	Steps     int

	View Rect[float64]
}

func NewMandel(w, h int) *Mandel {
	//	topLeft := NewPt(-2.5, -1.5)
	//	botRight := NewPt(1.0, 1.5)
	centre := NewPt(0.0, 0.0)
	xWidth := 2.0
	yHeight := float64(h) / float64(w) * xWidth
	size := NewPt(xWidth, yHeight)

	m := &Mandel{
		magMap:    make([][]float64, h),
		zoomScale: 1.2,

		Steps: 100,

		View:      NewRectCentered(centre, size),
		Threshold: 1000,
	}
	for j := 0; j < h; j++ {
		m.magMap[j] = make([]float64, w)
	}
	return m
}

func (m *Mandel) OnScroll(xProp, yProp, dxProp, dyProp float64) {
	log.Printf("Scrolled: x %v y %v dx %v dy %v", xProp, yProp, dxProp, dyProp)
	log.Printf("Scrolled: Old View: %+v", m.View)
	if dyProp > 0 {
		newCentre := m.proportionToPt(xProp, yProp)
		m.setCentre(newCentre)
		m.ZoomIn()
	} else {
		m.ZoomOut()
	}
	log.Printf("Scrolled: New View: %+v", m.View)
}

func (m *Mandel) OnTap(xProportion, yProportion float64) {
	log.Printf("Tapped: x %v y %v", xProportion, yProportion)
	newCentre := m.proportionToPt(xProportion, yProportion)
	m.View = NewRectCentered(newCentre, m.View.HalfSize())
}

func (m *Mandel) setCentre(newCentre Pt[float64]) {
	size := m.View.Size()
	m.View = NewRectCentered(newCentre, size)
}

func (m *Mandel) getCentre() Pt[float64] {
	return m.View.Centre()
}

func (m *Mandel) calcPoint(p Pt[float64]) float64 {
	c := complex(p.X, p.Y)

	z := complex(0, 0)
	steps := m.Steps
	for steps > 0 {
		z = z*z + c
		if cmplx.Abs(z) > m.Threshold {
			return m.Threshold
		}
		steps--
	}
	return cmplx.Abs(z)
}

func (m *Mandel) ZoomOut() {
	centre := m.View.Centre()
	size := m.View.Size().Scale(m.zoomScale)
	m.View = NewRectCentered(centre, size)
}

func (m *Mandel) ZoomIn() {
	centre := m.View.Centre()
	size := m.View.Size().Scale(1 / m.zoomScale)
	m.View = NewRectCentered(centre, size)
}

func (m *Mandel) proportionToPt(xProportion, yProportion float64) Pt[float64] {
	p := m.View.Size()
	p.X *= xProportion
	p.Y *= yProportion
	return p.Add(m.View.TopLeft)
}

func widthHeight(img image.Image) (int, int) {
	bounds := img.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y
	return w, h
}

func (m *Mandel) Draw(animTick, tickMax int, img *image.RGBA) {
	w, h := widthHeight(img)

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			c := magnitudeToColour(animTick, tickMax, m.magMap[j][i])
			img.Set(i, j, c)
		}
	}
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

func (m *Mandel) UpdateMagMap(animTick, tickMax int) {
	h := len(m.magMap)
	w := len(m.magMap[0])

	wg := sync.WaitGroup{}
	wg.Add(h)
	for j := 0; j < h; j++ {
		j := j
		go func() {
			//			y := m.Y + (m.Height * float64(j) / float64(h))
			yProp := float64(j) / float64(h)
			for i := 0; i < w; i++ {
				xProp := float64(i) / float64(w)
				p := m.proportionToPt(xProp, yProp)
				m.magMap[j][i] = m.calcPoint(p)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
