package particle

type Vec2 struct {
	X float64
	Y float64
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
