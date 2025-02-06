package goman

import (
	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

type Pt[T Numeric] struct {
	X T
	Y T
}

func NewPt[T Numeric](x, y T) Pt[T] {
	return Pt[T]{X: x, Y: y}
}

func (p Pt[T]) Zero() Pt[T] {
	return Pt[T]{}
}

func (p Pt[T]) Add(q Pt[T]) Pt[T] {
	return Pt[T]{p.X + q.X, p.Y + q.Y}
}

func (p Pt[T]) Sub(q Pt[T]) Pt[T] {
	return Pt[T]{p.X - q.X, p.Y - q.Y}
}

func (p Pt[T]) Scale(v T) Pt[T] {
	return Pt[T]{p.X * v, p.Y * v}
}

func (p Pt[T]) Inv(v T) Pt[T] {
	return p.Zero().Sub(p)
}

type Rect[T Numeric] struct {
	TopLeft  Pt[T]
	BotRight Pt[T]
}

func NewRectCorners[T Numeric](topleft, botRight Pt[T]) Rect[T] {
	return Rect[T]{TopLeft: topleft, BotRight: botRight}
}

func NewRectCentered[T Numeric](centre, size Pt[T]) Rect[T] {
	size.X /= 2
	size.Y /= 2
	return Rect[T]{TopLeft: centre.Sub(size), BotRight: centre.Add(size)}
}

func (r Rect[T]) ChangeSizeCentered(newSize Pt[T]) Rect[T] {
	centre := r.Centre()
	return NewRectCentered(centre, newSize)
}

func (r Rect[T]) Size() Pt[T] {
	return r.BotRight.Sub(r.TopLeft)
}

func (r Rect[T]) HalfSize() Pt[T] {
	step := r.Size()
	step.X /= 2
	step.Y /= 2
	return step
}

func (r Rect[T]) Centre() Pt[T] {
	return r.TopLeft.Add(r.HalfSize())
}
