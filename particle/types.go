package particle

import "math"

type Vec2 [2]float64

func Vec2Mulf(v1 *Vec2, v float64) {
	v1[0] *= v
	v1[1] *= v
}

func Vec2Add(v1 *Vec2, v2 *Vec2) {
	v1[0] += v2[0]
	v1[1] += v2[1]
}

func Vec2Sub(v1 *Vec2, v2 *Vec2) {
	v1[0] -= v2[0]
	v1[1] -= v2[1]
}

func Vec2Normalize(v *Vec2) {
	ln := 1.0 / math.Sqrt(float64(v[0]*v[0]+v[1]*v[1]))
	v[0] *= ln
	v[1] *= ln
}

type Point struct {
	X int
	Y int
}

type Size struct {
	Width  uint32
	Height uint32
}

type Range struct {
	Min uint32
	Max uint32
}

type Cfg struct {
	Cnt uint32

	Seed uint64

	MinRadius float64
	AvrRadius float64
	MaxRadius float64

	MaxChecks  uint32
	Iterations uint32
}
