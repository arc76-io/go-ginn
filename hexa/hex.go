package hexa

import "fmt"

type Axial struct {
	Q, R int32
}

type Cuboid struct {
	Q, R, S int32
}

func (r Axial) Cuboid() Cuboid {
	return Cuboid{r.Q, r.R, -r.Q-r.R}
}

func (r Cuboid) Axial() Axial {
	return Axial{r.Q, r.R}
}

func (r Axial) String() string {
	return fmt.Sprintf("%s:%s", toRunes(r.Q), toRunes(r.R))
}

func (r Cuboid) String() string {
	return fmt.Sprintf("%s:%s:%s", toRunes(r.Q), toRunes(r.R), toRunes(r.S))
}

func (r Axial) Abs() Axial {
	return Axial{abs(r.Q), abs(r.R)}
}

func (r Cuboid) Abs() Cuboid {
	return Cuboid{abs(r.Q), abs(r.R), abs(r.S)}
}

func (r Axial) Add(o Axial) Axial {
	return Axial{r.Q + o.Q, r.R + o.R}
}

func (r Cuboid) Add(o Cuboid) Cuboid {
	return Cuboid{r.Q + o.Q, r.R + o.R, r.S + o.S}
}

func (r Axial) Subtract(o Axial) Axial {
	return Axial{r.Q - o.Q, r.R - o.R}
}

func (r Cuboid) Subtract(o Cuboid) Cuboid {
	return Cuboid{r.Q - o.Q, r.R - o.R, r.S - o.S}
}

func (r Axial) Multiply(o Axial) Axial {
	return Axial{r.Q * o.Q, r.R * o.R}
}

func (r Cuboid) Multiply(o Cuboid) Cuboid {
	return Cuboid{r.Q * o.Q, r.R * o.R, r.S * o.S}
}

func (r Axial) Length() int32 {
	return (abs(r.Q) + abs(r.R) + abs(-r.Q-r.R)) >> 1
}

func (r Cuboid) Length() int32 {
	return (abs(r.Q) + abs(r.R) + abs(r.S)) >> 1
}

func (r Axial) Distance(o Axial) int32 {
	return r.Subtract(o).Length();
}

func (r Cuboid) Distance(o Cuboid) int32 {
	return r.Subtract(o).Length();
}

func (r Axial) Neighbor(d Direction) Axial {
	return add(r, d.Delta())
}

func (r Axial) Direction(to Axial) Direction {
	return subtract(r, to).Direction()
}

func (r Cuboid) Direction() Direction {
	abs := r.Abs()
	if abs.Q >= abs.R && abs.Q >= abs.S {
		if r.Q < 0 {
			return DirectionNegQ
		}
		return DirectionPosQ
	}
	if abs.R >= abs.S {
		if r.R < 0 {
			return DirectionNegR
		}
		return DirectionPosR
	}
	if r.S < 0 {
		return DirectionNegS
	}
	return DirectionPosS
}