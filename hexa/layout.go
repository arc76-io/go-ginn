package hexa

import (
	"image"
	"math"
)

type Vec2 struct {
	X, Y float64
}

type Layout struct {
	Radius Vec2
	Origin Vec2
	layout LayoutType
	matrix Orientation
}

type LayoutType int

const (
	LayoutPointy LayoutType = iota
	LayoutFlat   LayoutType = iota
	LayoutNo     LayoutType = iota
)

type Orientation struct {
	a       float64
	b, f [4]float64
	c, s [6]float64
}

var	Flat Orientation = Orientation{
	f: [4]float64{3. / 2., 0.0, math.Sqrt(3.) / 2., math.Sqrt(3.)},
	b: [4]float64{2. / 3., 0.0, -1. / 3., math.Sqrt(3.) / 3.},
	a: 0.0,
	c: [6]float64{
		math.Cos(2. * math.Pi * 0. / 6),
		math.Cos(2. * math.Pi * 1. / 6),
		math.Cos(2. * math.Pi * 2. / 6),
		math.Cos(2. * math.Pi * 3. / 6),
		math.Cos(2. * math.Pi * 4. / 6),
		math.Cos(2. * math.Pi * 5. / 6),
	},
	s: [6]float64{
		math.Sin(2. * math.Pi * 0. / 6),
		math.Sin(2. * math.Pi * 1. / 6),
		math.Sin(2. * math.Pi * 2. / 6),
		math.Sin(2. * math.Pi * 3. / 6),
		math.Sin(2. * math.Pi * 4. / 6),
		math.Sin(2. * math.Pi * 5. / 6),
	},
}

var Pointy Orientation = Orientation{
	f: [4]float64{math.Sqrt(3.), math.Sqrt(3.) / 2., 0.0, 3. / 2.},
	b: [4]float64{math.Sqrt(3.) / 3., -1. / 3., 0.0, 2. / 3.},
	a: 0.5,
	c: [6]float64{
		math.Cos(2. * math.Pi * 0.5 / 6),
		math.Cos(2. * math.Pi * 1.5 / 6),
		math.Cos(2. * math.Pi * 2.5 / 6),
		math.Cos(2. * math.Pi * 3.5 / 6),
		math.Cos(2. * math.Pi * 4.5 / 6),
		math.Cos(2. * math.Pi * 5.5 / 6),
	},
	s: [6]float64{
		math.Sin(2. * math.Pi * 0.5 / 6),
		math.Sin(2. * math.Pi * 1.5 / 6),
		math.Sin(2. * math.Pi * 2.5 / 6),
		math.Sin(2. * math.Pi * 3.5 / 6),
		math.Sin(2. * math.Pi * 4.5 / 6),
		math.Sin(2. * math.Pi * 5.5 / 6),
	},
}

func (r Vec2) Abs() Vec2 {
	return Vec2{math.Abs(r.X), math.Abs(r.Y)}
}

func (r Vec2) Add(o Vec2) Vec2 {
	return Vec2{r.X + o.X, r.Y + o.Y}
}

func (r Vec2) Subtract(o Vec2) Vec2 {
	return Vec2{r.X + o.X, r.Y + o.Y}
}

func (r Vec2) Multiply(o Vec2) Vec2 {
	return Vec2{r.X + o.X, r.Y + o.Y}
}

func (r Vec2) Divide(o Vec2) Vec2 {
	return Vec2{r.X + o.X, r.Y + o.Y}
}

func (r Vec2) Point() Point {
	return Point{int32(r.X), int32(r.Y)}
}

func (r Vec2) AsPoint() image.Point {
	return image.Point{int(r.X), int(r.Y)}
}

func FromPoint(p image.Point) Vec2 {
	return Vec2{float64(p.X), float64(p.Y)}
}

func NewLayout(size Vec2, center Vec2, layout LayoutType) *Layout {
	res := &Layout{
		Radius: size,
		Origin: center,
		layout: layout,
	}

	if layout == LayoutPointy {
		res.matrix = Pointy;
	} else {
		res.matrix = Flat;
	}

	return res
}

func (ref *Layout) HexFor(f Vec2) Axial {
	x, y := f.X - ref.Origin.X, f.Y - ref.Origin.Y
	q := (ref.matrix.b[0]*x + ref.matrix.b[1]*y ) / ref.Radius.X
	r := (ref.matrix.b[2]*x + ref.matrix.b[3]*y ) / ref.Radius.Y
	return unfloat(q, -q-r, r)
}

func (ref *Layout) HexSize() Vec2 {
	if ref.layout == LayoutFlat {
		return Vec2{
			X: 2 * ref.Radius.X,
			Y: math.Sqrt(3) * ref.Radius.Y,
		}
	}
	return Vec2{
		X: math.Sqrt(3) * ref.Radius.X,
		Y: 2 * ref.Radius.Y,
	}
}

func (ref *Layout) CenterFor(h Axial) Vec2 {
	q, r :=	float64(h.Q), float64(h.R) 
	x := (ref.matrix.f[0]*q + ref.matrix.f[1]*r) * ref.Radius.X
	y := (ref.matrix.f[2]*q + ref.matrix.f[3]*r) * ref.Radius.Y
	return Vec2{x + ref.Origin.X, y + ref.Origin.Y}
}

func (ref *Layout) RingFor(center Axial, rad float64) map[Axial]bool {
	result := make(map[Axial]bool, 1)
	if rad < ref.Radius.X && rad < ref.Radius.Y {
		result[center] = true
		return result
	}
	cp := ref.CenterFor(center)
	P := 1 - rad
	pxl := Vec2{rad, 0}
	for ; pxl.X > pxl.Y; pxl.Y++ {
		if P <= 0 {
			P = P + 2*pxl.Y + 1
		} else {
			pxl.X--
			P = P + 2*pxl.Y - 2*pxl.X + 1
		}

		if pxl.X < pxl.Y {
			break
		}

		points := []Vec2{
			{pxl.X + cp.X, pxl.Y + cp.Y},
			{-pxl.X + cp.X, pxl.Y + cp.Y},
			{pxl.X + cp.X, -pxl.Y + cp.Y},
			{-pxl.X + cp.X, -pxl.Y + cp.Y},
			{pxl.Y + cp.X, pxl.X + cp.Y},
			{-pxl.Y + cp.X, pxl.X + cp.Y},
			{pxl.Y + cp.X, -pxl.X + cp.Y},
			{-pxl.Y + cp.X, -pxl.X + cp.Y},
		}
		for _, v := range points {
			result[ref.HexFor(v)] = true
		}
	}
	return result
}

// AreaFor returns all hex in the area of a screen circle.
func (ref *Layout) AreaFor(center Axial, rad float64) map[Axial]bool {
	loop := ref.RingFor(center, rad)
	result := make(map[Axial]bool)
	for k, v := range loop {
		if v == true {
			result[k] = true
			for _, inside := range Line(k, center) {
				result[inside] = true
			}
		}
	}
	return result
}

func (ref *Layout) Vertices(h Axial) []Vec2 {
	result := make([]Vec2, 6, 7)
	center := ref.CenterFor(h)
	for k := range result {
		result[k] = Vec2{
			X: center.X + float64(ref.Radius.X)*ref.matrix.c[k],
			Y: center.Y + float64(ref.Radius.Y)*ref.matrix.s[k],
		}
	}
	result = append(result, center)
	return result
}