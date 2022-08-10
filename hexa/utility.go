package hexa

import (
	"fmt"
	"math"
	"strconv"
)

type Diagonal int
type Direction int

type Stats struct {
	Base		uint32	`json:"base"`
	Core		uint32	`json:"core"`
	Nano		uint32	`json:"nano"`
	Micro		uint32	`json:"micro"`
	Active	uint32	`json:"active"`
}

const (
	DiagonalPosQ Diagonal = iota
	DiagonalNegR Diagonal = iota
	DiagonalPosS Diagonal = iota
	DiagonalNegQ Diagonal = iota
	DiagonalPosR Diagonal = iota
	DiagonalNegS Diagonal = iota
	DiagonalNone Diagonal = iota
)

const (
	DirectionPosQ Direction = iota
	DirectionNegR Direction = iota
	DirectionPosS Direction = iota
	DirectionNegQ Direction = iota
	DirectionPosR Direction = iota
	DirectionNegS Direction = iota
	DirectionNone Direction = iota
)

var diagonals = []Cuboid {
	{2, -1, -1}, {1, -2, 1}, {-1, -1, 2},
	{-2, 1, 1}, {-1, 2, -1}, {1, 1, -2},
	{},
}

var directions = []Cuboid {
	{1, 0, -1}, {1, -1, 0}, {0, -1, 1},
	{-1, 0, 1}, {-1, 1, 0}, {0, 1, -1},
	{},
}

func (r Diagonal) Delta() Cuboid {
    return diagonals[r]
}

func (r Direction) Delta() Cuboid {
    return directions[r]
}

func (d Diagonal) String() string {
	switch d {
	case DiagonalPosQ:
		return "DiagonalPosQ"
	case DiagonalPosR:
		return "DiagonalPosR"
	case DiagonalPosS:
		return "DiagonalPosS"
	case DiagonalNegQ:
		return "DiagonalNegQ"
	case DiagonalNegR:
		return "DiagonalNegR"
	case DiagonalNegS:
		return "DiagonalNegS"
    default:
        return "DiagonalUndefined"
	}
}

func (r Direction) String() string {
	switch r {
	case DirectionPosQ:
		return "DirectionPosQ"
	case DirectionPosR:
		return "DirectionPosR"
	case DirectionPosS:
		return "DirectionPosS"
	case DirectionNegQ:
		return "DirectionNegQ"
	case DirectionNegR:
		return "DirectionNegR"
	case DirectionNegS:
		return "DirectionNegS"
	default:
		return "DirectionNone"
	}
}

func abs(x int32) int32 {
	if x > 0 {
		return x;
	}
	return -x;
}

func max(a, b int32) int32 {
    if a > b {
        return a
    }
    return b
}

func min(a, b int32) int32 {
    if a < b {
        return a
    }
    return b
}

func toRunes(v int32) string {
    if v >= 0 {
        return fmt.Sprintf("P%d", abs(v))
    }
    return fmt.Sprintf("N%d", v)
}

func fromRunes(v string) int32 {
    s := v[:1]
    n, err := strconv.Atoi(v[1:])
    if err != nil {
        fmt.Println(err)
        return 0;
    }
    if s == "N" {
        return int32(-n)
    }
    return int32(n)
}

func unfloat(x, y, z float64) Axial {
	rx, ry, rz := math.Round(x), math.Round(y), math.Round(z)
	dx, dy, dz := math.Abs(rx-x), math.Abs(ry-y), math.Abs(rz-z)

	if dx > dz && dx > dy {
		rx = -rz - ry
	} else if dz > dy {
		rz = -rx - ry
	} else {
		ry = -rx - rz
	}
	return Axial{
		Q: int32(math.Round(rx)),
		R: int32(math.Round(rz)),
	}
}

func add(a Axial, b Cuboid) Axial {
	return Axial {
		Q: a.Q + b.Q,
		R: a.R + b.R,
	}
}

func subtract(a, b Axial) Cuboid {
	return Cuboid {
		Q: a.Q - b.Q,
		R: a.R - b.R,
		S: -(a.Q - b.Q) - (a.R - b.R),
	}
}

func multiply(d Cuboid, k int32) Cuboid {
	return Cuboid {d.Q * k, d.R * k, d.S * k}
}

func Line(a, b Axial) []Axial {
	delta := subtract(a, b)
	n := delta.Length()
	dir := delta.Direction()

	results := make([]Axial, 0, n)
	visited := make(map[Axial]bool, n)

	ax, ay, az := float64(a.Q), float64(-a.Q-a.R), float64(a.R)
	bx, by, bz := float64(b.Q), float64(-b.Q-b.R), float64(b.R)
	x, y, z := bx-ax, by-ay, bz-az

	step := 1. / float64(n)
	for h := int32(0); h <= n; h++ {
		t := step * float64(h)
		pnt := unfloat(ax+x*t, ay+y*t, az+z*t)
		for visited[pnt] {
			pnt = pnt.Neighbor(dir)
		}
		results = append(results, pnt)
		visited[pnt] = true
	}
	if !visited[b] {
		results = append(results, b)
	}

	return results
}

func Range(h Axial, rad int32) map[Axial]bool {
	results := make(map[Axial]bool, rad*rad)
	if rad < 1 {
		return results
	}
	for x := -rad; x <= rad; x++ {
		for y := max(-rad, -x-rad); y <= min(rad, -x+rad); y++ {
			z := -x - y
			delta := Cuboid{
				Q: x,
				R: z,
				S: y,
			}
			results[add(h, delta)] = true
		}
	}
	return results
}

func Ring(h Axial, rad int32) map[Axial]bool {
	results := make(map[Axial]bool)
	if rad < 1 {
		return results
	}

	h = add(h, multiply(DirectionPosS.Delta(), rad))
	results[h] = true
	if rad > 1 {
		for i := 0; i < 6; i++ {
			for j := int32(0); j < rad; j++ {
				h = add(h, Direction(i).Delta())
				results[h] = true
			}
		}
	}
	return results
}